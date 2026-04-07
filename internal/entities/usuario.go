package entities

import "time"

type Usuario struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Rol          string    `json:"rol"`
	Activo       bool      `json:"activo"`
	CreadoEn     time.Time `json:"creado_en"`
}
