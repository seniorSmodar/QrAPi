package main

import (
	"module/configs"
	"module/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data":"pivas"})
	})

	routes.UserRoutes(app)

	routes.AuthRouter(app)

	
	app.Listen(configs.EnvPort())
}