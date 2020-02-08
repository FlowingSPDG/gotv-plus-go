package main

import (
	"flag"
	"github.com/FlowingSPDG/gotv-plus-go/server/src/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	addr  = flag.String("addr", "localhost:8080", "Address where GOTV+ hosted at")
	debug = flag.Bool("debug", false, "Debug mode option")
	delay = flag.Int("delay", 3, "How much frags to delay.")
	auth  = flag.String("auth", "gopher", "GOTV+ Auth password")
)

func init() {
	flag.Parse()

	log.Printf("DEBUG MODE : %v\n", *debug)
	if *debug == true {
		gin.SetMode(gin.ReleaseMode)
	}
	handlers.InitMatchEngine(*auth, uint32(*delay))
}

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("/", func(c *gin.Context) {
		m, _ := handlers.Matches.GetTokens()
		log.Printf("Matches : %v\n", m)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Title":   "GOTV+ for Gophers",
			"Matches": m,
			"Addr":    *addr,
		})
	})

	r.GET("/match/:token/:fragment_number", handlers.SyncHandler)
	r.GET("/match/:token/:fragment_number/:frametype", handlers.GetBodyHandler)
	r.POST("/:token/:fragment_number/:frametype", handlers.PostBodyHandler)

	r.GET("/matches", handlers.GetListHandler)

	log.Panicf("Failed to listen port %s : %v\n", *addr, r.Run(*addr))
}
