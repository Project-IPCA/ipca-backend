package requests

import (
	"github.com/go-playground/validator/v10"
)

type BasicUserInfo struct {
	Avatar             *string `json:"avatar"               example:"www.abc.com"`
	Nickname           *string `json:"nickname"             example:"Doody"                validate:"omitempty,max=20"`
	Gender             *string `json:"gender"               example:"MALE"                 validate:"omitempty,oneof=MALE FEMALE OTHER"`
	DOB                *string `json:"dob"                  example:"2006-01-02"           validate:"omitempty"`
	Email              *string `json:"email"                example:"john.doe@example.com" validate:"omitempty,email"`
	Tel                *string `json:"tel"                  example:"0870000000"           validate:"omitempty,len=10"`
	NewPassword        *string `json:"new_password"         example:"johndoenew"           validate:"omitempty,min=8"`
	ConfirmNewPassword *string `json:"confirm_new_password" example:"johndownew"           validate:"omitempty,min=8"`
	CurrentPassword    string  `json:"current_password"     example:"johndoe"`
}

func (bui BasicUserInfo) Validate() error {
	validate := validator.New()
	return validate.Struct(&bui)
}

type UpdateUserInfoRequest struct {
	BasicUserInfo
}
