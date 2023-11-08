package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Welcome to golang rest with biber and mongodb"})
	})

	//run database
	configs.ConnectMongoDB()

	routes.EmployeeRoute(app)

	app.Listen(":8081")
}
