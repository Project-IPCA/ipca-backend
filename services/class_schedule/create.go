package classschedule

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (classScheduleService *Service) Create(
	request *requests.CreateGroupRequest,
	supervisorId *uuid.UUID,
) (groupId uuid.UUID, err error) {
	groupId = uuid.New()
	group := builders.NewClassScheduleBuilder().
		SetGroupID(groupId).
		SetSupervisorID(supervisorId).
		SetNumber(request.Number).
		SetName(request.Name).
		SetDeptID(request.DeptID).
		SetSemester(request.Semester).
		SetDay(request.Day).
		SetYear(request.Year).
		SetTimeStart(request.TimeStart).
		SetTimeEnd(request.TimeEnd).
		SetLanguage(request.Language).
		Build()

	classScheduleService.DB.Create(&group)
	return groupId, nil
}
