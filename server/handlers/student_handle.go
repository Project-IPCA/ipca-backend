package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	minioclient "github.com/Project-IPCA/ipca-backend/minio_client"
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/pkg/utils"
	"github.com/Project-IPCA/ipca-backend/rabbitmq_client"
	"github.com/Project-IPCA/ipca-backend/redis_client"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
	activitylog "github.com/Project-IPCA/ipca-backend/services/activity_log"
	exercisesubmission "github.com/Project-IPCA/ipca-backend/services/exercise_submission"
	studentassignmentchapteritem "github.com/Project-IPCA/ipca-backend/services/student_assignment_chapter_item"
	"github.com/Project-IPCA/ipca-backend/services/token"
)

type StudentHandler struct {
	server *s.Server
}

func NewStudentHandle(server *s.Server) *StudentHandler {
	return &StudentHandler{server: server}
}

// @Description Exercise Submit
// @ID student-exercise-submit
// @Tags Student
// @Accept json
// @Produce json
// @Param params body	requests.ExcerciseSubmitRequest	true	"Exercise Submit Request"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/student/exercise_submit [post]
func (studentHandler *StudentHandler) ExerciseSubmit(c echo.Context) error {
	exerciseSubmitReq := new(requests.ExcerciseSubmitRequest)
	if err := c.Bind(exerciseSubmitReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := exerciseSubmitReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}
	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID

	existUser := models.User{}
	userRepository := repositories.NewUserRepository(studentHandler.server.DB)
	userRepository.GetUserByUserID(&existUser, userId)

	if *existUser.Role != constants.Role.Student || !existUser.Student.CanSubmit {
		return responses.ErrorResponse(c, http.StatusForbidden, "You Can't Submmit Now")
	}

	chapterUuid, err := uuid.Parse(exerciseSubmitReq.ChapterID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "GroupID Is Not UUID")
	}

	// TODO add validate group permission and assign chapter item
	var studentAssignChapterItem models.StudentAssignmentChapterItem
	studentAssignChapterItemRepo := repositories.NewStudentAssignChapterItemRepository(
		studentHandler.server.DB,
	)
	studentAssignChapterItemRepo.GetStudentAssignChapterItem(
		&studentAssignChapterItem,
		userId,
		chapterUuid,
		exerciseSubmitReq.ItemId,
	)

	var submissionList []models.ExerciseSubmission
	exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(studentHandler.server.DB)
	exerciseSubmissionRepo.GetStudentSubmission(existUser.UserID, chapterUuid, &submissionList)

	attemps := len(submissionList) + 1
	filename := fmt.Sprintf("%s-%04d*.py", existUser.Username, attemps)
	tempFile, err := utils.CreateTempFile(filename, exerciseSubmitReq.Sourcecode)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Create Temp File Fail %s", err),
		)
	}
	defer os.Remove(tempFile.Name())

	minioAction := minioclient.NewMinioAction(studentHandler.server.Minio)
	uploadFileName, err := minioAction.UploadToMinio(
		tempFile,
		studentHandler.server.Config.Minio.BucketStudentCode,
		true,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	exerciseSubmissionService := exercisesubmission.NewExerciseSubmissionService(
		studentHandler.server.DB,
	)
	isInfLoop := false
	submissionId, err := exerciseSubmissionService.Create(
		existUser.UserID,
		*studentAssignChapterItem.ExerciseID,
		uploadFileName,
		0,
		&isInfLoop,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	exerciseTestcaseRepo := repositories.NewExerciseTestcaseRepository(studentHandler.server.DB)
	var testcaseList []models.ExerciseTestcase
	exerciseTestcaseRepo.GetTestcasesByExerciseID(
		*studentAssignChapterItem.ExerciseID,
		&testcaseList,
	)

	filterTestcase := make([]models.ExerciseTestcase, 0)
	for _, testcase := range testcaseList {
		if testcase.IsActive != nil && *testcase.IsActive {
			filterTestcase = append(filterTestcase, testcase)
		}
	}

	labExerciseRepo := repositories.NewLabExerciseRepository(studentHandler.server.DB)
	var labExercise models.LabExercise
	labExerciseRepo.GetLabExerciseByID(*studentAssignChapterItem.ExerciseID, &labExercise)

	logAction := models.LogExerciseSubmissionAction{
		StuId:              userId,
		JobId:              exerciseSubmitReq.JobId,
		Status:             "Pending",
		SubmissionId:       submissionId,
		Attempt:            fmt.Sprintf("%04d", attemps),
		SourcecodeFilename: filename,
	}

	logActionString, err := json.Marshal(logAction)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	remoteIP := c.RealIP()
	if remoteIP == "" {
		remoteIP = c.Request().RemoteAddr
	}
	userAgent := c.Request().UserAgent()

	logData := requests.LogDataInfo{
		GroupID:  *existUser.Student.GroupID,
		Username: existUser.Username,
		RemoteIP: remoteIP,
		Agent:    userAgent,
		PageName: "exercise_submit",
		Action:   logAction,
	}

	rabbitMessage := requests.ExerciseSubmissionRabbitMessage{
		JobId:        exerciseSubmitReq.JobId,
		JobType:      "exercise-submit",
		LogData:      logData,
		SubmissionId: submissionId,
		SourceCode:   exerciseSubmitReq.Sourcecode,
		TestCaseList: filterTestcase,
		StudentID:    userId,
		ChapterId:    chapterUuid,
		ItemId:       exerciseSubmitReq.ItemId,
	}

	rabbit := rabbitmq_client.NewRabbitMQAction(
		studentHandler.server.RabitMQ,
		studentHandler.server.Config,
	)
	err = rabbit.SendQueue(rabbitMessage)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Error While Send Queue RabbitMQ",
		)
	}

	mockPort := 0

	// TODO remote port
	activitylogService := activitylog.NewActivityLogService(studentHandler.server.DB)
	insertLog, err := activitylogService.Create(
		existUser.Student.GroupID,
		existUser.Username,
		remoteIP,
		&mockPort,
		&userAgent,
		constants.LogPage.ExerciseSubmit,
		string(logActionString),
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Save Activity Log")
	}

	redis := redis_client.NewRedisAction(studentHandler.server.Redis)
	redis.PublishMessage(fmt.Sprintf("logs:%s", existUser.Student.GroupID), insertLog)

	response := responses.NewExerciseSubmitResponse(exerciseSubmitReq.JobId.String())

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get All Chapter
// @ID student-get-all-chapter
// @Tags Student
// @Accept json
// @Produce json
// @Success 200		{array}	responses.GetAllChapterResponse
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/student/all_chapter [get]
func (StudentHandler *StudentHandler) GetALLChapter(c echo.Context) error {
	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID

	var existUser models.User
	userRepo := repositories.NewUserRepository(StudentHandler.server.DB)
	userRepo.GetUserByUserID(&existUser, userId)
	if *existUser.Role != constants.Role.Student {
		return responses.ErrorResponse(c, http.StatusForbidden, "This User Not Student")
	}

	var labClassInfos []models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(StudentHandler.server.DB)
	labClassInfoRepo.GetAllLabClassInfos(&labClassInfos)

	var groupChapterPermission []models.GroupChapterPermission
	groupChapterPermissionRepo := repositories.NewGroupChapterPermissionRepository(
		StudentHandler.server.DB,
	)
	groupChapterPermissionRepo.GetGroupChapterPermissionByGroupID(
		&groupChapterPermission,
		*existUser.Student.GroupID,
	)

	studentAssignItemRepo := repositories.NewStudentAssignChapterItemRepository(
		StudentHandler.server.DB,
	)
	for _, item := range labClassInfos {
		var studentAssignChapterItems []models.StudentAssignmentChapterItem
		studentAssignItemRepo.GetStudentAssignChapter(
			&studentAssignChapterItems,
			userId,
			item.ChapterID,
		)
		if len(studentAssignChapterItems) < item.NoItems {
			maxIdxItem := 0
			if len(studentAssignChapterItems) > 0 {
				maxIdxItem = studentAssignChapterItems[len(studentAssignChapterItems)-1].ItemID
			}
			var chapter models.GroupChapterPermission
			for _, chapterPermission := range groupChapterPermission {
				if chapterPermission.ChapterID == item.ChapterID {
					chapter = chapterPermission
				}
			}
			studentAssignChapterItemService := studentassignmentchapteritem.NewStudentAssignmentChapterItem(
				StudentHandler.server.DB,
			)
			for i := maxIdxItem; i < item.NoItems; i++ {
				_, err := studentAssignChapterItemService.Create(
					userId,
					chapter.ChapterID,
					i+1,
					nil,
					item.FullMark,
					0,
					chapter.TimeStart,
					chapter.TimeEnd,
				)
				if err != nil {
					return responses.ErrorResponse(c, http.StatusInternalServerError, "Create Student Assigned Chapter Item Fail")
				}
			}
		}
	}

	var allStudentAssignChapterItems []models.StudentAssignmentChapterItem
	studentAssignItemRepo.GetAllStudentAssignChapter(&allStudentAssignChapterItems, userId)

	response := responses.NewGetAllChapter(
		groupChapterPermission,
		allStudentAssignChapterItems,
		labClassInfos,
	)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Chapter List
// @ID student-get-chapter-list
// @Tags Student
// @Accept json
// @Produce json
// @Success 200		{array}	responses.GetChapterListResponse
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/student/chapter_list [get]
func (StudentHandler *StudentHandler) GetChapterList(c echo.Context) error {
	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID
	var existUser models.User
	userRepo := repositories.NewUserRepository(StudentHandler.server.DB)
	userRepo.GetUserByUserID(&existUser, userId)
	if *existUser.Role != constants.Role.Student {
		return responses.ErrorResponse(c, http.StatusForbidden, "This User Not Student")
	}

	var labClassInfos []models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(StudentHandler.server.DB)
	labClassInfoRepo.GetAllLabClassInfos(&labClassInfos)

	var groupChapterPermission []models.GroupChapterPermission
	groupChapterPermissionRepo := repositories.NewGroupChapterPermissionRepository(
		StudentHandler.server.DB,
	)
	groupChapterPermissionRepo.GetGroupChapterPermissionByGroupID(
		&groupChapterPermission,
		*existUser.Student.GroupID,
	)

	var allGroupChapterItems []models.GroupAssignmentChapterItem
	groupChapterItemRepo := repositories.NewGroupAssignmentChapterItemRepository(
		StudentHandler.server.DB,
	)
	groupChapterItemRepo.GetAllGroupAssignmentChapterItemsByGroupId(
		&allGroupChapterItems,
		*existUser.Student.GroupID,
	)

	var allStudentAssignChapterItems []models.StudentAssignmentChapterItem
	studentAssignItemRepo := repositories.NewStudentAssignChapterItemRepository(
		StudentHandler.server.DB,
	)
	studentAssignItemRepo.GetAllStudentAssignChapter(&allStudentAssignChapterItems, userId)

	response := responses.NewGetChapterListResponse(
		groupChapterPermission,
		allGroupChapterItems,
		allStudentAssignChapterItems,
	)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Assigned Exercxise
// @ID student-get-assigned-exercise
// @Tags Student
// @Accept json
// @Produce json
// @Param chapter_idx query string false "chapter_idx"
// @Param item_id query string false "item_id"
// @Success 200		{object}	responses.StudentAssignmentItemResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/student/assigned_exercise [get]
func (StudentHandler *StudentHandler) GetStudentAssignedExercise(c echo.Context) error {
	chapterIdx := c.QueryParam("chapter_idx")
	itemId := c.QueryParam("item_id")

	chapterInt, err := strconv.Atoi(chapterIdx)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Can't Convert Chapter Index",
		)
	}

	itemInt, err := strconv.Atoi(itemId)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Convert Item ID")
	}

	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID

	var existUser models.User
	userRepo := repositories.NewUserRepository(StudentHandler.server.DB)
	userRepo.GetUserByUserID(&existUser, userId)

	if *existUser.Role != constants.Role.Student {
		return responses.ErrorResponse(c, http.StatusForbidden, "This User Not Student")
	}

	var labClassInfo models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(StudentHandler.server.DB)
	labClassInfoRepo.GetLabClassInfoByChapterIndex(&labClassInfo, chapterInt)

	var studentAssignChapterItems models.StudentAssignmentChapterItem
	studentAssignItemRepo := repositories.NewStudentAssignChapterItemRepository(
		StudentHandler.server.DB,
	)
	err = studentAssignItemRepo.GetStudentAssignChapterItem(
		&studentAssignChapterItems,
		existUser.UserID,
		labClassInfo.ChapterID,
		itemInt,
	)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			var groupChapterPermission models.GroupChapterPermission
			groupChapterPermissionRepo := repositories.NewGroupChapterPermissionRepository(StudentHandler.server.DB)
			groupChapterPermissionRepo.GetGroupChapterPermissionByPK(&groupChapterPermission, *existUser.Student.GroupID, labClassInfo.ChapterID)
			studentAssignChapterItemService := studentassignmentchapteritem.NewStudentAssignmentChapterItem(
				StudentHandler.server.DB,
			)
			createData, err := studentAssignChapterItemService.Create(
				userId,
				labClassInfo.ChapterID,
				itemInt,
				nil,
				labClassInfo.FullMark,
				0,
				groupChapterPermission.TimeStart,
				groupChapterPermission.TimeEnd,
			)
			if err != nil {
				return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			studentAssignChapterItems = *createData
		} else {
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	if studentAssignChapterItems.ExerciseID == nil {
		var selectItem []models.GroupChapterSelectedItem
		groupChapterSelectedItemRepo := repositories.NewGroupChapterSelectedItemRepository(
			StudentHandler.server.DB,
		)
		groupChapterSelectedItemRepo.GetSelectedItemByGroupChapterItemId(
			&selectItem,
			*existUser.Student.GroupID,
			labClassInfo.ChapterID,
			itemInt,
		)
		noSelectItem := len(selectItem)
		if noSelectItem < 1 {
			return responses.ErrorResponse(c, http.StatusBadRequest, "No Exercise Available")
		} else {
			studentAssigmItemService := studentassignmentchapteritem.NewStudentAssignmentChapterItem(StudentHandler.server.DB)
			if noSelectItem == 1 {
				studentAssigmItemService.UpdateAssignExercise(&studentAssignChapterItems, &selectItem[0].ExerciseID)
			} else {
				idx := rand.Intn(noSelectItem - 1)
				studentAssigmItemService.UpdateAssignExercise(&studentAssignChapterItems, &selectItem[idx].ExerciseID)
			}
			studentAssignItemRepo.GetStudentAssignChapterItem(&studentAssignChapterItems, existUser.UserID, labClassInfo.ChapterID, itemInt)
		}
	}

	var labExercise models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(StudentHandler.server.DB)
	labExerciseRepo.GetLabExerciseByID(*studentAssignChapterItems.ExerciseID, &labExercise)

	response := responses.NewGetStudentAssignmentItemResponse(labExercise)

	return responses.Response(c, http.StatusOK, response)
}
