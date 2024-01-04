package primary

import (
	"flag"
	"log"

	domain "raftinstance/domain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Driver() {
	addr := flag.String("raft", "", "raft server address")
	join := flag.String("join", "", "join cluster address")
	api := flag.String("api", "", "api server address")
	state := flag.String("state_dir", "", "raft state directory (WAL, Snapshots)")
	flag.Parse()
	domain.StartRaft(addr, state, join)
	e := echo.New()

	// Echo middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define routes
	e.PUT("/", domain.Save)
	e.GET("/:key", domain.Get)
	e.GET("/mgmt/nodes", domain.Nodes)
	e.DELETE("/mgmt/nodes/:id", domain.RemoveNode)

	// Start the Echo HTTP server
	if err := e.Start(*api); err != nil {
		log.Fatal(err)
	}
}
