package repositories

import (
	"strconv"
	"strings"

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
	classSchedules *[]models.ClassSchedule,
) {
	classScheduleRepository.DB.Find(classSchedules)
}

func (classScheduleRepository *ClassScheduleRepository) GetAllClassSchedulesByQuery(
	classSchedules *[]models.ClassSchedule,
	instructorId string,
	staffIds []string,
	year string,
	semester string,
	day string,
	page string,
	pageSize string,
) int64 {
	baseQuery := classScheduleRepository.DB.Model(models.ClassSchedule{}).
		Preload("Supervisor.User").
		Preload("Department").
		Preload("ClassLabStaffs.Supervisor.User").
		Preload("Students")

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

	instructorIdUuid, err := uuid.Parse(instructorId)
	if instructorIdUuid != uuid.Nil && instructorId != "" {
		baseQuery = baseQuery.Where("supervisor_id = ?", instructorIdUuid)
	}

	var staffIdList []string
	for _, id := range staffIds {
		staffIdList = append(staffIdList, strings.Split(id, ",")...)
	}

	if len(staffIdList) > 0 {
		baseQuery = baseQuery.Joins("JOIN class_lab_staffs ON class_lab_staffs.class_id = class_schedules.group_id").
			Where("class_lab_staffs.staff_id IN (?)", staffIdList).
			Group("class_schedules.group_id").
			Having("COUNT(DISTINCT class_lab_staffs.staff_id) = ?", len(staffIdList))
	}

	yearInt, err := strconv.Atoi(year)
	if err == nil && year != "" {
		baseQuery = baseQuery.Where("year = ?", yearInt)
	}

	semesterInt, err := strconv.Atoi(semester)
	if err == nil && semester != "" {
		baseQuery = baseQuery.Where("semester = ?", semesterInt)
	}

	if day != "" {
		baseQuery = baseQuery.Where("day = ?", day)
	}

	var totalGroups int64
	baseQuery.Count(&totalGroups)

	offset := (pageInt - 1) * pageSizeInt

	baseQuery.Offset(offset).Limit(pageSizeInt).Find(classSchedules)

	return totalGroups
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
		Preload("GroupChapterPermissions.LabClassInfo").
		Where("group_id = ?", classScheduleId).Find(classSchedule)
}

func (classScheduleRepository *ClassScheduleRepository) GetMyClassSchedulesByQuery(
	classSchedules *[]models.ClassSchedule,
	supervisorId uuid.UUID,
	year string,
	page string,
	pageSize string,
) int64 {
	baseQuery := classScheduleRepository.DB.Model(&models.ClassSchedule{}).
		Preload("Department").
		Preload("Supervisor.User").
		Preload("Department").
		Preload("ClassLabStaffs.Supervisor.User").
		Preload("Students")

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
	if err == nil && year != "" {
		baseQuery = baseQuery.Where("year = ?", yearInt)
	}

	baseQuery = baseQuery.Where(
		classScheduleRepository.DB.Where("supervisor_id = ?", supervisorId).Or(
			"EXISTS (SELECT 1 FROM class_lab_staffs WHERE class_lab_staffs.class_id = class_schedules.group_id AND class_lab_staffs.staff_id = ?)",
			supervisorId,
		),
	)

	var totalClassSchedules int64
	baseQuery.Count(&totalClassSchedules)

	offset := (pageInt - 1) * pageSizeInt
	baseQuery.Offset(offset).Limit(pageSizeInt).Find(classSchedules)

	return totalClassSchedules
}
