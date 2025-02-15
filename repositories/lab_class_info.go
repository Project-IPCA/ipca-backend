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
) {
	labClassInfoRepository.DB.Order("chapter_index ASC").Find(labClassInfos)
}

func (labClassInfoRepository *LabClassInfoRepository) GetLabClassInfoByChapterIndex(
	labClassInfo *models.LabClassInfo,
	chapterIdx int,
){
	labClassInfoRepository.DB.Where("chapter_index",chapterIdx).Find(labClassInfo)
}
