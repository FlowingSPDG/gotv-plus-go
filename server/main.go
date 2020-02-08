package main

import (
	"github.com/FlowingSPDG/gotv-plus-go/server/src"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	Matches := &models.Matches{
		Matches: make(map[string]*models.Match),
		Auth:    "gopher",
	}

	r := gin.Default()

	// GET  /match/:token/sync
	// GET  /match/:token/:fragment_number/:frametype
	// POST /:token/:fragment_number/:frametype

	// SYNC MATCH

	r.GET("/", func(c *gin.Context) {
		// index...
	})

	// SYNC!!
	r.GET("/match/:token/:fragment_number", func(c *gin.Context) {
		if c.Params.ByName("fragment_number") != "sync" { // Rejects all requests other than /sync
			c.String(http.StatusBadRequest, "Unknown Request")
			return
		}
		t := c.Params.ByName("token")
		m, err := Matches.GetMatch(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, "NotFound")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"tick":            m.Tick,
			"rtdelay":         m.RtDelay,
			"rcvage":          m.RcVage,
			"fragment":        m.Fragment,
			"signup_fragment": m.SignupFragment,
			"tps":             m.Tps,
			"protocol":        m.Protocol,
		})
	})

	r.GET("/match/:token/:fragment_number/:frametype", func(c *gin.Context) {
		t := c.Params.ByName("token")
		f := c.Params.ByName("fragment_number")

		fragment, err := strconv.Atoi(f)
		if err != nil {
			c.String(http.StatusBadRequest, "Fragment should be int")
			return
		}
		ft := c.Params.ByName("frametype")

		m, err := Matches.GetMatch(t)
		if err != nil {
			c.String(http.StatusNotFound, "Match not found")
			return
		}
		frags, err := m.GetBody(ft, uint32(fragment))
		if err != nil {
			c.String(http.StatusNotFound, "Fragment not found")
			return
		}
		c.Data(200, "application/octet-stream", frags)
	})

	r.POST("/:token/:fragment_number/:frametype", func(c *gin.Context) {
		t := c.Params.ByName("token")
		f := c.Params.ByName("fragment_number")
		fragment, err := strconv.Atoi(f)
		if err != nil {
			c.String(http.StatusBadRequest, "Fragment should be int")
			return
		}
		ft := c.Params.ByName("frametype")
		auth := c.Request.Header.Get("x-origin-auth")
		// log.Printf("token : [%s], fragment_number:[%s], frametype=[%s] auth=[%s]\n", t, f, ft, auth)
		// log.Printf("Queries : %v\n", c.Request.URL.Query())

		if auth != Matches.Auth {
			c.String(http.StatusForbidden, "Auth not match")
			return
		}

		reqBody, err := c.GetRawData() // body
		if err != nil {
			c.String(http.StatusForbidden, "Failed to fetch request body")
			return
		}

		switch ft {
		case "start":
			tick, err := strconv.Atoi(c.Query("tick"))
			tps, err := strconv.ParseFloat(c.Query("tps"), 10)
			protocol, err := strconv.Atoi(c.Query("protocol"))
			if err != nil {
				c.String(http.StatusBadRequest, "fragment,tps,protocol should be float or int")
				return
			}
			/*
				&models.Startframe{
						At:   time.Now(),
						Body: reqBody,
					},
			*/
			Matches.Register(&models.Match{
				Token:          t,
				Startframe:     make(map[uint32]*models.Startframe),
				Fullframes:     make(map[uint32]*models.Fullframe),
				Deltaframes:    make(map[uint32]*models.Deltaframes),
				SignupFragment: uint32(fragment),
				Tps:            tps,
				Map:            c.Query("map"),
				Protocol:       uint8(protocol),
				Auth:           auth,
				Tick:           uint32(tick),
				RtDelay:        2, // ?
				RcVage:         2, // ?
				// Fragment:       uint32(fragment),
			})
			m, err := Matches.GetMatch(t)
			m.Startframe[uint32(fragment)] = &models.Startframe{
				At:   time.Now(),
				Body: reqBody,
			}
			c.String(http.StatusOK, "OK")
		case "full":
			m, err := Matches.GetMatch(t)
			if err != nil {
				log.Printf("ERR : %v\n", err)
				c.String(http.StatusResetContent, "RESET")
				return
			}
			tick, err := strconv.Atoi(c.Query("tick"))
			if err != nil {
				c.String(http.StatusBadRequest, "tick should be float or int")
				return
			}
			log.Printf("tick = %d\n", tick)

			m.Tick = uint32(tick)
			m.Fragment = uint32(fragment)
			m.Fullframes[uint32(fragment)] = &models.Fullframe{
				At:   time.Now(),
				Tick: tick,
				Body: reqBody,
			}
			c.String(http.StatusOK, "OK")
		case "delta":
			m, err := Matches.GetMatch(t)
			if err != nil {
				log.Printf("ERR : %v\n", err)
				c.String(http.StatusResetContent, "RESET")
				return
			}
			endtick, err := strconv.Atoi(c.Query("endtick"))
			if err != nil {
				c.String(http.StatusBadRequest, "endtick should be float or int")
				return
			}
			log.Printf("endtick = %d\n", endtick)
			log.Printf("final = %s\n", c.Query("final"))

			// final...?

			m.Fragment = uint32(fragment)
			m.Deltaframes[uint32(fragment)] = &models.Deltaframes{
				Body: reqBody,
			}
			c.String(http.StatusOK, "OK")

		default:
			log.Println("frametype : unknown...")
			c.String(http.StatusBadRequest, "Unknown")
		}
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
