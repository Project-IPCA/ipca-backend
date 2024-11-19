package responses

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GetLabChapterInfoResponse struct {
	ChapterName       string               `json:"chapter_name"`
	ChapterId         string               `json:"chapter_id"`
	ChapterIdx        int                  `json:"chapter_idx"`
	GroupId           uuid.UUID            `json:"group_id"`
	GroupSelectedLabs map[string][]string  `json:"group_selected_labs"`
	LabList           map[string][]LabData `json:"lab_list"`
}

type LabData struct {
	ExerciseId string `json:"exercise_id"`
	Name       string `json:"name"`
	ItemId     string `json:"item_id"`
	ChapterId  string `json:"chapter_id"`
}

func NewGetLabChapterInfoResponse(
	labClassInfo models.LabClassInfo,
	groupId uuid.UUID,
	groupChapterSelectItems []models.GroupChapterSelectedItem,
	exerciseList []models.LabExercise,
) GetLabChapterInfoResponse {
	groupSelectedLabs := make(map[string][]string)
	labList := make(map[string][]LabData)

	for i := 1; i <= labClassInfo.NoItems; i++ {
		key := strconv.Itoa(i)
		groupSelectedLabs[key] = []string{}
		labList[key] = []LabData{}
	}

	for _, item := range groupChapterSelectItems {
		itemIdStr := strconv.Itoa(int(item.ItemID))
		if _, exists := groupSelectedLabs[itemIdStr]; exists {
			groupSelectedLabs[itemIdStr] = append(
				groupSelectedLabs[itemIdStr],
				item.ExerciseID.String(),
			)
		}
	}

	for _, exercise := range exerciseList {
		if exercise.Level == nil {
			continue
		}
		levelStr := *exercise.Level
		if _, exists := labList[levelStr]; exists {
			labList[levelStr] = append(labList[levelStr], LabData{
				ExerciseId: exercise.ExerciseID.String(),
				Name:       *exercise.Name,
				ItemId:     levelStr,
				ChapterId:  exercise.ChapterID.String(),
			})
		}
	}

	return GetLabChapterInfoResponse{
		ChapterName:       labClassInfo.Name,
		ChapterId:         labClassInfo.ChapterID.String(),
		ChapterIdx:        labClassInfo.ChapterIndex,
		GroupId:           groupId,
		GroupSelectedLabs: groupSelectedLabs,
		LabList:           labList,
	}
}

type SetChapterPermissionResponse struct {
	ClassID         uuid.UUID  `json:"group_id"`
	ChapterID       uuid.UUID  `json:"chapter_id"`
	AllowAccessType string     `json:"allow_access_type"`
	AccessTimeStart *time.Time `json:"access_time_start"`
	AccessTimeEnd   *time.Time `json:"access_time_end"`
	AllowSubmitType string     `json:"allow_submit_type"`
	SubmitTimeStart *time.Time `json:"submit_time_start"`
	SubmitTimeEnd   *time.Time `json:"submit_time_end"`
	AllowSubmit     bool       `json:"allow_submit"`
	Status          string     `json:"status"`
	AllowAccess     bool       `json:"allow_access"`
	TimeStart       *string    `json:"time_start"`
	TimeEnd         *string    `json:"time_end"`
}

func NewSetChapterPermissionResponse(
	groupChapterPermission models.GroupChapterPermission,
) SetChapterPermissionResponse {
	response := SetChapterPermissionResponse{
		ClassID:         groupChapterPermission.ClassID,
		ChapterID:       groupChapterPermission.ChapterID,
		AllowAccessType: groupChapterPermission.AllowAccessType,
		AccessTimeStart: groupChapterPermission.AccessTimeStart,
		AccessTimeEnd:   groupChapterPermission.AccessTimeEnd,
		AllowSubmitType: groupChapterPermission.AllowSubmitType,
		SubmitTimeStart: groupChapterPermission.SubmitTimeStart,
		SubmitTimeEnd:   groupChapterPermission.SubmitTimeEnd,
		AllowSubmit:     groupChapterPermission.AllowSubmit,
		Status:          groupChapterPermission.Status,
		AllowAccess:     groupChapterPermission.AllowAccess,
		TimeStart:       groupChapterPermission.TimeStart,
		TimeEnd:         groupChapterPermission.TimeEnd,
	}
	return response
}

type SupervisorsResponse struct {
	SupervisorID uuid.UUID `json:"supervisor_id"`
	FirstName    string    `json:"f_name"`
	LastName     string    `json:"l_name"`
}

func NewSupervisorsResponse(supervisors []models.Supervisor) *[]SupervisorsResponse {
	response := make([]SupervisorsResponse, 0)
	for _, supervisor := range supervisors {
		response = append(response, SupervisorsResponse{
			SupervisorID: supervisor.SupervisorID,
			FirstName:    *supervisor.User.FirstName,
			LastName:     *supervisor.User.LastName,
		})
	}
	return &response
}

type StudentResponse struct {
	StudentID uuid.UUID `json:"student_id"`
	FirstName string    `json:"f_name"`
	LastName  string    `json:"l_name"`
}

type LogoutAllStudentResponse struct {
	Message       string            `json:"message"`
	StudentLogout []StudentResponse `json:"student_logout"`
}

func NewLogoutAllStudentResponse(count int, studentList []models.Student) LogoutAllStudentResponse {
	studentLogout := make([]StudentResponse, 0)
	for _, student := range studentList {
		studentLogout = append(studentLogout, StudentResponse{
			StudentID: student.StuID,
			FirstName: *student.User.FirstName,
			LastName:  *student.User.LastName,
		})
	}
	response := LogoutAllStudentResponse{
		Message:       fmt.Sprintf("Logout %d Students Successfully", count),
		StudentLogout: studentLogout,
	}

	return response
}
