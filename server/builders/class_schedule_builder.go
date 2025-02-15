package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ClassScheduleBuilder struct {
	GroupID        uuid.UUID
	Number         *int
	Name           string
	DeptID         uuid.UUID
	SupervisorID   *uuid.UUID
	Day            *string
	TimeStart      *string
	TimeEnd        *string
	Year           *int
	Semester       *int
	AllowUploadPic bool
	AllowLogin     bool
	AllowSubmit    bool
	AllowExercise  bool
	Language       *string
}

func NewClassScheduleBuilder() *ClassScheduleBuilder {
	return &ClassScheduleBuilder{}
}

func (classScheduleBuilder *ClassScheduleBuilder) SetGroupID(
	groupId uuid.UUID,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.GroupID = groupId
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetNumber(
	number *int,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.Number = number
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetName(
	name string,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.Name = name
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetDeptID(
	deptId uuid.UUID,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.DeptID = deptId
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetSupervisorID(
	supervisorId *uuid.UUID,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.SupervisorID = supervisorId
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetDay(
	day *string,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.Day = day
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetTimeStart(
	timeStart *string,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.TimeStart = timeStart
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetTimeEnd(
	timeEnd *string,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.TimeEnd = timeEnd
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetYear(
	year *int,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.Year = year
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetSemester(
	semester *int,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.Semester = semester
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetAllowUploadPic(
	allowUploadPic bool,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.AllowUploadPic = allowUploadPic
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetAllowLogin(
	allowLogin bool,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.AllowLogin = allowLogin
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetAllowSubmit(
	allowSubmit bool,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.AllowSubmit = allowSubmit
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetAllowExercise(
	allowExercise bool,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.AllowExercise = allowExercise
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) SetLanguage(
	Language *string,
) (cs *ClassScheduleBuilder) {
	classScheduleBuilder.Language = Language
	return classScheduleBuilder
}

func (classScheduleBuilder *ClassScheduleBuilder) Build() models.ClassSchedule {
	classSchedule := models.ClassSchedule{
		GroupID:        classScheduleBuilder.GroupID,
		Number:         classScheduleBuilder.Number,
		Name:           classScheduleBuilder.Name,
		DeptID:         classScheduleBuilder.DeptID,
		SupervisorID:   classScheduleBuilder.SupervisorID,
		Day:            classScheduleBuilder.Day,
		TimeStart:      classScheduleBuilder.TimeStart,
		TimeEnd:        classScheduleBuilder.TimeEnd,
		Year:           classScheduleBuilder.Year,
		Semester:       classScheduleBuilder.Semester,
		AllowUploadPic: classScheduleBuilder.AllowUploadPic,
		AllowLogin:     classScheduleBuilder.AllowLogin,
		AllowSubmit:    classScheduleBuilder.AllowSubmit,
		AllowExercise:  classScheduleBuilder.AllowExercise,
		Language:       classScheduleBuilder.Language,
	}
	return classSchedule
}
