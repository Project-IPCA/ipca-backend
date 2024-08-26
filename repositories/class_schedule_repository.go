package repositories

import (
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ClassScheduleRepositoryQ interface {
	GetClassScheduleByGroupID(
		classSchedule *models.ClassSchedule,
		groupId uuid.UUID,
	)
}

type ClassScheduleRepository struct {
	DB *gorm.DB
}

func NewClassScheduleRepository(db *gorm.DB) *ClassScheduleRepository {
	return &ClassScheduleRepository{DB: db}
}

func (classScheduleRepository *ClassScheduleRepository) GetClassScheduleByGroupID(
	classSchedule *models.ClassSchedule,
	groupId uuid.UUID,
) {
	classScheduleRepository.DB.Where("group_id = ?", groupId).Find(classSchedule)
}

func (classScheduleRepository *ClassScheduleRepository) GetClassScheduleByNumber(
	classSchedule *models.ClassSchedule,
	number int,
) {
	classScheduleRepository.DB.Where("number = ?", number).Find(classSchedule)
}

func (classScheduleRepository *ClassScheduleRepository) GetAllClassSchedules(
	classSchedule *[]models.ClassSchedule,
	page string,
	pageSize string,
) {
	defaultPage := 1
	defaultPageSize := 10

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = defaultPage
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = defaultPageSize
	}

	offset := (pageInt - 1) * pageSizeInt

	classScheduleRepository.DB.Preload("Supervisor.User").
		Preload("Department").
		Preload("ClassLabStaffs.Supervisor.User").Preload("Students").
		Offset(offset).
		Limit(pageSizeInt).
		Find(classSchedule)
}

func (classScheduleRepository *ClassScheduleRepository) GetMyClassSchedules(
	classSchedules *[]models.ClassSchedule,
	supervisorId uuid.UUID,
) {
	classScheduleRepository.DB.Where("supervisor_id = ?", supervisorId).Find(classSchedules)
}

func (classScheduleRepository *ClassScheduleRepository) GetClassSchedulePreloadByGroupID(
	classSchedule *models.ClassSchedule,
	classScheduleId uuid.UUID,
) {
	classScheduleRepository.DB.Preload("Supervisor.User").
		Preload("Department").
		Preload("ClassLabStaffs.Supervisor.User").
		Where("group_id = ?", classScheduleId).Find(classSchedule)
}

func (classScheduleRepository *ClassScheduleRepository) GetMyClassSchedulesByQuery(
	classSchedules *[]models.ClassSchedule,
	supervisorId uuid.UUID,
	year string,
	page string,
	pageSize string,
) {
	query := classScheduleRepository.DB.Preload("Department")

	defaultPage := 1
	defaultPageSize := 10

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = defaultPage
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = defaultPageSize
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil && year != "" {
		query = query.Where("year = ?", yearInt)
	}

	offset := (pageInt - 1) * pageSizeInt

	query = query.Where("supervisor_id = ?", supervisorId).Offset(offset).
		Limit(pageSizeInt)

	query.Find(classSchedules)
}
