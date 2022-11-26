package main

import (
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/FlowingSPDG/gotv-plus-go/examples/disk"
	"github.com/FlowingSPDG/gotv-plus-go/gotv"
)

func main() {
	m := disk.NewDiskGOTV("SuperSecureStringDoNotShare", "gotv_plus_binary")
	app := fiber.New()
	g := app.Group("/gotv") // /gotv
	g.Use(logger.New())
	gotv.SetupStoreHandlersFiber(m, g)
	gotv.SetupBroadcasterHandlersFiber(m, g)

	p := net.JoinHostPort("localhost", "8080")

	// Start server
	log.Println("Start listening on:", p)
	if err := app.Listen(p); err != nil {
		panic(err)
	}
}
