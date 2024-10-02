package responses

import (
	"sort"
	"strconv"
	"time"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/google/uuid"
)

type GetAllChapterResponse struct {
	Index int `json:"index"`
	ChapterId string `json:"chapter_id"`
	Name string `json:"name"`
	Marking int `json:"marking"`
	FullMark int `json:"full_mark"`
	IsOpen bool `json:"is_open"`
	LastExerciseSuceess int `json:"last_exercise_success"`
}

func NewGetAllChapter (
	chapterPermission []models.GroupChapterPermission,
	studentChapterItem []models.StudentAssignmentChapterItem,
	labClassInfos	[]models.LabClassInfo,
)*[]GetAllChapterResponse{
	getAllChapter := make([]GetAllChapterResponse,0)
	for _,chapter := range chapterPermission{
		canAccess := false
		if(chapter.AllowAccessType == constants.AccessType.Always || chapter.AllowAccessType == constants.AccessType.TimerPaused){
			canAccess = true
		}else if (chapter.AllowAccessType == constants.AccessType.Timer || chapter.AllowAccessType == constants.AccessType.DateTime){
			if(chapter.AccessTimeStart != nil && chapter.AccessTimeEnd != nil ){
				now := time.Now()
				canAccess = now.After(*chapter.AccessTimeStart) && now.Before(*chapter.AccessTimeEnd)
			}
		}
		var labClassInfo models.LabClassInfo
		for _,labInfo := range labClassInfos{
			if(labInfo.ChapterID == chapter.ChapterID){
				labClassInfo = labInfo
				break
			}
		}
		marking := 0
		currntItem := 0
		studentNotDoneItemList := make([]models.StudentAssignmentChapterItem,0)
		for _,studentItem := range studentChapterItem{
			if(studentItem.ChapterID == chapter.ChapterID){
				marking = marking + studentItem.Marking
				currntItem = currntItem + 1
				if(studentItem.Marking == 0){
					studentNotDoneItemList = append(studentNotDoneItemList, studentItem)
				}
			}
			if(currntItem >= labClassInfo.NoItems){
				break
			}
		}
		
		minNotDone := 5
		if(len(studentNotDoneItemList) >0){
			minNotDone = studentNotDoneItemList[0].ItemID
			for _,notDone := range studentNotDoneItemList{
				if(notDone.ItemID < minNotDone){
					minNotDone = notDone.ItemID
				}
			}
		}
		
		getAllChapter = append(getAllChapter, GetAllChapterResponse{
			Index: chapter.LabClassInfo.ChapterIndex,
			ChapterId: chapter.ChapterID.String(),
			Name: chapter.LabClassInfo.Name,
			Marking: marking,
			FullMark: chapter.LabClassInfo.FullMark,
			IsOpen: canAccess,
			LastExerciseSuceess: minNotDone,
		})
	}
	sort.Slice(getAllChapter, func(i, j int) bool {
		return getAllChapter[i].Index < getAllChapter[j].Index
	})
	return &getAllChapter
}

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
	ExerciseID             uuid.UUID        `json:"exercise_id"`
	ChapterIdx              int       `json:"chapter_index"`
	Level                  string          `json:"level"`
	Name                   string          `json:"name"`
	Content                string          `json:"content"`
	Testcase               string           `json:"testcase"`
	FullMark               int              `json:"full_mark"`
	UserDefinedConstraints *string `json:"user_defined_constraints"`
	SuggestedConstraints   *string `json:"suggested_constraints"`
	TestcaseList		   []TestcaseResponse `json:"testcase_list"`
}

type TestcaseResponse struct{
	TestcaseID       *uuid.UUID `json:"testcase_id"`
	TestcaseContent  string `json:"testcase_content"`
	IsShowStudent    *bool   `json:"show_to_student"`
	TestcaseNote     *string `json:"testcase_note"`
	TestcaseOutput   *string `json:"testcase_output"`
}

func NewGetStudentAssignmentItemResponse(labExercise models.LabExercise)*StudentAssignmentItemResponse{
	testcaseListResponse := make([]TestcaseResponse,0)
	testcaseValid := constants.Testcase.NoInput

	if(labExercise.Testcase == constants.Testcase.Yes){
		testcaseValid = constants.Testcase.Yes
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
	
	response := StudentAssignmentItemResponse{
		ExerciseID: labExercise.ExerciseID,
		ChapterIdx: labExercise.Chapter.ChapterIndex,
		Level: *labExercise.Level,
		Name: *labExercise.Name,
		Content: *labExercise.Content,
		Testcase: testcaseValid,
		FullMark: labExercise.FullMark,
		UserDefinedConstraints: nil,
		SuggestedConstraints: nil,
		TestcaseList: testcaseListResponse,
	}
	return &response
}

type GetStudentWithAssigmentScoreResponse struct{
	GroupId uuid.UUID `json:"group_id"`
	LabInfo []LabInfo `json:"lab_info"`
	StudentList []StudentWithScore `json:"student_list"`
}

type LabInfo struct {
	Name string `json:"name"`
	ChapterId uuid.UUID `json:"chapter_id"`
	ChapterIdx int `json:"chapter_idx"`
	NoItem int `json:"no_item"`
	FullMark int `json:"full_mark"`
}

type StudentWithScore struct {
	ChapterScore map[string]int `json:"chapter_score"`
	Active bool `json:"active"`
	FirstName string `json:"f_name"`
	LastName string `json:"l_name"`
	CanSubmit bool `json:"can_submit"`
	StuID uuid.UUID `json:"stu_id"`
	KmitlID string `json:"kmitl_id"`
	MidtermScore int `json:"midterm_score"`
	Status bool `json:"status"`
}

func NewGetStudentWithAssigmentScoreByGroupID(
	labClassInfo []models.LabClassInfo,
	students []models.Student,
	groupId uuid.UUID,
) GetStudentWithAssigmentScoreResponse {
	labInfo := make([]LabInfo, 0, len(labClassInfo))
	chapterIdMap := make(map[uuid.UUID]int)

	for _, lab := range labClassInfo {
		labInfo = append(labInfo, LabInfo{
			Name:       lab.Name,
			ChapterId:  lab.ChapterID,
			ChapterIdx: lab.ChapterIndex,
			NoItem:     lab.NoItems,
			FullMark:   lab.FullMark,
		})
		chapterIdMap[lab.ChapterID] = lab.ChapterIndex
	}

	studentScore := make([]StudentWithScore, 0, len(students))
	for _, student := range students {
		chapterScore := make(map[string]int)
		assignmentMap := make(map[uuid.UUID]int)

		if(len(student.Assignments)>0){
			for _, assignment := range student.Assignments {
				assignmentMap[assignment.ChapterID] += assignment.Marking
			}
		}

		for chapterId, chapterIdx := range chapterIdMap {
			chapterIdxStr := strconv.Itoa(chapterIdx)
			if score, exists := assignmentMap[chapterId]; exists {
				chapterScore[chapterIdxStr] = score
			} else {
				chapterScore[chapterIdxStr] = 0
			}
		}

		studentScore = append(studentScore, StudentWithScore{
			ChapterScore:  chapterScore,
			Active:        student.User.IsActive,
			FirstName:     *student.User.FirstName,
			LastName:      *student.User.LastName,
			CanSubmit:     student.CanSubmit,
			StuID:         student.StuID,
			KmitlID:       student.KmitlID,
			MidtermScore:  int(student.MidCore),
			Status:        student.User.IsOnline,
		})
	}

	return GetStudentWithAssigmentScoreResponse{
		GroupId:     groupId,
		LabInfo:     labInfo,
		StudentList: studentScore,
	}
}