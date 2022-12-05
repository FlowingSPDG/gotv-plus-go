package gotv

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

// CheckAuthMiddlewareGin Check Auth on Gin
func CheckAuthMiddlewareGin(g Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		auth := c.Request.Header.Get("X-Origin-Auth")
		if auth == "" {
			c.String(http.StatusUnauthorized, "tv_broadcast_origin_auth required")
			c.Abort()
			return
		}
		if err := g.Auth(token, auth); err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

// OnStartFragmentHandlerGin Register start fragment on Gin
func OnStartFragmentHandlerGin(g Store) func(c *gin.Context) {
	return (func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		q := StartQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		if err := g.OnStart(token, fragment, StartFrame{
			At:       time.Now(),
			Tps:      q.TPS,
			Protocol: q.Protocol,
			Map:      q.Map,
			Body:     b,
		}); err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusResetContent, "RESET CONTENT")
				c.Abort()
				return
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		return
	})
}

// OnFullFragmentHandlerGin Register start fragment on Gin
func OnFullFragmentHandlerGin(g Store) func(c *gin.Context) {
	return (func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		q := FullQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		if err := g.OnFull(token, fragment, q.Tick, time.Now(), b); err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusResetContent, "RESET CONTENT")
				c.Abort()
				return
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		return
	})
}

// OnDeltaFragmentHandlerGin Register start fragment on Gin
func OnDeltaFragmentHandlerGin(g Store) func(c *gin.Context) {
	return (func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		q := DeltaQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		if err := g.OnDelta(token, fragment, q.EndTick, time.Now(), q.Final, b); err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusResetContent, "RESET CONTENT")
				c.Abort()
				return
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		return
	})
}

// GetSyncRequestHandlerGin get sync JSON on Gin
func GetSyncRequestHandlerGin(b Broadcaster) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		q := SyncQuery{}
		if err := c.BindQuery(&q); err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		var s Sync
		var err error
		if q.Fragment != 0 {
			s, err = b.GetSync(token, q.Fragment)
		} else {
			s, err = b.GetSyncLatest(token)
		}
		if err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
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
func GetStartRequestHandlerGin(b Broadcaster) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		b, err := b.GetStart(token, fragment)
		if err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		c.Data(http.StatusOK, "application/octet-stream", b)
		return
	}
}

// GetFullRequestHandlerGin get full on Gin
func GetFullRequestHandlerGin(b Broadcaster) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		b, err := b.GetFull(token, fragment)
		if err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		c.Data(http.StatusOK, "application/octet-stream", b)
		return
	}
}

// GetDeltaRequestHandlerGin get full on Gin
func GetDeltaRequestHandlerGin(b Broadcaster) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		fragment, err := strconv.Atoi(c.Param("fragment_number"))
		if err != nil {
			c.String(http.StatusBadRequest, "BadRequest:"+err.Error())
			c.Abort()
			return
		}
		b, err := b.GetDelta(token, fragment)
		if err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				c.String(http.StatusNotFound, "MATCH NOT FOUND")
				c.Abort()
				return
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				c.String(http.StatusNotFound, "FRAGMENT NOT FOUND")
				c.Abort()
				return
			}
			return
		}
		c.Data(http.StatusOK, "application/octet-stream", b)
		return
	}
}

// SetupStoreHandlersGin setup Store handlers to gin.RouterGroup
func SetupStoreHandlersGin(g Store, r *gin.RouterGroup) {
	r.POST("/:token/:fragment_number/start", CheckAuthMiddlewareGin(g), OnStartFragmentHandlerGin(g))
	r.POST("/:token/:fragment_number/full", CheckAuthMiddlewareGin(g), OnFullFragmentHandlerGin(g))
	r.POST("/:token/:fragment_number/delta", CheckAuthMiddlewareGin(g), OnDeltaFragmentHandlerGin(g))
}

// SetupBroadcasterHandlersGin setup Broadcaster handlers to specified gin.RouterGroup
func SetupBroadcasterHandlersGin(b Broadcaster, r *gin.RouterGroup) {
	r.GET("/:token/sync", GetSyncRequestHandlerGin(b))
	r.GET("/:token/:fragment_number/start", GetStartRequestHandlerGin(b))
	r.GET("/:token/:fragment_number/full", GetFullRequestHandlerGin(b))
	r.GET("/:token/:fragment_number/delta", GetDeltaRequestHandlerGin(b))
}
