package main

// import (
// 	"flag"
// 	"log"
// 	"net"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"github.com/shaj13/raft"
// 	"github.com/shaj13/raft/transport"
// 	"github.com/shaj13/raft/transport/raftgrpc"
// 	"google.golang.org/grpc"
// )

// type entry struct {
// 	Key   string
// 	Value string
// }

// var (
// 	node *raft.Node
// 	fsm  *stateMachine
// )

// func Driver() {
// 	addr := flag.String("raft", "", "raft server address")
// 	join := flag.String("join", "", "join cluster address")
// 	api := flag.String("api", "", "api server address")
// 	state := flag.String("state_dir", "", "raft state directory (WAL, Snapshots)")
// 	flag.Parse()

// 	router := mux.NewRouter()
// 	router.HandleFunc("/", http.HandlerFunc(save)).Methods("PUT", "POST")
// 	router.HandleFunc("/{key}", http.HandlerFunc(get)).Methods("GET")
// 	router.HandleFunc("/mgmt/nodes", http.HandlerFunc(nodes)).Methods("GET")
// 	router.HandleFunc("/mgmt/nodes/{id}", http.HandlerFunc(removeNode)).Methods("DELETE")

// 	var (
// 		opts      []raft.Option
// 		startOpts []raft.StartOption
// 	)

// 	startOpts = append(startOpts, raft.WithAddress(*addr))
// 	opts = append(opts, raft.WithStateDIR(*state))
// 	if *join != "" {
// 		opt := raft.WithFallback(
// 			raft.WithJoin(*join, time.Second),
// 			raft.WithRestart(),
// 		)
// 		startOpts = append(startOpts, opt)
// 	} else {
// 		opt := raft.WithFallback(
// 			raft.WithInitCluster(),
// 			raft.WithRestart(),
// 		)
// 		startOpts = append(startOpts, opt)
// 	}

// 	raftgrpc.Register(
// 		raftgrpc.WithDialOptions(grpc.WithInsecure()),
// 	)
// 	fsm = newstateMachine()
// 	node = raft.NewNode(fsm, transport.GRPC, opts...)
// 	raftServer := grpc.NewServer()
// 	raftgrpc.RegisterHandler(raftServer, node.Handler())

// 	go func() {
// 		lis, err := net.Listen("tcp", *addr)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		err = raftServer.Serve(lis)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	go func() {
// 		err := node.Start(startOpts...)
// 		if err != nil && err != raft.ErrNodeStopped {
// 			log.Fatal(err)
// 		}
// 	}()

// 	if err := http.ListenAndServe(*api, router); err != nil {
// 		log.Fatal(err)
// 	}
// }
