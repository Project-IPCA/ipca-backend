package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

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

	// TODO add validate group permission and assign chapter item

	chaperUuid, err := uuid.Parse(exerciseSubmitReq.ChapterID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "GroupID Is Not UUID")
	}

	var submissionList []models.ExerciseSubmission
	exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(studentHandler.server.DB)
	exerciseSubmissionRepo.GetStudentSubmission(existUser.UserID, chaperUuid, &submissionList)

	attemps := len(submissionList) + 1
	// TODO Improve filename and dir path to env
	filename := fmt.Sprintf("%s-%04d.py", existUser.Username, attemps)
	err = utils.CreateSourcecode("./bucket/student", filename, exerciseSubmitReq.Sourcecode)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	exerciseUuid, err := uuid.Parse(exerciseSubmitReq.ItemId)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "ItemID Is Not UUID")
	}

	exerciseSubmissionService := exercisesubmission.NewExerciseSubmissionService(
		studentHandler.server.DB,
	)
	isInfLoop := false
	submissionId, err := exerciseSubmissionService.Create(
		existUser.UserID,
		exerciseUuid,
		filename,
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
	exerciseTestcaseRepo.GetTestcasesByExerciseID(exerciseUuid, &testcaseList)

	filterTestcase := make([]models.ExerciseTestcase, 0)
	for _, testcase := range testcaseList {
		if testcase.IsActive != nil && *testcase.IsActive {
			filterTestcase = append(filterTestcase, testcase)
		}
	}

	labExerciseRepo := repositories.NewLabExerciseRepository(studentHandler.server.DB)
	var labExercise models.LabExercise
	labExerciseRepo.GetLabExerciseByID(exerciseSubmitReq.ItemId, &labExercise)

	logAction := models.LogExerciseSubmissionAction{
		StuId:              existUser.UserID,
		JobId:              exerciseSubmitReq.JobID,
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
		JobId:        exerciseSubmitReq.JobID,
		JobType:      "exercise-submit",
		LogData:      logData,
		SubmissionId: submissionId,
		SourceCode:   exerciseSubmitReq.Sourcecode,
		TestCaseList: filterTestcase,
	}

	rabbit := rabbitmq_client.NewRabbitMQAction(
		studentHandler.server.RabitMQ,
		studentHandler.server.Config,
	)
	fmt.Println(rabbitMessage)
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
		"exercise-submit",
		string(logActionString),
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Save Activity Log")
	}

	redis := redis_client.NewRedisAction(studentHandler.server.Redis)
	redis.PublishMessage(fmt.Sprintf("logs:%s", existUser.Student.GroupID), insertLog)

	return responses.MessageResponse(c, http.StatusOK, "Submission are being run")
}

// @Description Get Lab Exercise By ID
// @ID get-lab-exercise-by-id
// @Tags Student
// @Accept json
// @Produce json
// @Param exercise_id path string true "Exercise ID"
// @Success 200		{array}	responses.Data
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/student/exercise/{exercise_id} [get]
func (studentHandler *StudentHandler) GetLabExerciseByID(c echo.Context) error {
	exerciseId := c.Param("exercise_id")

	existExercise := models.LabExercise{}
	exerciseRepo := repositories.NewLabExerciseRepository(studentHandler.server.DB)
	exerciseRepo.GetLabExerciseByID(exerciseId, &existExercise)
	return responses.Response(c, http.StatusOK, existExercise)
}
