package classschedule

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
)

func (classScheduleService *Service) UpdateMyGroup(
	classSchedule *models.ClassSchedule,
	request *requests.UpdateMyGroupRequest,
) {
	classSchedule.Semester = request.Semester
	classSchedule.Number = request.Number
	classSchedule.Name = request.Name
	classSchedule.DeptID = request.DeptID
	classSchedule.Day = request.Day
	classSchedule.Year = request.Year
	classSchedule.TimeStart = request.TimeStart
	classSchedule.TimeEnd = request.TimeEnd
	classScheduleService.DB.Save(classSchedule)
}

func (classScheduleService *Service) UpdateAllowLogin(
	classSchedule *models.ClassSchedule,
	allowLogin *bool,
) {
	classSchedule.AllowLogin = *allowLogin
	classScheduleService.DB.Save(classSchedule)
}

func (classScheduleService *Service) UpdateAllowUploadPicture(
	classSchedule *models.ClassSchedule,
	allowUploadPicture *bool,
) {
	classSchedule.AllowUploadPic = *allowUploadPicture
	classScheduleService.DB.Save(classSchedule)
}
