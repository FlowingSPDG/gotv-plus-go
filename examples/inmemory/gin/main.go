package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/FlowingSPDG/gotv-plus-go/examples/inmemory"
	"github.com/FlowingSPDG/gotv-plus-go/gotv"
)

var (
	auth string
	port int
)

func main() {
	flag.StringVar(&auth, "auth", "SuperSecureStringDoNotShare", "tv_broadcast_origin_auth \"SuperSecureStringDoNotShare\"")
	flag.IntVar(&port, "port", 8080, "Port to listen")
	flag.Parse()

	m := inmemory.NewInmemoryGOTV(auth)
	app := gin.Default()
	g := app.Group("/gotv") // /gotv
	gotv.SetupStoreHandlersGin(m, g)
	gotv.SetupBroadcasterHandlersGin(m, g)

	p := fmt.Sprintf("%s:%d", "", port)

	// Start server
	log.Println("Start listening on:", p)
	if err := app.Run(p); err != nil {
		panic(err)
	}
}
