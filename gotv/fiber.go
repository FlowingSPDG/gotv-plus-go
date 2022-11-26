package gotv

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/xerrors"
)

// CheckAuthMiddlewareFiber Check Auth on Fiber
func CheckAuthMiddlewareFiber(g Store) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		hs := c.GetReqHeaders()
		auth, ok := hs["X-Origin-Auth"]
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
		fragment, err := strconv.Atoi(c.Params("fragment_number"))
		if err != nil {
			return err
		}
		q := StartQuery{}
		if err := c.QueryParser(&q); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("BadRequest:" + err.Error())
		}
		if err := g.OnStart(token, fragment, StartFrame{
			At:       time.Now(),
			Tps:      q.TPS,
			Protocol: q.Protocol,
			Map:      q.Map,
			Body:     utils.CopyBytes(c.Body()),
		}); err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				return c.Status(fiber.StatusResetContent).SendString("RESET CONTENT")
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("FRAGMENT NOT FOUND")
			}
			return err
		}
		return nil
	})
}

// OnFullFragmentHandlerFiber Register start fragment on Fiber
func OnFullFragmentHandlerFiber(g Store) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		fragment, err := strconv.Atoi(c.Params("fragment_number"))
		if err != nil {
			return err
		}
		q := FullQuery{}
		if err := c.QueryParser(&q); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("BadRequest:" + err.Error())
		}
		if err := g.OnFull(token, fragment, q.Tick, time.Now(), utils.CopyBytes(c.Body())); err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				return c.Status(fiber.StatusResetContent).SendString("RESET CONTENT")
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("FRAGMENT NOT FOUND")
			}
			return err
		}
		return nil
	})
}

// OnDeltaFragmentHandlerFiber Register start fragment on Fiber
func OnDeltaFragmentHandlerFiber(g Store) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		fragment, err := strconv.Atoi(c.Params("fragment_number"))
		if err != nil {
			return err
		}
		q := DeltaQuery{}
		if err := c.QueryParser(&q); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("BadRequest:" + err.Error())
		}
		if err := g.OnDelta(token, fragment, q.EndTick, time.Now(), q.Final, utils.CopyBytes(c.Body())); err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				return c.Status(fiber.StatusResetContent).SendString("RESET CONTENT")
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("FRAGMENT NOT FOUND")
			}
			return err
		}
		return nil
	})
}

// GetSyncRequestHandlerFiber Register start fragment on Fiber
func GetSyncRequestHandlerFiber(b Broadcaster) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		q := SyncQuery{}
		if err := c.QueryParser(&q); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("BadRequest:" + err.Error())
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
				return c.Status(fiber.StatusNotFound).SendString("MATCH NOT FOUND")
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("FRAGMENT NOT FOUND")
			}
			return err
		}
		return c.JSON(s)
	})
}

// GetStartRequestHandlerFiber Get start fragment on Fiber
func GetStartRequestHandlerFiber(b Broadcaster) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		fragment, err := strconv.Atoi(c.Params("fragment_number"))
		if err != nil {
			return err
		}
		s, err := b.GetStart(token, fragment)
		if err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("MATCH NOT FOUND")
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("FRAGMENT NOT FOUND")
			}
			return err
		}
		return c.Status(fiber.StatusOK).Send(s)
	})
}

// GetFullRequestHandlerFiber Get start fragment on Fiber
func GetFullRequestHandlerFiber(b Broadcaster) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		fragment, err := strconv.Atoi(c.Params("fragment_number"))
		if err != nil {
			return err
		}
		s, err := b.GetFull(token, fragment)
		if err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("MATCH NOT FOUND")
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("FRAGMENT NOT FOUND")
			}
			return err
		}
		return c.Status(fiber.StatusOK).Send(s)
	})
}

// GetDeltaRequestHandlerFiber Get delta fragment on Fiber
func GetDeltaRequestHandlerFiber(b Broadcaster) func(c *fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		token := utils.CopyString(c.Params("token"))
		fragment, err := strconv.Atoi(c.Params("fragment_number"))
		if err != nil {
			return err
		}
		d, err := b.GetDelta(token, fragment)
		if err != nil {
			if xerrors.Is(err, ErrMatchNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("MATCH NOT FOUND")
			}
			if xerrors.Is(err, ErrFragmentNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("FRAGMENT NOT FOUND")
			}
			return err
		}
		return c.Status(fiber.StatusOK).Send(d)
	})
}

// SetupStoreHandlers setup Store handlers to specified fiber.Router
func SetupStoreHandlersFiber(g Store, r fiber.Router) {
	r.Post("/:token/:fragment_number/start", CheckAuthMiddlewareFiber(g), OnStartFragmentHandlerFiber(g))
	r.Post("/:token/:fragment_number/full", CheckAuthMiddlewareFiber(g), OnFullFragmentHandlerFiber(g))
	r.Post("/:token/:fragment_number/delta", CheckAuthMiddlewareFiber(g), OnDeltaFragmentHandlerFiber(g))
}

// SetupBroadcasterHandlers setup Broadcaster handlers to specified fiber.Router
func SetupBroadcasterHandlersFiber(b Broadcaster, r fiber.Router) {
	r.Get("/:token/sync", GetSyncRequestHandlerFiber(b))
	r.Get("/:token/:fragment_number/start", GetStartRequestHandlerFiber(b))
	r.Get("/:token/:fragment_number/full", GetFullRequestHandlerFiber(b))
	r.Get("/:token/:fragment_number/delta", GetDeltaRequestHandlerFiber(b))
}
