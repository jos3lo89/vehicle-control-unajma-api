package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jos3lo89/vehicle-control-unajma-api/config"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Vehicle Control API v1.0.0",
		ServerHeader: "Fiber",
	})

	app.Use(logger.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"app": "vehicle control unajma api",
			"env": cfg.Env,
		})
	})

	log.Println("Servidor corriendo en puerto:", cfg.Port)
	app.Listen(":" + cfg.Port)
}
