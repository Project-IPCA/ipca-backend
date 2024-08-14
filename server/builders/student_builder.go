package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type StudentBuilder struct {
	StuID     uuid.UUID
	KmitlID   string
	GroupID   *uuid.UUID
	Note      *string
	DeptID    *uuid.UUID
	MidCore   float64
	CanSubmit bool
}

func NewUserStudentBuilder() *StudentBuilder {
	return &StudentBuilder{}
}

func (studentBuilder *StudentBuilder) SetStuID(stuId uuid.UUID) (us *StudentBuilder) {
	studentBuilder.StuID = stuId
	return studentBuilder
}

func (studentBuilder *StudentBuilder) SetKmitlID(kmitlId string) (us *StudentBuilder) {
	studentBuilder.KmitlID = kmitlId
	return studentBuilder
}

func (studentBuilder *StudentBuilder) SetGroupID(groupId string) (us *StudentBuilder) {
	studentBuilder.KmitlID = groupId
	return studentBuilder
}

func (studentBuilder *StudentBuilder) SetNote(note string) (us *StudentBuilder) {
	studentBuilder.Note = &note
	return studentBuilder
}

func (studentBuilder *StudentBuilder) SetDeptID(deptId uuid.UUID) (us *StudentBuilder) {
	studentBuilder.DeptID = &deptId
	return studentBuilder
}

func (studentBuilder *StudentBuilder) SetMidCore(
	midCore float64,
) (us *StudentBuilder) {
	studentBuilder.MidCore = midCore
	return studentBuilder
}

func (studentBuilder *StudentBuilder) SetCanSubmit(
	canSubmit bool,
) (us *StudentBuilder) {
	studentBuilder.CanSubmit = canSubmit
	return studentBuilder
}

func (studentBuilder *StudentBuilder) Build() models.Student {
	student := models.Student{
		StuID:     studentBuilder.StuID,
		KmitlID:   studentBuilder.KmitlID,
		GroupID:   studentBuilder.GroupID,
		Note:      studentBuilder.Note,
		MidCore:   studentBuilder.MidCore,
		DeptID:    studentBuilder.DeptID,
		CanSubmit: studentBuilder.CanSubmit,
	}
	return student
}
