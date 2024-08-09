package builders

import (
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type UserStudentBuilder struct {
	StuID     uuid.UUID
	StuStuID  string
	FirstName *string
	LastName  *string
	Nickname  *string
	Gender    *string
	DOB       *time.Time
	Avatar    *string
	Email     *string
	Tel       *string
	Group     *int
	Note      *string
	DeptID    *int
	MidCore   float64
	CanSubmit string
}

func NewUserStudentBuilder() *UserStudentBuilder {
	return &UserStudentBuilder{}
}

func (userStudentBuilder *UserStudentBuilder) SetStuID(stuId uuid.UUID) (us *UserStudentBuilder) {
	userStudentBuilder.StuID = stuId
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetStuStuID(
	stuStuId string,
) (us *UserStudentBuilder) {
	userStudentBuilder.StuStuID = stuStuId
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetFirstName(
	firstName string,
) (us *UserStudentBuilder) {
	userStudentBuilder.FirstName = &firstName
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetLastName(
	lastName string,
) (us *UserStudentBuilder) {
	userStudentBuilder.LastName = &lastName
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetNickname(
	nickname string,
) (us *UserStudentBuilder) {
	userStudentBuilder.Nickname = &nickname
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetGender(
	gender string,
) (us *UserStudentBuilder) {
	userStudentBuilder.Gender = &gender
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetDOB(
	dob time.Time,
) (us *UserStudentBuilder) {
	userStudentBuilder.DOB = &dob
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetAvatar(
	avatar string,
) (us *UserStudentBuilder) {
	userStudentBuilder.Avatar = &avatar
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetEmail(
	email string,
) (us *UserStudentBuilder) {
	userStudentBuilder.Email = &email
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetTel(
	tel string,
) (us *UserStudentBuilder) {
	userStudentBuilder.Tel = &tel
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetGroup(
	group int,
) (us *UserStudentBuilder) {
	userStudentBuilder.Group = &group
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetNote(
	note string,
) (us *UserStudentBuilder) {
	userStudentBuilder.Note = &note
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetDeptID(
	deptId int,
) (us *UserStudentBuilder) {
	userStudentBuilder.DeptID = &deptId
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetMidCore(
	midCore float64,
) (us *UserStudentBuilder) {
	userStudentBuilder.MidCore = midCore
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) SetCanSubmit(
	canSubmit string,
) (us *UserStudentBuilder) {
	userStudentBuilder.CanSubmit = canSubmit
	return userStudentBuilder
}

func (userStudentBuilder *UserStudentBuilder) Build() models.UserStudent {
	userStudent := models.UserStudent{
		StuID:        userStudentBuilder.StuID,
		StuStuID:     userStudentBuilder.StuStuID,
		StuFirstname: userStudentBuilder.FirstName,
		StuLastname:  userStudentBuilder.LastName,
		StuNickname:  userStudentBuilder.Nickname,
		StuGender:    userStudentBuilder.Gender,
		StuDOB:       userStudentBuilder.DOB,
		StuAvatar:    userStudentBuilder.Avatar,
		StuEmail:     userStudentBuilder.Email,
		StuTel:       userStudentBuilder.Tel,
		StuGroup:     userStudentBuilder.Group,
		Note:         userStudentBuilder.Note,
		StuDeptID:    userStudentBuilder.DeptID,
		MidCore:      userStudentBuilder.MidCore,
		CanSubmit:    userStudentBuilder.CanSubmit,
	}
	return userStudent
}
