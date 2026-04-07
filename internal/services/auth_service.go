package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/entities"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	CrearUsuario(ctx context.Context, username, password, rol string) (*entities.Usuario, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	usuario, err := s.repo.BuscarPorUsername(ctx, username)
	if err != nil {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	if !usuario.Activo {
		return "", errors.New("usuario inactivo")
	}

	// Comparar contraseñas
	err = bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	// Generar Token JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secreto_por_defecto_cambiame" // Solo para evitar crasheos si olvidas el .env
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  usuario.ID,
		"username": usuario.Username,
		"rol":      usuario.Rol,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expira en 24h
	})

	return token.SignedString([]byte(secret))
}

func (s *authService) CrearUsuario(ctx context.Context, username, password, rol string) (*entities.Usuario, error) {
	// 1. Encriptar contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al encriptar la contraseña")
	}

	// 2. Preparar entidad
	usuario := &entities.Usuario{
		Username:     username,
		PasswordHash: string(hash),
		Rol:          rol,
	}

	// 3. Guardar en BD
	err = s.repo.CrearUsuario(ctx, usuario)
	if err != nil {
		return nil, errors.New("error al crear usuario, es posible que el username ya exista")
	}

	return usuario, nil
}
