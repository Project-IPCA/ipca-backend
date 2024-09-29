package groupchapterselecteditem

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (service *Service) Create(
	groupID uuid.UUID,
	chapterID uuid.UUID,
	itemID int64,
	exerciseID uuid.UUID,
) error {
	groupChapterSelectedItem := builders.NewGroupChapterSelectedItemBuilder().
		SetGroupID(groupID).
		SetChapterID(chapterID).
		SetItemID(itemID).
		SetExerciseID(exerciseID).
		Build()

	err := service.DB.Create(&groupChapterSelectedItem)
	if err.Error != nil {
		return fmt.Errorf("fail when insert data to group_chapter_select_items : %v",err.Error)
	}
	
	return nil
}