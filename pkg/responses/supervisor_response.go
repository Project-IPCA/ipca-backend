package responses

import (
	"strconv"
	"time"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type GetLabChapterInfoResponse struct {
	ChapterName       string                 `json:"chapter_name"`
	ChapterId         string                 `json:"chapter_id"`
	ChapterIdx        int                    `json:"chapter_idx"`
	GroupId           uuid.UUID              `json:"group_id"`
	GroupSelectedLabs map[string][]string    `json:"group_selected_labs"`
	LabList           map[string][]LabData   `json:"lab_list"`
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
			groupSelectedLabs[itemIdStr] = append(groupSelectedLabs[itemIdStr], item.ExerciseID.String())
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
	ClassID         uuid.UUID     `json:"group_id"`
	ChapterID       uuid.UUID     `json:"chapter_id"`
	AllowAccessType string        `json:"allow_access_type"`
	AccessTimeStart *time.Time    `json:"access_time_start"`
	AccessTimeEnd   *time.Time    `json:"access_time_end"`
	AllowSubmitType string        `json:"allow_submit_type"`
	SubmitTimeStart *time.Time    `json:"submit_time_start"`
	SubmitTimeEnd   *time.Time    `json:"submit_time_end"`
	AllowSubmit     bool          `json:"allow_submit"`
	Status          string        `json:"status"`
	AllowAccess     bool          `json:"allow_access"`
	TimeStart       *string       `json:"time_start"`
	TimeEnd         *string       `json:"time_end"`
}

func NewSetChapterPermissionResponse (
	groupChapterPermission models.GroupChapterPermission,
) SetChapterPermissionResponse{
	response := SetChapterPermissionResponse{
		ClassID: groupChapterPermission.ClassID,
		ChapterID: groupChapterPermission.ChapterID,
		AllowAccessType: groupChapterPermission.AllowAccessType,
		AccessTimeStart: groupChapterPermission.AccessTimeStart,
		AccessTimeEnd: groupChapterPermission.AccessTimeEnd,
		AllowSubmitType: groupChapterPermission.AllowSubmitType,
		SubmitTimeStart: groupChapterPermission.SubmitTimeStart,
		SubmitTimeEnd: groupChapterPermission.SubmitTimeEnd,
		AllowSubmit: groupChapterPermission.AllowSubmit,
		Status: groupChapterPermission.Status,
		AllowAccess: groupChapterPermission.AllowAccess,
		TimeStart: groupChapterPermission.TimeStart,
		TimeEnd: groupChapterPermission.TimeEnd,
	}
	return response
}