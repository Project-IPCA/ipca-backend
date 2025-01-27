package responses

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
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

type StaffsResponse struct {
	StaffID   uuid.UUID `json:"staff_id"`
	FirstName string    `json:"f_name"`
	LastName  string    `json:"l_name"`
	Role      string    `json:"role"`
}

func NewStaffsResponse(users []models.User) *[]StaffsResponse {
	response := make([]StaffsResponse, 0)
	for _, user := range users {
		response = append(response, StaffsResponse{
			StaffID:   user.UserID,
			FirstName: *user.FirstName,
			LastName:  *user.LastName,
			Role:      *user.Role,
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

type GetExerciseDataResponse struct {
	ExerciseID             uuid.UUID                 `json:"exercise_id"`
	Name                   string                    `json:"name"`
	Content                string                    `json:"content"`
	Sourcecode             string                    `json:"sourcecode"`
	UserDefinedConstraints *json.RawMessage          `json:"user_defined_constraints"`
	SuggestedConstraints   *json.RawMessage          `json:"suggested_constraints"`
	AddedBy                string                    `json:"added_by"`
	CreatedBy              uuid.UUID                 `json:"created_by"`
	TestcaseList           []models.ExerciseTestcase `json:"testcase_list"`
}

func NewGetExerciseDataResponse(
	labExercise models.LabExercise,
	sourceCode string,
) GetExerciseDataResponse {
	response := GetExerciseDataResponse{
		ExerciseID:             labExercise.ExerciseID,
		Name:                   *labExercise.Name,
		Content:                *labExercise.Content,
		Sourcecode:             sourceCode,
		UserDefinedConstraints: labExercise.UserDefinedConstraints,
		SuggestedConstraints:   labExercise.SuggestedConstraints,
		AddedBy:                *labExercise.AddedBy,
		CreatedBy:              *labExercise.CreatedBy,
		TestcaseList:           labExercise.TestcaseList,
	}

	return response
}

type GetAssginStudentExerciseResponse struct {
	ExerciseID  uuid.UUID `json:"exercise_id"`
	ChapterID   uuid.UUID `json:"chapter_id"`
	ChapterName string    `json:"chapter_name"`
	ChapterIdx  int       `json:"chapter_index"`
	Level       string    `json:"level"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	Testcase    string    `json:"testcase"`
	FullMark    int       `json:"full_mark"`
}

func NewGetAssginStudentExerciseResponse(
	labExercise models.LabExercise,
) GetAssginStudentExerciseResponse {
	response := GetAssginStudentExerciseResponse{
		ExerciseID:  labExercise.ExerciseID,
		ChapterID:   *labExercise.ChapterID,
		ChapterName: labExercise.Chapter.Name,
		ChapterIdx:  labExercise.Chapter.ChapterIndex,
		Level:       *labExercise.Level,
		Name:        *labExercise.Name,
		Content:     *labExercise.Content,
		Testcase:    labExercise.Testcase,
		FullMark:    labExercise.FullMark,
	}

	return response
}

type GetRolePermissionResponse struct {
	Permission []string `json:"permission"`
}

func NewGetRolePermissionResponse(rolePermission []models.RolePermission, user models.User) GetRolePermissionResponse {
	permisisonList := make([]string, 0)
	if *user.Role == constants.Role.Supervisor || *user.Role == constants.Role.Beyonder {
		permisisonList = append(permisisonList, constants.PermissionType.DashboardAdmin, constants.PermissionType.ExerciseAdmin, constants.PermissionType.GroupAdmin, constants.PermissionType.StudentAdmin)
	} else {
		for _, permission := range rolePermission {
			permisisonList = append(permisisonList, permission.Permission)
		}
	}

	response := GetRolePermissionResponse{
		Permission: permisisonList,
	}

	return response
}

type GetAllRolePermissionResponse struct {
	Role       string   `json:"role"`
	Permission []string `json:"permission"`
}

// TODO Improve get key unique from db and the filter
func NewGetAllRolePermissionResponse(rolePermission []models.RolePermission) []GetAllRolePermissionResponse {
	permisisonList := make([]GetAllRolePermissionResponse, 0)
	executivePermission := make([]string, 0)
	taPermission := make([]string, 0)

	for _, permisson := range rolePermission {
		if permisson.Role == constants.Role.Executive {
			executivePermission = append(executivePermission, permisson.Permission)
		} else if permisson.Role == constants.Role.Ta {
			taPermission = append(taPermission, permisson.Permission)
		}
	}

	permisisonList = append(permisisonList, GetAllRolePermissionResponse{
		Role:       constants.Role.Executive,
		Permission: executivePermission,
	}, GetAllRolePermissionResponse{
		Role:       constants.Role.Ta,
		Permission: taPermission,
	})

	return permisisonList
}
