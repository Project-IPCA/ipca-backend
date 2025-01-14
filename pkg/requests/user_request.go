package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BasicUserInfo struct {
	Avatar             *string    `json:"avatar"               example:"www.abc.com"`
	DeptID             *uuid.UUID `json:"dept_id"              example:"00000000-0000-0000-0000-000000000000" validate:"omitempty"`
	Nickname           *string    `json:"nickname"             example:"Doody"                                validate:"omitempty,max=20"`
	Gender             *string    `json:"gender"               example:"MALE"                                 validate:"omitempty,oneof=MALE FEMALE OTHER"`
	DOB                *string    `json:"dob"                  example:"2006-01-02"                           validate:"omitempty"`
	Email              *string    `json:"email"                example:"john.doe@example.com"                 validate:"omitempty,email"`
	Tel                *string    `json:"tel"                  example:"0870000000"                           validate:"omitempty,len=10"`
	NewPassword        *string    `json:"new_password"         example:"johndoenew"                           validate:"omitempty,min=8"`
	ConfirmNewPassword *string    `json:"confirm_new_password" example:"johndownew"                           validate:"omitempty,min=8"`
	CurrentPassword    string     `json:"current_password"     example:"johndoe"`
}

func (bui BasicUserInfo) Validate() error {
	validate := validator.New()
	return validate.Struct(&bui)
}

type UpdateUserInfoRequest struct {
	BasicUserInfo
}

type CreateAdminRequest struct {
	Username  string    `json:"username" example:"kanut" validate:"omitempty"`
	Firstname string    `json:"f_name" example:"kanut" validate:"omitempty"`
	Lastname  string    `json:"l_name" example:"handsome" validate:"omitempty"`
	Role      string    `json:"role" example:"SUPERVISOR | TA" validate:"omitempty"`
	Gender    string    `json:"gender" example:"MALE | FEMALE | OTHER" validate:"omitempty"`
	DeptID    uuid.UUID `json:"dept_id" example:"00000000-0000-0000-0000-000000000000" validate:"omitempty"`
}

func (bui CreateAdminRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&bui)
}
