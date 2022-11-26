package gotv

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Query Query for request
type Query struct {
	Tick     int    `query:"tick" form:"tick"`          // the starting tick of the broadcast
	EndTick  int    `query:"endtick" query:"endtick"`   // endtick of delta frame
	Final    bool   `query:"final" query:"final"`       // is final fragment
	TPS      int    `query:"tps" query:"tps"`           // the tickrate of the GOTV broadcast
	Map      string `query:"map" query:"map"`           // the name of the map
	Protocol int    `query:"protocol" query:"protocol"` // Currently 4
}

// CheckAuthMiddlewareFiber Check Auth on Fiber
func CheckAuthMiddlewareFiber(g Store) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		hs := c.GetReqHeaders()
		auth, ok := hs[authHeader]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("tv_broadcast_origin_auth required")
		}
		token := utils.CopyString(c.Params("token"))
		if err := g.Auth(token, auth); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		return c.Next()
	})
}

// OnStartFragmentHandlerFiber Register start fragment on Fiber
func OnStartFragmentHandlerFiber(g Store) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		q := Query{}
		if err := c.QueryParser(&q); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("BadRequest:" + err.Error())
		}
		return g.OnStart(token, StartFrame{
			At:       time.Now(),
			Tps:      q.TPS,
			Protocol: q.Protocol,
			Map:      q.Map,
			Body:     utils.CopyBytes(c.Body()),
		})
	})
}

// OnFullFragmentHandlerFiber Register start fragment on Fiber
func OnFullFragmentHandlerFiber(g Store) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		q := Query{}
		if err := c.QueryParser(&q); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("BadRequest:" + err.Error())
		}
		return g.OnFull(token, FullFrame{
			At:   time.Now(),
			Tick: q.Tick,
			Body: utils.CopyBytes(c.Body()),
		})
	})
}

// OnDeltaFragmentHandlerFiber Register start fragment on Fiber
func OnDeltaFragmentHandlerFiber(g Store) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		q := Query{}
		if err := c.QueryParser(&q); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("BadRequest:" + err.Error())
		}
		return g.OnDelta(token, DeltaFrame{
			At:      time.Now(),
			Final:   q.Final,
			EndTick: q.EndTick,
			Body:    utils.CopyBytes(c.Body()),
		})
	})
}

// SetupStoreHandlers setup Store handlers to specified fiber.Router
func SetupStoreHandlers(g Store, r fiber.Router) {
	r.Post("/:token/:fragment_number/start", CheckAuthMiddlewareFiber(g), OnStartFragmentHandlerFiber(g))
	r.Post("/:token/:fragment_number/full", CheckAuthMiddlewareFiber(g), OnFullFragmentHandlerFiber(g))
	r.Post("/:token/:fragment_number/delta", CheckAuthMiddlewareFiber(g), OnDeltaFragmentHandlerFiber(g))
}
