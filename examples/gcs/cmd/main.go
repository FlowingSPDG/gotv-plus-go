package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/FlowingSPDG/gotv-plus-go/examples/gcs"
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

	ctx := context.Background()
	s, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	m := gcs.NewCloudStorageGOTV(s, auth, 8)
	app := fiber.New()
	g := app.Group("/gotv") // /gotv
	g.Use(logger.New())
	gotv.SetupStoreHandlersFiber(m, g)
	gotv.SetupBroadcasterHandlersFiber(m, g)

	p := fmt.Sprintf("%s:%d", "", port)

	// Start server
	log.Println("Start listening on:", p)
	if err := app.Listen(p); err != nil {
		panic(err)
	}
}
