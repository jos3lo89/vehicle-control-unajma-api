package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/controllers"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/entities"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/middlewares"
)

func SetupRoutes(app *fiber.App, authCtrl *controllers.AuthController) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/login", authCtrl.Login)

	protected := v1.Group("/", middlewares.RequireAuth())

	admin := protected.Group("/admin", middlewares.RequireRole(entities.RolAdministrador))

	admin.Post("/usuarios", authCtrl.CrearUsuario)

	reportes := protected.Group("/reportes", middlewares.RequireRole(entities.RolAdministrador, entities.RolConsulta))
	reportes.Get("/ver", func(c fiber.Ctx) error {
		return c.SendString("Reportes vistos exitosamente")
	})
}
