package main

import (
	"flag"
	grpc "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc"
	"github.com/FlowingSPDG/gotv-plus-go/server/src/handlers"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	addr     = flag.String("addr", "localhost:8080", "Address where GOTV+ hosted at")
	debug    = flag.Bool("debug", false, "Debug mode option")
	grpcaddr = flag.String("grpc", "localhost:50055", "gRPC API Address")
	delay    = flag.Int("delay", 3, "How much frags to delay.")
	auth     = flag.String("auth", "gopher", "GOTV+ Auth password")
)

func init() {
	flag.Parse()

	log.Printf("DEBUG MODE : %v\n", *debug)
	if *debug == true {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.DefaultWriter = ioutil.Discard
	}
	handlers.InitMatchEngine(*auth, uint32(*delay))
	go grpc.StartGRPC(*grpcaddr)
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

	// playcast "http://localhost:8080/match/token/s1m1...."
	r.GET("/match/token/:token/:fragment_number", handlers.SyncHandler)
	r.GET("/match/token/:token/:fragment_number/:frametype", handlers.GetBodyHandler)

	// playcast "http://localhost:8080/match/id/YOUR_MATCH_ID"
	r.GET("/match/id/:id/:fragment_number", handlers.SyncByIDHandler)
	r.GET("/match/id/:id/:fragment_number/:frametype", handlers.GetBodyByIDHandler)

	// tv_broadcast_url "http://localhost:8080/token/"
	r.POST("/token/:token/:fragment_number/:frametype", handlers.PostBodyHandler)

	// tv_broadcast_url "http://localhost:8080/id/YOUR_MATCH_ID"
	r.POST("/id/:id/:token/:fragment_number/:frametype", handlers.PostBodyByIDHandler)

	r.GET("/matches", handlers.GetListHandler)

	log.Panicf("Failed to listen port %s : %v\n", *addr, r.Run(*addr))
}
