package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/services"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// ==========================================
// LOGIN
// ==========================================
func (ctrl *AuthController) Login(c fiber.Ctx) error {
	type loginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var input loginInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	token, err := ctrl.service.Login(c.Context(), input.Username, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

// ==========================================
// CREAR USUARIO
// ==========================================
func (ctrl *AuthController) CrearUsuario(c fiber.Ctx) error {
	type registroInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Rol      string `json:"rol"`
	}
	var input registroInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	// Validar que manden rol válido
	if input.Rol != "ADMINISTRADOR" && input.Rol != "CONSULTA" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Rol inválido. Debe ser ADMINISTRADOR o CONSULTA"})
	}

	usuario, err := ctrl.service.CrearUsuario(c.Context(), input.Username, input.Password, input.Rol)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuario creado exitosamente",
		"data":    usuario,
	})
}
