package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jos3lo89/vehicle-control-unajma-api/config"
	"github.com/jos3lo89/vehicle-control-unajma-api/config/database"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/controllers"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/repositories"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/routes"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/services"
)

func main() {

	// envs
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// base de datos
	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("No se pudo iniciar la BD: %v", err)
	}
	defer pool.Close()

	// Inyección de Dependencias
	authRepo := repositories.NewAuthRepository(pool)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	// app
	app := fiber.New(fiber.Config{
		AppName:      "Vehicle Control API v1.0.0",
		ServerHeader: "Fiber",
	})

	app.Use(logger.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"env":           cfg.Env,
			"app":           "vehicle control unajma api",
			"Base de datos": "conectado",
		})
	})

	routes.SetupRoutes(app, authController)

	// iniciar app
	log.Println("Servidor corriendo en puerto:", cfg.Port)
	app.Listen(":" + cfg.Port)
}
