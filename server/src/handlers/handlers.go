package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GET  /match/:token/sync
// GET  /match/:token/:fragment_number/:frametype
// POST /:token/:fragment_number/:frametype

// SyncHandler handlers request against /match/:token/sync
func SyncHandler(c *gin.Context) {
	if c.Params.ByName("fragment_number") != "sync" { // Rejects all requests other than /sync
		c.String(http.StatusBadRequest, "Unknown Request")
		return
	}

	t := c.Params.ByName("token")
	m, err := Matches.GetMatchByToken(t)
	if err != nil {
		log.Printf("ERR : %v\n", err)
		c.String(http.StatusNotFound, "NotFound")
		return
	}
	if c.Query("fragment") != "" {
		frag, err := strconv.Atoi(c.Query("fragment"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment should be int")
			return
		}
		full, err := m.GetFullFrame(m.Fragment)
		specifiedfull, err := m.GetFullFrame(uint32(frag))
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, "NotFound")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"tick":            specifiedfull.Tick,
			"rtdelay":         time.Since(specifiedfull.At).Seconds(),
			"rcvage":          time.Since(full.At).Seconds(),
			"fragment":        frag,
			"signup_fragment": m.SignupFragment,
			"tps":             m.Tps,
			"protocol":        m.Protocol,
		})
	} else {
		full, err := m.GetFullFrame(m.Fragment)
		delayedfull, err := m.GetFullFrame(m.Fragment - Matches.Delay)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, "NotFound")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"tick":            full.Tick,
			"rtdelay":         time.Since(delayedfull.At).Seconds(),
			"rcvage":          time.Since(full.At).Seconds(),
			"fragment":        m.Fragment - Matches.Delay,
			"signup_fragment": m.SignupFragment,
			"tps":             m.Tps,
			"protocol":        m.Protocol,
		})
	}
}

// SyncByIDHandler handlers request against /match/:token/sync by ID
func SyncByIDHandler(c *gin.Context) {
	if c.Params.ByName("fragment_number") != "sync" { // Rejects all requests other than /sync
		c.String(http.StatusBadRequest, "Unknown Request")
		return
	}
	id := c.Params.ByName("id")
	m, err := Matches.GetMatchByID(id)
	if err != nil {
		log.Printf("ERR : %v\n", err)
		c.String(http.StatusNotFound, "NotFound")
		return
	}
	if c.Query("fragment") != "" {
		frag, err := strconv.Atoi(c.Query("fragment"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment should be int")
			return
		}
		full, err := m.GetFullFrame(m.Fragment)
		specifiedfull, err := m.GetFullFrame(uint32(frag))
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, "NotFound")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"tick":            specifiedfull.Tick,
			"rtdelay":         time.Since(specifiedfull.At).Seconds(),
			"rcvage":          time.Since(full.At).Seconds(),
			"fragment":        frag,
			"signup_fragment": m.SignupFragment,
			"tps":             m.Tps,
			"protocol":        m.Protocol,
		})
	} else {
		full, err := m.GetFullFrame(m.Fragment)
		delayedfull, err := m.GetFullFrame(m.Fragment - Matches.Delay)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, "NotFound")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"tick":            full.Tick,
			"rtdelay":         time.Since(delayedfull.At).Seconds(),
			"rcvage":          time.Since(full.At).Seconds(),
			"fragment":        m.Fragment - Matches.Delay,
			"signup_fragment": m.SignupFragment,
			"tps":             m.Tps,
			"protocol":        m.Protocol,
		})
	}
}

// GetBodyHandler handles fragment request from CS:GO client
func GetBodyHandler(c *gin.Context) {
	t := c.Params.ByName("token")
	f := c.Params.ByName("fragment_number")

	fragment, err := strconv.Atoi(f)
	if err != nil {
		c.String(http.StatusBadRequest, "Fragment should be int")
		return
	}
	ft := c.Params.ByName("frametype")

	m, err := Matches.GetMatchByToken(t)
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
}

// GetBodyByIDHandler handles fragment request from CS:GO client
func GetBodyByIDHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	f := c.Params.ByName("fragment_number")

	fragment, err := strconv.Atoi(f)
	if err != nil {
		c.String(http.StatusBadRequest, "Fragment should be int")
		return
	}
	ft := c.Params.ByName("frametype")

	m, err := Matches.GetMatchByID(id)
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
		tick, err := strconv.Atoi(c.Query("tick"))
		tps, err := strconv.ParseFloat(c.Query("tps"), 10)
		protocol, err := strconv.Atoi(c.Query("protocol"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment,tps,protocol should be float or int")
			return
		}
		log.Printf("Received START Fragment. Register match... Token[%s] Tps[%f] Protocol[%d]\n", t, tps, protocol)
		Matches.Register(&Match{
			ID:             id,
			Token:          t,
			Startframe:     make(map[uint32]*Startframe),
			Fullframes:     make(map[uint32]*Fullframe),
			Deltaframes:    make(map[uint32]*Deltaframes),
			SignupFragment: uint32(fragment),
			Tps:            tps,
			Map:            c.Query("map"),
			Protocol:       uint8(protocol),
			Auth:           auth,
			Tick:           uint32(tick),
			// RtDelay:        10, // TODO?
			// RcVage:         10, // TODO?
			// Fragment:       uint32(fragment),
		})
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, "")
			return
		}
		m.Lock()
		defer m.Unlock()
		m.Startframe[uint32(fragment)] = &Startframe{
			At:   time.Now(),
			Body: reqBody,
		}
		c.String(http.StatusOK, "OK")
	case "full":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
			return
		}
		m.Lock()
		defer m.Unlock()
		tick, err := strconv.Atoi(c.Query("tick"))
		if err != nil {
			c.String(http.StatusBadRequest, "tick should be float or int")
			return
		}
		log.Printf("tick = %d\n", tick)

		m.Fragment = uint32(fragment)
		// m.Tick = uint32(tick)
		m.Fullframes[uint32(fragment)] = &Fullframe{
			At:   time.Now(),
			Tick: tick,
			Body: reqBody,
		}
		c.String(http.StatusOK, "OK")
	case "delta":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
			return
		}
		m.Lock()
		defer m.Unlock()
		endtick, err := strconv.Atoi(c.Query("endtick"))
		if err != nil {
			c.String(http.StatusBadRequest, "endtick should be float or int")
			return
		}
		log.Printf("endtick = %d\n", endtick)
		log.Printf("final = %s\n", c.Query("final"))

		// final...?

		m.Fragment = uint32(fragment)
		m.Deltaframes[uint32(fragment)] = &Deltaframes{
			Body: reqBody,
		}
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
		tick, err := strconv.Atoi(c.Query("tick"))
		tps, err := strconv.ParseFloat(c.Query("tps"), 10)
		protocol, err := strconv.Atoi(c.Query("protocol"))
		if err != nil {
			c.String(http.StatusBadRequest, "fragment,tps,protocol should be float or int")
			return
		}
		Matches.Register(&Match{
			Token:          t,
			Startframe:     make(map[uint32]*Startframe),
			Fullframes:     make(map[uint32]*Fullframe),
			Deltaframes:    make(map[uint32]*Deltaframes),
			SignupFragment: uint32(fragment),
			Tps:            tps,
			Map:            c.Query("map"),
			Protocol:       uint8(protocol),
			Auth:           auth,
			Tick:           uint32(tick),
			// RtDelay:        10, // TODO?
			// RcVage:         10, // TODO?
			// Fragment:       uint32(fragment),
		})
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusNotFound, "")
			return
		}
		m.Lock()
		defer m.Unlock()
		m.Startframe[uint32(fragment)] = &Startframe{
			At:   time.Now(),
			Body: reqBody,
		}
		c.String(http.StatusOK, "OK")
	case "full":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
			return
		}
		m.Lock()
		defer m.Unlock()
		tick, err := strconv.Atoi(c.Query("tick"))
		if err != nil {
			c.String(http.StatusBadRequest, "tick should be float or int")
			return
		}
		log.Printf("tick = %d\n", tick)

		m.Fragment = uint32(fragment)
		// m.Tick = uint32(tick)
		m.Fullframes[uint32(fragment)] = &Fullframe{
			At:   time.Now(),
			Tick: tick,
			Body: reqBody,
		}
		c.String(http.StatusOK, "OK")
	case "delta":
		m, err := Matches.GetMatchByToken(t)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.String(http.StatusResetContent, "")
			return
		}
		m.Lock()
		defer m.Unlock()
		endtick, err := strconv.Atoi(c.Query("endtick"))
		if err != nil {
			c.String(http.StatusBadRequest, "endtick should be float or int")
			return
		}
		log.Printf("endtick = %d\n", endtick)
		log.Printf("final = %s\n", c.Query("final"))

		// final...?

		m.Fragment = uint32(fragment)
		m.Deltaframes[uint32(fragment)] = &Deltaframes{
			Body: reqBody,
		}
		c.String(http.StatusOK, "OK")

	default:
		log.Println("frametype : unknown...")
		c.String(http.StatusBadRequest, "Unknown")
	}
}

// GetListHandler handles list of Matches
func GetListHandler(c *gin.Context) {
	m, err := Matches.GetTokens()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, m)
}
