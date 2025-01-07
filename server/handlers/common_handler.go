package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	minioclient "github.com/Project-IPCA/ipca-backend/minio_client"
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/pkg/utils"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
	userservice "github.com/Project-IPCA/ipca-backend/services/user"
)

type CommonHandler struct {
	server *s.Server
}

func NewCommonHandler(server *s.Server) *CommonHandler {
	return &CommonHandler{server: server}
}

// @Description Get User Info
// @ID common-get-user-info
// @Tags Common
// @Accept json
// @Produce json
// @Success 200		{object}	responses.UserInfoResponse
// @Security BearerAuth
// @Router			/api/common/user_info [get]
func (commonHandler *CommonHandler) GetUserInfo(c echo.Context) error {
	userRepository := repositories.NewUserRepository(commonHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var allDepts []models.Department
	deptRepository := repositories.NewDepartmentRepository(commonHandler.server.DB)
	deptRepository.GetAllDepts(&allDepts)

	var classSchedule models.ClassSchedule
	if *existUser.Role == constants.Role.Student {
		classScheduleRepo := repositories.NewClassScheduleRepository(commonHandler.server.DB)
		classScheduleRepo.GetClassSchedulePreloadByGroupID(&classSchedule, *existUser.Student.GroupID)
	}

	response := responses.NewUserInfoResponse(existUser, allDepts, &classSchedule)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Update User Info
// @ID common-update-user-info
// @Tags Common
// @Accept json
// @Produce json
// @Param params body	requests.UpdateUserInfoRequest	true	"User Info Request"
// @Success 200		{object}	responses.UserInfoResponse
// @Failure 400		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/common/user_info [put]
func (commonHandler *CommonHandler) UpdateUserInfo(c echo.Context) error {
	updateUserInfoReq := new(requests.UpdateUserInfoRequest)

	userRepository := repositories.NewUserRepository(commonHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if err := c.Bind(updateUserInfoReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Required Field")
	}

	if err := updateUserInfoReq.BasicUserInfo.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	if updateUserInfoReq.NewPassword != nil && updateUserInfoReq.ConfirmNewPassword != nil &&
		(*updateUserInfoReq.NewPassword != *updateUserInfoReq.ConfirmNewPassword) {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(existUser.Password),
		[]byte(updateUserInfoReq.CurrentPassword),
	) != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Password is not correct.",
		)
	}

	userService := userservice.NewUserService(commonHandler.server.DB)
	userService.UpdateUserInfo(&existUser, updateUserInfoReq)

	var allDepts []models.Department
	deptRepository := repositories.NewDepartmentRepository(commonHandler.server.DB)
	deptRepository.GetAllDepts(&allDepts)

	var classSchedule models.ClassSchedule
	if *existUser.Role == constants.Role.Student {
		classScheduleRepo := repositories.NewClassScheduleRepository(commonHandler.server.DB)
		classScheduleRepo.GetClassSchedulePreloadByGroupID(&classSchedule, *existUser.Student.GroupID)
	}

	response := responses.NewUserInfoResponse(existUser, allDepts, &classSchedule)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Keyword List
// @ID common-get-keyword-list
// @Tags Common
// @Accept json
// @Produce json
// @Param params body	requests.GetKeywordListRequest	true	"Get Keyword List Request"
// @Success 200		{object}	constants.ReceiveGetKeyWordData
// @Failure 400		{object}	responses.Error
// @Router			/api/common/get_keyword_list [post]
func (commonHandler *CommonHandler) GetKeywordList(c echo.Context) error {
	getKeywordListRequest := new(requests.GetKeywordListRequest)
	if err := c.Bind(getKeywordListRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := getKeywordListRequest.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}
	keywordList, err := utils.GetKeywordFromCode(getKeywordListRequest.Sourcecode)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Error While Running Sourcecode %s", err),
		)
	}
	return responses.Response(c, http.StatusOK, keywordList)
}

// @Description Keyword Check
// @ID common-keyword-check
// @Tags Common
// @Accept json
// @Produce json
// @Param params body	requests.CheckKeywordRequest	true	"Keyword Check"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Router			/api/common/keyword_check [post]
func (commonHandler *CommonHandler) KeywordCheck(c echo.Context) error {
	checkKeywordRequest := new(requests.CheckKeywordRequest)
	if err := c.Bind(checkKeywordRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := checkKeywordRequest.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}
	checkKeyword, err := utils.KeywordCheck(
		checkKeywordRequest.Sourcecode,
		checkKeywordRequest.ExerciseKeywordList,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Error While Running Sourcecode %s", err),
		)
	}
	return responses.Response(c, http.StatusOK, checkKeyword)
}

// @Description Get Student Submission
// @ID common-get-student-submission
// @Tags Common
// @Accept json
// @Produce json
// @Param stu_id query string false "stu_id"
// @Param chapter_idx query string false "chapter_idx"
// @Param item_id query string false "item_id"
// @Success 200		{array}		responses.StudentSubmssionResponse
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/common/student_submission [get]
func (commonHandle *CommonHandler) GetStudentSubmission(c echo.Context) error {
	stuId := c.QueryParam("stu_id")
	chapterIdx := c.QueryParam("chapter_idx")
	itemId := c.QueryParam("item_id")

	userRepository := repositories.NewUserRepository(commonHandle.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	stuUuid, err := uuid.Parse(stuId)
	if err != nil && *existUser.Role != constants.Role.Student {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Error While Parse Student ID %s", err),
		)
	}

	if *existUser.Role == constants.Role.Student {
		stuUuid = existUser.UserID
	}

	chapterInt, err := strconv.Atoi(chapterIdx)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Error While Convert Chapter Index String To Int %s", err),
		)
	}

	itemIdInt, err := strconv.Atoi(itemId)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Error While Convert Item Id To Int %s", err),
		)
	}

	studentAssignChapterItemRepo := repositories.NewStudentAssignChapterItemRepository(
		commonHandle.server.DB,
	)

	var labClassInfoData models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(commonHandle.server.DB)
	labClassInfoRepo.GetLabClassInfoByChapterIndex(&labClassInfoData, chapterInt)

	var assignItem models.StudentAssignmentChapterItem
	studentAssignChapterItemRepo.GetStudentAssignChapterItem(
		&assignItem,
		stuUuid,
		labClassInfoData.ChapterID,
		itemIdInt,
	)

	if assignItem.ExerciseID == nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"No Exercise Assigned To This Student.",
		)
	}

	exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(commonHandle.server.DB)
	var exerciseSubmissionList []models.ExerciseSubmission
	exerciseSubmissionRepo.GetStudentSubmission(
		stuUuid,
		*assignItem.ExerciseID,
		&exerciseSubmissionList,
	)

	response := responses.NewStudentSubmssionResponse(exerciseSubmissionList)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Upload User Profile
// @ID common-upload-user-profile
// @Tags Common
// @Accept  multipart/form-data
// @Produce json
// @Param file formData file true "Profile Image File"
// @Success 200 {object} responses.FileResponse
// @Failure 400 {object} responses.Error
// @Security BearerAuth
// @Router /api/common/user_profile [post]
func (commonHandler *CommonHandler) UploadUserProfile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Image Uploaded")
	}

	userRepository := repositories.NewUserRepository(commonHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if existUser.Role == &constants.Role.Student {
		var classSchedule models.ClassSchedule
		classScheduleRepo := repositories.NewClassScheduleRepository(commonHandler.server.DB)
		classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, *existUser.Student.GroupID)

		if !classSchedule.AllowUploadPic {
			return responses.ErrorResponse(c, http.StatusForbidden, "Not Allow To Upload Image")
		}
	}

	minioAction := minioclient.NewMinioAction(commonHandler.server.Minio)
	imageName, err := minioAction.UploadToMinio(
		file,
		commonHandler.server.Config.Minio.BucketProfile,
		false,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to Upload Image")
	}

	var imageUrl string
	if commonHandler.server.Config.Env.Enviroment == constants.EnviromentType.Develop {
		imageUrl = fmt.Sprintf(
			"http://localhost:%s/%s/%s",
			commonHandler.server.Config.Minio.Port,
			commonHandler.server.Config.Minio.BucketProfile,
			imageName,
		)
	} else if commonHandler.server.Config.Env.Enviroment == constants.EnviromentType.Production {
		imageUrl = fmt.Sprintf(
			"%s/%s/%s",
			commonHandler.server.Config.Minio.Endpoint,
			commonHandler.server.Config.Minio.BucketProfile,
			imageName,
		)
	}

	response := responses.NewFileResponse(imageName, imageUrl)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Departments
// @ID common-get-departments
// @Tags Common
// @Accept json
// @Produce json
// @Success 200 {array} responses.DepartmentResponse
// @Security BearerAuth
// @Router /api/common/departments [get]
func (commonHandler *CommonHandler) GetDepartments(c echo.Context) error {
	var depts []models.Department
	deptRepo := repositories.NewDepartmentRepository(commonHandler.server.DB)
	deptRepo.GetAllDepts(&depts)
	response := responses.NewDepartmentsResponse(depts)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Staffs
// @ID common-get-staffs
// @Tags Common
// @Accept json
// @Produce json
// @Success 200 {array} responses.SupervisorsResponse
// @Security BearerAuth
// @Router /api/common/staffs [get]
func (commonHandler *CommonHandler) GetStaffs(c echo.Context) error {
	var supervisors []models.Supervisor
	supervRepo := repositories.NewSupervisorRepository(commonHandler.server.DB)
	supervRepo.GetAllSupervisors(&supervisors)
	response := responses.NewSupervisorsResponse(supervisors)
	return responses.Response(c, http.StatusOK, response)
}
