package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/entities"
)

// RequireAuth valida que el usuario envíe un Bearer Token válido
func RequireAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Obtener el header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Falta el token de autorización"})
		}

		// 2. Extraer el token (Quitar la palabra "Bearer ")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Formato de token inválido. Usa: Bearer <token>"})
		}
		tokenString := parts[1]

		// 3. Validar y parsear el JWT
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secreto_por_defecto_cambiame"
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido o expirado"})
		}

		// 4. Extraer datos y guardarlos en la petición (Locals)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
			c.Locals("rol", claims["rol"])
		}

		return c.Next() // Si todo está bien, pasa a la siguiente función
	}
}

// RequireRole verifica si el usuario tiene el permiso necesario
// Ahora recibe variables de tipo entities.RolSistema
func RequireRole(allowedRoles ...entities.RolSistema) fiber.Handler {
	return func(c fiber.Ctx) error {

		// 1. Extraemos el rol que el JWT guardó en el contexto
		userRole, ok := c.Locals("rol").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Rol no encontrado en el token"})
		}

		// 2. Comparamos el rol del usuario con los roles permitidos
		for _, role := range allowedRoles {
			// Convertimos 'role' (que es RolSistema) a string para poder compararlo
			if userRole == string(role) {
				return c.Next() // Tiene permiso, avanzar
			}
		}

		// 3. Si termina el ciclo y no coincidió, bloqueamos
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Acceso denegado. No tienes los permisos necesarios.",
		})
	}
}
