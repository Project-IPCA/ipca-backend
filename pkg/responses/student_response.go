package responses

import (
	"time"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/google/uuid"
)

type GetChapterListResponse struct {
	AccessTimeEnd   *time.Time    `json:"access_time_end"`
	AccessTimeStart *time.Time    `json:"access_time_start"`
	AllowAccess     bool          `json:"allow_access"`
	AllowAccessType string        `json:"allow_access_type"`
	AllowSubmit     bool          `json:"allow_submit"`
	AllowSubmitType string        `json:"allow_submit_type"`
	ChapterFullMark int				`json:"chapter_full_mark"`
	ChapterID       uuid.UUID     `json:"chapter_id"`
	Name			string			`json:"chapter_name"`
	ClassID         uuid.UUID     `json:"class_id"`
	NoItems			int				`json:"no_items"`
	Status          string        `json:"status"`
	SubmitTimeEnd   *time.Time    `json:"submit_time_end"`
	SubmitTimeStart *time.Time    `json:"submit_time_start"`
	TimeEnd         *string       `json:"time_end"`
	TimeStart       *string       `json:"time_start"`
	Items			[]ChapterItemsResponse `json:"items"`
}

type ChapterItemsResponse struct{
	ChapterIdx	int `json:"chapter_idx"`
	FullMark 	int	`json:"full_mark"`
	ItemIdx		int `json:"item_idx"`
	Status      *string  `json:"status"`
	 Marking int	`json:"marking"`
	TimeEnd         *string       `json:"time_end"`
	TimeStart       *string       `json:"time_start"`
}

func NewGetChapterListResponse(
	chapterPermission []models.GroupChapterPermission,
	assignChapterItem []models.GroupAssignmentChapterItem,
	studentChapterItem []models.StudentAssignmentChapterItem,
) *[]GetChapterListResponse {
	getChapterList := make([]GetChapterListResponse,0)

	for _,chapter := range chapterPermission{
		itemData := make([]ChapterItemsResponse,0)
		
		for _,item := range assignChapterItem{
			if(item.ChapterID == chapter.ChapterID){
				studentMarking := 0
				for _,studentChapter := range studentChapterItem{
				if(studentChapter.ItemID == int(item.ItemID) && studentChapter.ChapterID == chapter.ChapterID){
					studentMarking = studentChapter.Marking
					}
				}
				chapterItem := ChapterItemsResponse{
					ChapterIdx: chapter.LabClassInfo.ChapterIndex,
					FullMark: item.FullMark,
					ItemIdx: int(item.ItemID),
					Status: item.Status,
					Marking: studentMarking,
					TimeEnd: item.TimeEnd,
					TimeStart: item.TimeStart,
				}
				itemData = append(itemData, chapterItem)
			}
		}

		chapterData := GetChapterListResponse{
			AccessTimeEnd: chapter.AccessTimeEnd,
			AccessTimeStart: chapter.AccessTimeStart,
			AllowAccess: chapter.AllowAccess,
			AllowAccessType: chapter.AllowAccessType,
			AllowSubmit: chapter.AllowSubmit,
			AllowSubmitType: chapter.AllowSubmitType,
			ChapterFullMark: chapter.LabClassInfo.FullMark,
			ChapterID: chapter.ChapterID,
			Name: chapter.LabClassInfo.Name,
			ClassID: chapter.ClassID,
			NoItems: chapter.LabClassInfo.NoItems,
			Status: chapter.Status,
			SubmitTimeEnd: chapter.SubmitTimeEnd,
			SubmitTimeStart: chapter.SubmitTimeStart,
			TimeEnd: chapter.TimeEnd,
			TimeStart: chapter.TimeStart,
			Items: itemData,
		}
		getChapterList = append(getChapterList, chapterData)
	}

	return &getChapterList
}

type StudentAssignmentItemResponse struct{
	
}

type TestcaseResponse struct{
	TestcaseID       *uuid.UUID `json:"testcase_id"`
	TestcaseContent  string `json:"testcase_content"`
	IsShowStudent    *bool   `json:"show_to_student"`
	TestcaseNote     *string `json:"testcase_note"`
	TestcaseOutput   *string `json:"testcase_output"`
}

func NewGetStudentAssignmentItemResponse(labExercise models.LabExercise){
	testcaseListResponse := make([]TestcaseResponse,0)

	if(labExercise.Testcase == constants.Testcase.Yes){
		for _,testcase := range labExercise.TestcaseList{
			if(testcase.IsReady == "yes" && *testcase.IsActive){
				testcaseContent := "Hidden"
				testcaseOutput := "Hidden"
				if(*testcase.IsShowStudent){
					testcaseContent = testcase.TestcaseContent
					testcaseOutput = *testcase.TestcaseOutput
				}
				testcaseListResponse = append(testcaseListResponse, TestcaseResponse{
					TestcaseID: testcase.TestcaseID,
					TestcaseContent: testcaseContent,
					IsShowStudent: testcase.IsShowStudent,
					TestcaseOutput: &testcaseOutput,
					TestcaseNote: testcase.TestcaseNote,
				})
			}
			
		}
	}
}