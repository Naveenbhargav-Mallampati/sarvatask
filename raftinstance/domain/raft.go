package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	secondary "raftinstance/adapters/secondary"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shaj13/raft"
	"github.com/shaj13/raft/transport"
	"github.com/shaj13/raft/transport/raftgrpc"
	"google.golang.org/grpc"
)

// var e = echo.New()

type entry struct {
	Key   string
	Value string
}

var (
	node *raft.Node
	fsm  *stateMachine
)

func StartRaft(addr, state, join *string) {
	var (
		opts      []raft.Option
		startOpts []raft.StartOption
	)

	startOpts = append(startOpts, raft.WithAddress(*addr))
	opts = append(opts, raft.WithStateDIR(*state))
	if *join != "" {
		opt := raft.WithFallback(
			raft.WithJoin(*join, time.Second),
			raft.WithRestart(),
		)
		startOpts = append(startOpts, opt)
	} else {
		opt := raft.WithFallback(
			raft.WithInitCluster(),
			raft.WithRestart(),
		)
		startOpts = append(startOpts, opt)
	}

	raftgrpc.Register(
		raftgrpc.WithDialOptions(grpc.WithInsecure()),
	)
	fsm = newstateMachine()
	node = raft.NewNode(fsm, transport.GRPC, opts...)
	raftServer := grpc.NewServer()
	raftgrpc.RegisterHandler(raftServer, node.Handler())

	// Start the gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", *addr)
		if err != nil {
			log.Fatal(err)
		}

		err = raftServer.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Start the Raft node
	go func() {
		err := node.Start(startOpts...)
		if err != nil && err != raft.ErrNodeStopped {
			log.Fatal(err)
		}
	}()
}

func Get(c echo.Context) error {
	key := c.Param("key")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second)
	defer cancel()

	if err := node.LinearizableRead(ctx); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	value := fsm.Read(key)

	response := struct {
		Name string `json:"name"`
		Size string `json:"size"`
	}{
		Name: key,
		Size: value,
	}

	return c.JSON(http.StatusOK, response)
}

func Nodes(c echo.Context) error {
	raws := []raft.RawMember{}
	membs := node.Members()
	for _, m := range membs {
		raws = append(raws, m.Raw())
	}

	return c.JSON(http.StatusOK, raws)
}

func RemoveNode(c echo.Context) error {
	sid := c.Param("id")
	id, err := strconv.ParseUint(sid, 0, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second)
	defer cancel()

	if err := node.RemoveMember(ctx, id); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func Save(c echo.Context) error {
	buf, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := json.Unmarshal(buf, new(entry)); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second)
	defer cancel()
	var data entry

	err = json.Unmarshal(buf, &data)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "Error decoding JSON request body")
	}
	data1 := secondary.GetData(data.Key)
	fmt.Println(data1)
	if data1 == data.Value {
		response := struct {
			Message string `json:"message"`
			Data    string `json:"data"`
		}{
			Message: "Duplicate:K/V",
			Data:    fmt.Sprintf("%s: %s", data.Key, data.Value),
		}

		return c.JSON(http.StatusOK, response)
	}

	if err := node.Replicate(ctx, buf); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	stat := secondary.SetData(data.Key, data.Value)

	response := struct {
		Message string `json:"message"`
		Data    bool   `json:"status"`
	}{
		Message: "Update",
		Data:    stat,
	}

	return c.JSON(http.StatusOK, response)
}
