package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jos3lo89/vehicle-control-unajma-api/internal/entities"
)

type AuthRepository interface {
	BuscarPorUsername(ctx context.Context, username string) (*entities.Usuario, error)
	CrearUsuario(ctx context.Context, usuario *entities.Usuario) error
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) BuscarPorUsername(ctx context.Context, username string) (*entities.Usuario, error) {
	query := `SELECT id, username, password_hash, rol, activo, creado_en FROM usuarios_web WHERE username = $1`
	var u entities.Usuario
	err := r.db.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Rol, &u.Activo, &u.CreadoEn)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &u, nil
}

func (r *authRepository) CrearUsuario(ctx context.Context, u *entities.Usuario) error {
	// Insertamos y pedimos que nos devuelva los datos autogenerados (ID, activo, creado_en)
	query := `
		INSERT INTO usuarios_web (username, password_hash, rol) 
		VALUES ($1, $2, $3) 
		RETURNING id, activo, creado_en
	`
	err := r.db.QueryRow(ctx, query, u.Username, u.PasswordHash, u.Rol).Scan(&u.ID, &u.Activo, &u.CreadoEn)
	return err
}
