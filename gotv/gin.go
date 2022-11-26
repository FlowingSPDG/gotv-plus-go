package gotv

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WIP: GOTV+の実装をGinで作る

// CheckAuthMiddlewareGin Check Auth on Gin
func CheckAuthMiddlewareGin(g Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		auth := c.Request.Header.Get(authHeader)
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
