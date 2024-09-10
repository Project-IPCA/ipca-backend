package groupassignmentchapteritem

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (groupAssignmentChapterItemService *Service) Create(
	groupId uuid.UUID,
	chapterId uuid.UUID,
	itemId int64,
	fullMark int,
	timeStart *string,
	timeEnd *string,
	status string,
) error {
	gaci := builders.NewGroupAssignmentChapterItemBuilder().
		SetGroupID(groupId).
		SetChapterID(chapterId).
		SetItemID(itemId).
		SetFullMark(fullMark).
		SetTimeStart(timeStart).
		SetTimeEnd(timeEnd).
		SetStatus(status).
		Build()

	groupAssignmentChapterItemService.DB.Create(&gaci)
	return nil
}
