package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BasicStaff struct {
	StaffID uuid.UUID `json:"staff_id" validate:"required" example:"00000000-0000-0000-0000-000000000000"`
}

type BasicGroup struct {
	SupervisorId *uuid.UUID `json:"supervisor_id" example:"00000000-0000-0000-0000-000000000000"`
	Semester     *int       `json:"semester"   validate:""         example:"1"`
	Number       *int       `json:"number"     validate:""         example:"1"`
	Name         string     `json:"name"       validate:"required" example:"lorem"`
	DeptID       uuid.UUID  `json:"dept_id"    validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	Day          *string    `json:"day"        validate:""         example:"SUNDAY"`
	Year         *int       `json:"year"       validate:""         example:"2023"`
	TimeStart    *string    `json:"time_start" validate:""         example:"09:00:00"`
	TimeEnd      *string    `json:"time_end"   validate:""         example:"12:00:00"`
	Language     *string    `json:"language"   validate:""         example:"12:00:00"`
}

func (bg BasicGroup) Validate() error {
	validate := validator.New()
	return validate.Struct(&bg)
}

type CreateGroupRequest struct {
	BasicGroup
	Staffs []BasicStaff `json:"staffs"`
}

type UpdateMyGroupRequest struct {
	BasicGroup
	Staffs []BasicStaff `json:"staffs"`
}
