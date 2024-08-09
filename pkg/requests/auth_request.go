package requests

import (
	"github.com/go-playground/validator/v10"
)

type BasicAuth struct {
	Username string `json:"stu_id"   validate:"required" example:"64010000"`
	Password string `json:"password" validate:"required" example:"64010000"`
}

func (ba BasicAuth) Validate() error {
	validate := validator.New()
	return validate.Struct(&ba)
}

type LoginRequest struct {
	BasicAuth
}
