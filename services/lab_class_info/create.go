package labclassinfo

import (
	"github.com/Project-IPCA/ipca-backend/server/builders"
	"github.com/google/uuid"
)

func (labClassInfoService *Service) Create(
	chapterIndex int,
    name         string,
    fullMark     int,
    noItems      int,
) error {
	chapterId := uuid.New()
	labClassInfo := builders.NewLabClassInfoBuilder().
		SetChapterID(chapterId).
		SetChapterIndex(chapterIndex).
		SetName(name).
		SetFullMark(fullMark).
		SetNoItems(noItems).
		Build()
	labClassInfoService.DB.Create(&labClassInfo)
	return nil
}
