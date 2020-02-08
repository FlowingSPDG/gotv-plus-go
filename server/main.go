package main

import (
	"github.com/FlowingSPDG/gotv-plus-go/server/src/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// Listen and Server in localhost:8080
	addr := "localhost:8080"

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("/", func(c *gin.Context) {
		m, _ := handlers.Matches.GetTokens()
		log.Printf("Matches : %v\n", m)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Title":   "GOTV+ for Gophers",
			"Matches": m,
			"Addr":    addr,
		})
	})

	r.GET("/match/:token/:fragment_number", handlers.SyncHandler)
	r.GET("/match/:token/:fragment_number/:frametype", handlers.GetBodyHandler)
	r.POST("/:token/:fragment_number/:frametype", handlers.PostBodyHandler)

	r.GET("/matches", handlers.GetListHandler)

	log.Panicf("Failed to listen port %s : %v\n", addr, r.Run(addr))
}
