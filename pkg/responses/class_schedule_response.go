package responses

import (
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ClassStaff struct {
	StaffID   uuid.UUID `json:"staff_id"`
	FirstName string    `json:"f_name"`
	LastName  string    `json:"l_name"`
	Role      string    `json:"role"`
}

type Instructor struct {
	SupervisorID uuid.UUID `json:"supervisor_id"`
	FirstName    string    `json:"f_name"`
	LastName     string    `json:"l_name"`
}

type ClassSchedule struct {
	GroupID       uuid.UUID    `json:"group_id"`
	GroupNo       int          `json:"group_no"`
	Department    Department   `json:"department"`
	Year          *int         `json:"year"`
	Semester      *int         `json:"semester"`
	Day           *string      `json:"day"`
	TimeStart     *string      `json:"time_start"`
	TimeEnd       *string      `json:"time_end"`
	StudentAmount int          `json:"student_amount"`
	Instructor    Instructor   `json:"instructor"`
	Staff         []ClassStaff `json:"staffs"`
}

type AvailableGroupFilter struct {
	Instructors []ClassStaff `json:"instructors"`
	Staffs      []ClassStaff `json:"staffs"`
	Years       []*int       `json:"years"`
}

type AvailableGroupsResponse struct {
	AvailableGroups []ClassSchedule      `json:"available_groups"`
	Filter          AvailableGroupFilter `json:"filters"`
	Pagination      Pagination           `json:"pagination"`
}

func NewClassSchedulesResponse(
	filteredClassSchedules []models.ClassSchedule,
	allClassSchedules []models.ClassSchedule,
	staffs []models.User,
	page string,
	pageSize string,
	totalClassSchedules int64,
) *AvailableGroupsResponse {
	classSchedules := make([]ClassSchedule, 0)
	for _, classSchedule := range filteredClassSchedules {
		classStaffResponse := make([]ClassStaff, 0)
		for _, labStaff := range classSchedule.ClassLabStaffs {
			classStaffResponse = append(classStaffResponse, ClassStaff{
				StaffID:   labStaff.User.UserID,
				FirstName: *labStaff.User.FirstName,
				LastName:  *labStaff.User.LastName,
				Role:      *labStaff.User.Role,
			})
		}
		departmentData := Department{
			DeptID: classSchedule.Department.DeptID,
			NameTH: classSchedule.Department.NameTH,
			NameEN: classSchedule.Department.NameEN,
		}
		classSchedules = append(classSchedules, ClassSchedule{
			GroupID:       classSchedule.GroupID,
			GroupNo:       *classSchedule.Number,
			Department:    departmentData,
			Year:          classSchedule.Year,
			Semester:      classSchedule.Semester,
			Day:           classSchedule.Day,
			TimeStart:     classSchedule.TimeStart,
			TimeEnd:       classSchedule.TimeEnd,
			StudentAmount: len(classSchedule.Students),
			Instructor: Instructor{
				SupervisorID: classSchedule.Supervisor.SupervisorID,
				FirstName:    *classSchedule.Supervisor.User.FirstName,
				LastName:     *classSchedule.Supervisor.User.LastName,
			},
			Staff: classStaffResponse,
		})
	}

	allStaffs := make([]ClassStaff, 0)
	for _, staff := range staffs {
		allStaffs = append(allStaffs, ClassStaff{
			StaffID:   staff.UserID,
			FirstName: *staff.FirstName,
			LastName:  *staff.LastName,
			Role:      *staff.Role,
		})
	}

	yearMap := make(map[int]bool)
	for _, classSchedule := range allClassSchedules {
		yearMap[*classSchedule.Year] = true
	}
	var uniqueYears []*int
	for year := range yearMap {
		uniqueYears = append(uniqueYears, &year)
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	pages := int(math.Ceil(float64(totalClassSchedules) / float64(pageSizeInt)))

	return &AvailableGroupsResponse{
		AvailableGroups: classSchedules,
		Filter: AvailableGroupFilter{
			Years:       uniqueYears,
			Instructors: allStaffs,
			Staffs:      allStaffs,
		},
		Pagination: Pagination{
			Page:     pageInt,
			PageSize: pageSizeInt,
			Pages:    pages,
		},
	}
}

type MyGroupFilter struct {
	Year []*int `json:"year"`
}

type MyGroupResponse struct {
	MyGroups   []ClassSchedule `json:"my_groups"`
	Filter     MyGroupFilter   `json:"filters"`
	Pagination Pagination      `json:"pagination"`
}

func NewMyClassSchedulesResponse(
	filteredClassSchedules []models.ClassSchedule,
	allClassSchedules []models.ClassSchedule,
	page string,
	pageSize string,
	totalClassScheduls int,
) *MyGroupResponse {
	classSchedules := make([]ClassSchedule, 0)
	for _, classSchedule := range filteredClassSchedules {
		classStaffResponse := make([]ClassStaff, 0)
		for _, labStaff := range classSchedule.ClassLabStaffs {
			classStaffResponse = append(classStaffResponse, ClassStaff{
				StaffID:   labStaff.User.UserID,
				FirstName: *labStaff.User.FirstName,
				LastName:  *labStaff.User.LastName,
				Role:      *labStaff.User.Role,
			})
		}
		departmentData := Department{
			DeptID: classSchedule.Department.DeptID,
			NameTH: classSchedule.Department.NameTH,
			NameEN: classSchedule.Department.NameEN,
		}
		classSchedules = append(classSchedules, ClassSchedule{
			GroupID:       classSchedule.GroupID,
			GroupNo:       *classSchedule.Number,
			Department:    departmentData,
			Year:          classSchedule.Year,
			Semester:      classSchedule.Semester,
			Day:           classSchedule.Day,
			TimeStart:     classSchedule.TimeStart,
			TimeEnd:       classSchedule.TimeEnd,
			StudentAmount: len(classSchedule.Students),
			Instructor: Instructor{
				SupervisorID: classSchedule.Supervisor.SupervisorID,
				FirstName:    *classSchedule.Supervisor.User.FirstName,
				LastName:     *classSchedule.Supervisor.User.LastName,
			},
			Staff: classStaffResponse,
		})
	}

	yearMap := make(map[int]bool)

	for _, classSchedule := range allClassSchedules {
		yearMap[*classSchedule.Year] = true
	}

	var uniqueYears []*int
	for year := range yearMap {
		uniqueYears = append(uniqueYears, &year)
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	pages := int(math.Ceil(float64(totalClassScheduls) / float64(pageSizeInt)))

	return &MyGroupResponse{
		MyGroups: classSchedules,
		Filter: MyGroupFilter{
			Year: uniqueYears,
		},
		Pagination: Pagination{
			Page:     pageInt,
			PageSize: pageSizeInt,
			Pages:    pages,
		},
	}
}

type Permission struct {
	Type      string     `json:"type"`
	TimeStart *time.Time `json:"time_start"`
	TimeEnd   *time.Time `json:"time_end"`
}

type GroupChapterPermission struct {
	ChapterID    uuid.UUID  `json:"chapter_id"`
	ChapterIndex int        `json:"chapter_index"`
	Name         string     `json:"name"`
	AllowAccess  Permission `json:"allow_access"`
	AllowSubmit  Permission `json:"allow_submit"`
	FullMark     int        `json:"full_mark"`
}

type ClassScheduleInfoResponse struct {
	GroupID                 uuid.UUID                `json:"group_id"`
	GroupNo                 int                      `json:"group_no"`
	Name                    string                   `json:"name"`
	Department              Department               `json:"department"`
	Year                    *int                     `json:"year"`
	Semester                *int                     `json:"semester"`
	Day                     *string                  `json:"day"`
	TimeStart               *string                  `json:"time_start"`
	TimeEnd                 *string                  `json:"time_end"`
	AllowLogin              bool                     `json:"allow_login"`
	AllowUploadProfile      bool                     `json:"allow_upload_profile"`
	StudentAmount           int                      `json:"student_amount"`
	Instructor              Instructor               `json:"instructor"`
	Staffs                  []ClassStaff             `json:"staffs"`
	GroupChapterPermissions []GroupChapterPermission `json:"group_chapter_permissions"`
}

func NewClassScheduleInfoResponse(classSchedule models.ClassSchedule) *ClassScheduleInfoResponse {
	classStaffResponse := make([]ClassStaff, 0)
	for _, labStaff := range classSchedule.ClassLabStaffs {
		classStaffResponse = append(classStaffResponse, ClassStaff{
			StaffID:   labStaff.User.UserID,
			FirstName: *labStaff.User.FirstName,
			LastName:  *labStaff.User.LastName,
			Role:      *labStaff.User.Role,
		})
	}
	groupChapterPermResponse := make([]GroupChapterPermission, 0)
	for _, gcp := range classSchedule.GroupChapterPermissions {
		groupChapterPermResponse = append(groupChapterPermResponse, GroupChapterPermission{
			ChapterID:    gcp.ChapterID,
			ChapterIndex: gcp.LabClassInfo.ChapterIndex,
			Name:         gcp.LabClassInfo.Name,
			AllowAccess: Permission{
				Type:      gcp.AllowAccessType,
				TimeStart: gcp.AccessTimeStart,
				TimeEnd:   gcp.AccessTimeEnd,
			},
			AllowSubmit: Permission{
				Type:      gcp.AllowSubmitType,
				TimeStart: gcp.SubmitTimeStart,
				TimeEnd:   gcp.SubmitTimeEnd,
			},
			FullMark: gcp.LabClassInfo.FullMark,
		})
	}

	departmentData := Department{
		DeptID: classSchedule.Department.DeptID,
		NameTH: classSchedule.Department.NameTH,
		NameEN: classSchedule.Department.NameEN,
	}

	return &ClassScheduleInfoResponse{
		GroupID:       classSchedule.GroupID,
		GroupNo:       *classSchedule.Number,
		Name:          classSchedule.Name,
		Department:    departmentData,
		Year:          classSchedule.Year,
		Semester:      classSchedule.Semester,
		Day:           classSchedule.Day,
		TimeStart:     classSchedule.TimeStart,
		TimeEnd:       classSchedule.TimeEnd,
		StudentAmount: len(classSchedule.Students),
		Instructor: Instructor{
			SupervisorID: classSchedule.Supervisor.SupervisorID,
			FirstName:    *classSchedule.Supervisor.User.FirstName,
			LastName:     *classSchedule.Supervisor.User.LastName,
		},
		AllowLogin:              classSchedule.AllowLogin,
		AllowUploadProfile:      classSchedule.AllowUploadPic,
		Staffs:                  classStaffResponse,
		GroupChapterPermissions: groupChapterPermResponse,
	}
}

type Department struct {
	DeptID uuid.UUID `json:"dept_id"`
	NameTH string    `json:"name_th"`
	NameEN string    `json:"name_en"`
}

type MyGroupInfoResponse struct {
	GroupID    uuid.UUID    `json:"group_id"`
	GroupNo    int          `json:"group_no"`
	Name       string       `json:"name"`
	Department Department   `json:"department"`
	Year       *int         `json:"year"`
	Semester   *int         `json:"semester"`
	Day        *string      `json:"day"`
	TimeStart  *string      `json:"time_start"`
	TimeEnd    *string      `json:"time_end"`
	Staff      []ClassStaff `json:"staffs"`
}

func NewMyClassScheduleInfoResponse(
	classSchedule models.ClassSchedule,
) *MyGroupInfoResponse {
	classStaffResponse := make([]ClassStaff, 0)
	for _, labStaff := range classSchedule.ClassLabStaffs {
		classStaffResponse = append(classStaffResponse, ClassStaff{
			StaffID:   labStaff.User.UserID,
			FirstName: *labStaff.User.FirstName,
			LastName:  *labStaff.User.LastName,
			Role:      *labStaff.User.Role,
		})
	}

	return &MyGroupInfoResponse{
		GroupID: classSchedule.GroupID,
		GroupNo: *classSchedule.Number,
		Name:    classSchedule.Name,
		Department: Department{
			DeptID: classSchedule.Department.DeptID,
			NameTH: classSchedule.Department.NameTH,
			NameEN: classSchedule.Department.NameEN,
		},
		Year:      classSchedule.Year,
		Semester:  classSchedule.Semester,
		Day:       classSchedule.Day,
		TimeStart: classSchedule.TimeStart,
		TimeEnd:   classSchedule.TimeEnd,
		Staff:     classStaffResponse,
	}
}
