package responses

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ClassStaff struct {
	SupervisorID uuid.UUID `json:"supervisor_id"`
	FirstName    string    `json:"f_name"`
	LastName     string    `json:"l_name"`
}

type Instructor struct {
	SupervisorID uuid.UUID `json:"supervisor_id"`
	FirstName    string    `json:"f_name"`
	LastName     string    `json:"l_name"`
}

type ClassScheduleResponse struct {
	GroupID       uuid.UUID    `json:"group_id"`
	GroupNo       int          `json:"group_no"`
	Department    string       `json:"department"`
	Year          *int         `json:"year"`
	Semester      *int         `json:"semester"`
	Day           *string      `json:"day"`
	TimeStart     *string      `json:"time_start"`
	TimeEnd       *string      `json:"time_end"`
	StudentAmount int          `json:"student_amount"`
	Instructor    Instructor   `json:"instructor"`
	Staff         []ClassStaff `json:"staffs"`
}

func NewClassSchedulesResponse(classSchedules []models.ClassSchedule) *[]ClassScheduleResponse {
	classSchedulesResponse := make([]ClassScheduleResponse, 0)
	for _, classSchedule := range classSchedules {
		classStaffResponse := make([]ClassStaff, 0)
		for _, labStaff := range classSchedule.ClassLabStaffs {
			classStaffResponse = append(classStaffResponse, ClassStaff{
				SupervisorID: labStaff.Supervisor.SupervisorID,
				FirstName:    *labStaff.Supervisor.User.FirstName,
				LastName:     *labStaff.Supervisor.User.LastName,
			})
		}
		classSchedulesResponse = append(classSchedulesResponse, ClassScheduleResponse{
			GroupID:       classSchedule.GroupID,
			GroupNo:       *classSchedule.Number,
			Department:    classSchedule.Department.Name,
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
	return &classSchedulesResponse
}

type MyClassScheduleResponse struct {
	GroupID    uuid.UUID `json:"group_id"`
	GroupNo    int       `json:"group_no"`
	Department string    `json:"department"`
	Year       *int      `json:"year"`
	Semester   *int      `json:"semester"`
	Day        *string   `json:"day"`
	TimeStart  *string   `json:"time_start"`
	TimeEnd    *string   `json:"time_end"`
}

func NewMyClassSchedulesResponse(classSchedules []models.ClassSchedule) *[]MyClassScheduleResponse {
	classSchedulesResponse := make([]MyClassScheduleResponse, 0)
	for _, classSchedule := range classSchedules {
		classSchedulesResponse = append(classSchedulesResponse, MyClassScheduleResponse{
			GroupID:    classSchedule.GroupID,
			GroupNo:    *classSchedule.Number,
			Department: classSchedule.Department.Name,
			Year:       classSchedule.Year,
			Semester:   classSchedule.Semester,
			Day:        classSchedule.Day,
			TimeStart:  classSchedule.TimeStart,
			TimeEnd:    classSchedule.TimeEnd,
		})
	}
	return &classSchedulesResponse
}

type GroupChapterPermission struct {
	ChapterID       uuid.UUID `json:"chapter_id"`
	ChapterIndex    int       `json:"chapter_index"`
	Name            string    `json:"name"`
	AllowAccessType string    `json:"allow_access_type"`
	AllowSubmitType string    `json:"allow_submit_type"`
	FullMark        int       `json:"full_mark"`
}

type ClassScheduleInfoResponse struct {
	GroupID                 uuid.UUID                `json:"group_id"`
	GroupNo                 int                      `json:"group_no"`
	Department              string                   `json:"department"`
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
			SupervisorID: labStaff.Supervisor.SupervisorID,
			FirstName:    *labStaff.Supervisor.User.FirstName,
			LastName:     *labStaff.Supervisor.User.LastName,
		})
	}
	groupChapterPermResponse := make([]GroupChapterPermission, 0)
	for _, gcp := range classSchedule.GroupChapterPermissions {
		groupChapterPermResponse = append(groupChapterPermResponse, GroupChapterPermission{
			ChapterID:       gcp.ChapterID,
			ChapterIndex:    gcp.LabClassInfo.ChapterIndex,
			Name:            gcp.LabClassInfo.Name,
			AllowAccessType: gcp.AllowAccessType,
			AllowSubmitType: gcp.AllowSubmitType,
			FullMark:        gcp.LabClassInfo.FullMark,
		})
	}
	return &ClassScheduleInfoResponse{
		GroupID:       classSchedule.GroupID,
		GroupNo:       *classSchedule.Number,
		Department:    classSchedule.Department.Name,
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
