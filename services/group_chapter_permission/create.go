package groupchapterpermission

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (groupChapterPermissionService *Service) Create(
	classId uuid.UUID,
	chapterId uuid.UUID,
	timeStart string,
	timeEnd string,
	allowSubmit bool,
	allowAccess bool,
	status string,
) error {
	gcp := builders.NewGroupChapterPermissionBuilder().
		SetClassID(classId).
		SetChapterID(chapterId).
		SetTimeStart(timeStart).
		SetTimeEnd(timeEnd).
		SetAllowSubmit(allowSubmit).
		SetAllowAccess(allowAccess).
		SetStatus(status).
		Build()

	groupChapterPermissionService.DB.Create(&gcp)
	return nil
}
