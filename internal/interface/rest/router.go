package rest

import (
	"strings"
	"sync"
	"teng231/goapp/internal/app"
	"teng231/goapp/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

type API struct {
	router     *fiber.App
	wsConn     map[string]*websocket.Conn
	mt         sync.Mutex
	liveHubApp *app.LiveCommentApp
}

func New(hub *app.LiveCommentApp) *API {
	api := &API{mt: sync.Mutex{}, liveHubApp: hub, wsConn: make(map[string]*websocket.Conn)}
	api.router = api.createRouter()
	return api
}

func (a *API) Router() *fiber.App {
	return a.router
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	// Lowercase để so sánh không phân biệt hoa/thường
	o := strings.ToLower(origin)
	// Cho phép nếu chứa "localhost"
	if strings.Contains(o, "localhost") {
		return true
	}
	// Cho phép nếu chứa "tikcart"
	if strings.Contains(o, "tikcart") {
		return true
	}
	return false
}

func (a *API) createRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Network:      fiber.NetworkTCP,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 3 * time.Minute,
	})
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: isAllowedOrigin,
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Range",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,HEAD", // Allowed methods
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))
	// Middlewares
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(logger.New())
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	apiRouter := app.Group("/api")
	apiRouter.Get("/ping", handleSimpleHealthcheck)
	apiRouter.Post("/cart", a.handleAddCartItem)
	wsRouter := app.Group("/ws")
	// Middleware để nâng cấp lên WS
	wsRouter.Use("/conn/:roomId", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	wsRouter.Get("/conn/:roomId", websocket.New(a.wsConnHandler))

	////////////////////////
	// UNPROTECTED ROUTES //
	////////////////////////

	// unprotectedAPIRouter := app.Group("/")
	// unprotectedAPIRouter.Get("/v1/ping", handleSimpleHealthcheck)

	//////////////////////
	// PROTECTED ROUTES //
	//////////////////////

	// protectedAPIRouter := app.Group("/api/")
	// protectedAPIRouter.Get("/stack", a.handleAppStack)
	return app
}

func handleSimpleHealthcheck(c *fiber.Ctx) error {
	return c.SendString("pong")
}
func (a *API) handleAddCartItem(c *fiber.Ctx) error {
	var comment *domain.Comment
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if comment.Username == "" || comment.Comment == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username and message are required"})
	}
	if err := a.liveHubApp.AppendSheet(c.Context(), comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to append comment"})
	}
	return c.SendString("pong")

}
