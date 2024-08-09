package userstudent

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (userStudentService *Service) Create(
	userId uuid.UUID,
	stuId string,
	firstName string,
	lastName string,
) error {
	userStudent := builders.NewUserStudentBuilder().
		SetStuID(userId).
		SetStuStuID(stuId).
		SetFirstName(firstName).
		SetLastName(lastName).
		Build()
	userStudentService.DB.Create(&userStudent)
	return nil
}
