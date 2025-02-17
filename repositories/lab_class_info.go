package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type LabClassInfoRepositoryQ interface {
	GetCount() int64
}

type LabClassInfoRepository struct {
	DB *gorm.DB
}

func NewLabClassInfoRepository(db *gorm.DB) *LabClassInfoRepository {
	return &LabClassInfoRepository{DB: db}
}

func (labClassInfoRepository *LabClassInfoRepository) GetCount() int64 {
	var count int64
	labClassInfoRepository.DB.Model(models.LabClassInfo{}).Count(&count)
	return count
}

func (labClassInfoRepository *LabClassInfoRepository) GetAllLabClassInfos(
	labClassInfos *[]models.LabClassInfo,
	language string,
) {
	labClassInfoRepository.DB.Where("language = ?",language).Order("chapter_index ASC").Find(labClassInfos)
}

func (labClassInfoRepository *LabClassInfoRepository) GetLabClassInfoByChapterIndex(
	labClassInfo *models.LabClassInfo,
	chapterIdx int,
	language string,
){
	labClassInfoRepository.DB.Where("chapter_index = ? AND language = ?",chapterIdx,language).Find(labClassInfo)
}
