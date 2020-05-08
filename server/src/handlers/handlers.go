package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GET  /match/:token/sync
// GET  /match/:token/:fragment_number/:frametype
// POST /:token/:fragment_number/:frametype

// SyncHandler handlers request against /match/:token/sync
func SyncHandler(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3")
	// Set "Expires" for 3sec...

	if c.Params.ByName("fragment_number") != "sync" { // Rejects all requests other than /sync
		c.String(http.StatusBadRequest, "Unknown Request")
		return
	}

	t := c.Params.ByName("token")
	m, err := Matches.GetMatchByToken(t)
	if err != nil {
		log.Printf("ERR : %v\n", err)
		c.String(http.StatusNotFound, err.Error())
		return
	}
	if c.Query("fragment") != "" {
		fragment, err := strconv.Atoi(c.Query("fragment"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment should be int")
			return
		}
		frag := uint32(fragment)
		json, err := m.Sync(frag)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, json)
	} else {
		json, err := m.Sync(m.Latest)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, json)
	}
}

// SyncByIDHandler handlers request against /match/:token/sync by ID
func SyncByIDHandler(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=5")
	if c.Params.ByName("fragment_number") != "sync" { // Rejects all requests other than /sync
		c.String(http.StatusBadRequest, "Unknown Request")
		return
	}
	id := c.Params.ByName("id")
	m, err := Matches.GetMatchByID(id)
	if err != nil {
		log.Printf("ERR : %v\n", err)
		c.String(http.StatusNotFound, err.Error())
		return
	}
	if c.Query("fragment") != "" {
		fragment, err := strconv.Atoi(c.Query("fragment"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment should be int")
			return
		}
		frag := uint32(fragment)
		json, err := m.Sync(frag)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, json)
	} else {
		json, err := m.Sync(m.Latest)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, json)
	}
}

// GetBodyHandler handles fragment request from CS:GO client
func GetBodyHandler(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000")
	t := c.Params.ByName("token")
	f := c.Params.ByName("fragment_number")

	fragment, err := strconv.Atoi(f)
	if err != nil {
		c.String(http.StatusBadRequest, "Fragment should be int")
		return
	}
	frag := uint32(fragment)
	ft := c.Params.ByName("frametype")

	m, err := Matches.GetMatchByToken(t)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	switch ft {
	case "full":
		full, err := m.GetFullFrame(frag)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.Data(200, "application/octet-stream", full.Body)
		return
	case "delta":
		delta, err := m.GetDeltaFrame(frag)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.Data(200, "application/octet-stream", delta.Body)
		return
	case "start":
		start, err := m.GetStartFrame(frag)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.Data(200, "application/octet-stream", start.Body)
		return
	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("INVALID FRAME TYPE"))
	}

}

// GetBodyByIDHandler handles fragment request from CS:GO client
func GetBodyByIDHandler(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000")
	id := c.Params.ByName("id")
	f := c.Params.ByName("fragment_number")

	fragment, err := strconv.Atoi(f)
	if err != nil {
		c.String(http.StatusBadRequest, "Fragment should be int")
		return
	}
	frag := uint32(fragment)
	ft := c.Params.ByName("frametype")

	m, err := Matches.GetMatchByID(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	switch ft {
	case "full":
		full, err := m.GetFullFrame(frag)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.Data(200, "application/octet-stream", full.Body)
		return
	case "delta":
		delta, err := m.GetDeltaFrame(frag)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.Data(200, "application/octet-stream", delta.Body)
		return
	case "start":
		start, err := m.GetStartFrame(frag)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.Data(200, "application/octet-stream", start.Body)
		return
	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("INVALID FRAME TYPE"))
	}
}

func PostBodyByIDHandler(c *gin.Context) {
	t := c.Params.ByName("token")
	id := c.Params.ByName("id")
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
		tpsF, err := strconv.ParseFloat(c.Query("tps"), 10)
		protocol, err := strconv.Atoi(c.Query("protocol"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment,tps,protocol should be float or int")
			return
		}
		tps := uint32(tpsF)
		log.Printf("Received START Fragment. Register match... Token[%s] Tps[%d] Protocol[%d]\n", t, tps, protocol)
		Matches.Register(&Match{
			ID:             id,
			Token:          t,
			Startframe:     make(map[uint32]*Startframe),
			Fullframes:     make(map[uint32]*Fullframe),
			Deltaframes:    make(map[uint32]*Deltaframes),
			SignupFragment: uint32(fragment),
			Tps:            uint32(tps),
			Map:            c.Query("map"),
			Protocol:       uint8(protocol),
			Auth:           auth,
		})
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, err.Error())
			return
		}
		m.RegisterStartFrame(uint32(fragment), &Startframe{At: time.Now(), Body: reqBody}, uint32(tps))

		c.String(http.StatusOK, "OK")
	case "full":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
			return
		}
		tick, err := strconv.Atoi(c.Query("tick"))
		if err != nil {
			c.String(http.StatusBadRequest, "tick should be float or int")
			return
		}
		log.Printf("tick = %d\n", tick)

		m.RegisterFullFrame(uint32(fragment), &Fullframe{
			At:   time.Now(),
			Tick: uint64(tick),
			Body: reqBody,
		})
		c.String(http.StatusOK, "OK")
	case "delta":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
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

		m.RegisterDeltaFrame(uint32(fragment), &Deltaframes{
			At:      time.Now(),
			Body:    reqBody,
			EndTick: uint64(endtick),
		})
		c.String(http.StatusOK, "OK")

	default:
		log.Println("frametype : unknown...")
		c.String(http.StatusBadRequest, "Unknown")
	}
}

// PostBodyHandler handles fragment registration from CS:GO Server
func PostBodyHandler(c *gin.Context) {
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
		tpsF, err := strconv.ParseFloat(c.Query("tps"), 10)
		protocol, err := strconv.Atoi(c.Query("protocol"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment,tps,protocol should be float or int")
			return
		}
		tps := uint32(tpsF)
		match := &Match{
			Token:          t,
			Startframe:     make(map[uint32]*Startframe),
			Fullframes:     make(map[uint32]*Fullframe),
			Deltaframes:    make(map[uint32]*Deltaframes),
			SignupFragment: uint32(fragment),
			Tps:            uint32(tps),
			Map:            c.Query("map"),
			Protocol:       uint8(protocol),
			Auth:           auth,
		}
		log.Printf("Received START Fragment. Register match... Token[%s] Tps[%d] Protocol[%d]\n", t, tps, protocol)
		match.RegisterStartFrame(uint32(fragment), &Startframe{At: time.Now(), Body: reqBody}, uint32(tps))
		Matches.Register(match)
		c.String(http.StatusOK, "OK")
	case "full":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
			return
		}
		tick, err := strconv.Atoi(c.Query("tick"))
		if err != nil {
			c.String(http.StatusBadRequest, "tick should be float or int")
			return
		}
		log.Printf("tick = %d\n", tick)

		m.RegisterFullFrame(uint32(fragment), &Fullframe{
			At:   time.Now(),
			Tick: uint64(tick),
			Body: reqBody,
		})
		c.String(http.StatusOK, "OK")
	case "delta":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
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

		m.RegisterDeltaFrame(uint32(fragment), &Deltaframes{
			At:      time.Now(),
			Body:    reqBody,
			EndTick: uint64(endtick),
		})
		c.String(http.StatusOK, "OK")

	default:
		log.Println("frametype : unknown...")
		c.String(http.StatusBadRequest, "Unknown")
	}
}

// GetListHandler handles list of Matches
func GetListHandler(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=10")
	m, err := Matches.GetTokens()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, m)
}
