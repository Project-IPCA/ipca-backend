package classschedule

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (classScheduleService *Service) Create(
	request *requests.CreateGroupRequest,
) (groupId uuid.UUID, err error) {
	groupId = uuid.New()
	group := builders.NewClassScheduleBuilder().
		SetGroupID(groupId).
		SetSupervisorID(request.SupervisorID).
		SetNumber(request.Number).
		SetName(request.Name).
		SetDeptID(request.DeptID).
		SetDay(request.Day).
		SetYear(request.Year).
		SetTimeStart(request.TimeStart).
		SetTimeEnd(request.TimeEnd).
		Build()

	classScheduleService.DB.Create(&group)
	return groupId, nil
}
