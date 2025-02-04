package domain

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string
	RegisteredAt time.Time
}

type UserInput struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (i UserInput) Validate() error {
	return validate.Struct(i)
}
