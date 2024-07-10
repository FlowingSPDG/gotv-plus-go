package gotv

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

type GinCSTV struct {
	a Auth
	s Store
	b Broadcaster
}

func NewGinCSTV(a Auth, s Store, b Broadcaster) *GinCSTV {
	return &GinCSTV{a: a, s: s, b: b}
}

// CheckAuthMiddlewareGin Check Auth
func (g *GinCSTV) CheckAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		auth := c.Request.Header.Get("X-Origin-Auth")
		if err := g.a.Auth(token, auth); err != nil {
			c.AbortWithError(http.StatusUnauthorized, xerrors.Errorf("Unauthorized: %w", err))
			return
		}
		c.Next()
	}
}

// OnStartFragmentHandlerGin Register start fragment on Gin
func (g *GinCSTV) OnStartFragment() func(c *gin.Context) {
	return (func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		q := StartQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err := g.s.OnFull(token, fragment, q.Tick, time.Now(), c.Request.Body); err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.AbortWithStatus(http.StatusResetContent)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})
}

// OnFullFragmentHandlerGin Register start fragment on Gin
func (g *GinCSTV) OnFullFragment() func(c *gin.Context) {
	return (func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		q := FullQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err := g.s.OnFull(token, fragment, q.Tick, time.Now(), c.Request.Body); err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.AbortWithStatus(http.StatusResetContent)
				return
			}
			if errors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
		}
	})
}

// OnDeltaFragmentHandlerGin Register start fragment on Gin
func (g *GinCSTV) OnDeltaFragment() func(c *gin.Context) {
	return (func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		q := DeltaQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		final := false
		if q.Final != nil {
			final = *q.Final
		}
		if err := g.s.OnDelta(token, fragment, q.EndTick, time.Now(), final, c.Request.Body); err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.AbortWithStatus(http.StatusResetContent)
				return
			}
			if errors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
		}
	})
}

// OnSyncRequest get sync JSON on Gin
func (g *GinCSTV) OnSyncRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		q := SyncQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var s Sync
		var err error
		if q.Fragment != nil {
			s, err = g.b.GetSync(token, *q.Fragment)
		} else {
			s, err = g.b.GetSyncLatest(token)
		}
		if err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if errors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		c.JSON(http.StatusOK, s)
		return
	}
}

// GetStartRequestHandlerGin get start on Gin
func (g *GinCSTV) OnGetStartRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		r, err := g.b.GetStart(token, fragment)
		if err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if errors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		defer r.Close()
		defer c.Writer.Flush()
		if _, err := io.Copy(c.Writer, r); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		return
	}
}

// GetFullRequestHandlerGin get full on Gin
func (g *GinCSTV) OnGetFullRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		r, err := g.b.GetFull(token, fragment)
		if err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if errors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		defer r.Close()
		defer c.Writer.Flush()
		if _, err := io.Copy(c.Writer, r); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		return
	}
}

// GetDeltaRequestHandlerGin get full on Gin
func (g *GinCSTV) OnGetDeltaRequest(b Broadcaster) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		r, err := g.b.GetDelta(token, fragment)
		if err != nil {
			if errors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if errors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		defer r.Close()
		defer c.Writer.Flush()
		if _, err := io.Copy(c.Writer, r); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		return
	}
}

// SetupStoreHandlersGin setup Store handlers to gin.RouterGroup
func SetupStoreHandlersGin(g *GinCSTV, r *gin.RouterGroup) {
	r.Use(g.CheckAuthMiddleware())
	r.POST("/:token/:fragment_number/start", g.OnStartFragment())
	r.POST("/:token/:fragment_number/full", g.OnFullFragment())
	r.POST("/:token/:fragment_number/delta", g.OnDeltaFragment())
}

// SetupBroadcasterHandlersGin setup Broadcaster handlers to specified gin.RouterGroup
func SetupBroadcasterHandlersGin(g *GinCSTV, r *gin.RouterGroup) {
	r.Use(g.CheckAuthMiddleware())
	r.GET("/:token/sync", g.OnSyncRequest())
	r.GET("/:token/:fragment_number/start", g.OnGetStartRequest())
	r.GET("/:token/:fragment_number/full", g.OnGetFullRequest())
	r.GET("/:token/:fragment_number/delta", g.OnDeltaFragment())
}
