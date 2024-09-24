package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	minioclient "github.com/Project-IPCA/ipca-backend/minio_client"
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/pkg/utils"

	"github.com/Project-IPCA/ipca-backend/rabbitmq_client"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
	classlabstaff "github.com/Project-IPCA/ipca-backend/services/class_lab_staff"
	classschedule "github.com/Project-IPCA/ipca-backend/services/class_schedule"
	groupassignmentchapteritem "github.com/Project-IPCA/ipca-backend/services/group_assignment_chapter_item"
	groupassignmentexercise "github.com/Project-IPCA/ipca-backend/services/group_assignment_exercise"
	groupchapterpermission "github.com/Project-IPCA/ipca-backend/services/group_chapter_permission"
	labexercise "github.com/Project-IPCA/ipca-backend/services/lab_exercise"
	"github.com/Project-IPCA/ipca-backend/services/student"
	"github.com/Project-IPCA/ipca-backend/services/token"
	"github.com/Project-IPCA/ipca-backend/services/user"
)

type SupervisorHandler struct {
	server *s.Server
}

func NewSupervisorHandler(server *s.Server) *SupervisorHandler {
	return &SupervisorHandler{server: server}
}

// @Description Add Students
// @ID supervisor-add-students
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.AddStudentsTextRequest	true	"Add Students Request"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/students [post]
func (supervisorHandler *SupervisorHandler) AddStudents(c echo.Context) error {
	addStudentsReq := new(requests.AddStudentsTextRequest)

	if err := c.Bind(addStudentsReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}

	if err := addStudentsReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	arrStudents := strings.Split(addStudentsReq.StudentsData, "\n")

	userService := user.NewUserService(supervisorHandler.server.DB)
	studentService := student.NewStudentService(supervisorHandler.server.DB)
	studentRepository := repositories.NewStudentRepository(supervisorHandler.server.DB)

	var wg sync.WaitGroup
	errChan := make(chan error, len(arrStudents))

	for _, item := range arrStudents {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			data := strings.Split(item, " ")
			kmitlId := data[1]
			firstName := data[2]
			lastName := data[3]
			gender := constants.Gender.Other

			if strings.Contains(firstName, constants.Prefix.Mr) {
				firstName = strings.TrimPrefix(firstName, constants.Prefix.Mr)
				gender = constants.Gender.Male
			} else if strings.Contains(firstName, constants.Prefix.Miss) {
				firstName = strings.TrimPrefix(firstName, constants.Prefix.Miss)
				gender = constants.Gender.Female
			}
			existUserStudent := models.Student{}
			studentRepository.GetUserByKmitlID(&existUserStudent, kmitlId)
			if existUserStudent.KmitlID == kmitlId {
				errChan <- fmt.Errorf("User Student with ID %s is Already Exist", kmitlId)
				return
			}

			userId, err := userService.CreateQuick(
				kmitlId,
				kmitlId,
				firstName,
				lastName,
				gender,
				constants.Role.Student,
			)
			if err != nil {
				errChan <- fmt.Errorf("Failed to create user for student %s: %v", kmitlId, err)
				return
			}

			err = studentService.Create(userId, kmitlId, &addStudentsReq.GroupID)
			if err != nil {
				errChan <- fmt.Errorf("Failed to create user student for %s: %v", kmitlId, err)
				return
			}
		}(item)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
	}
	return responses.MessageResponse(c, http.StatusCreated, "Add Student Successful")
}

// @Description Create Group
// @ID supervisor-create-group
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.CreateGroupRequest	true	"Create Group Request"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/group [post]
func (supervisorHandler *SupervisorHandler) CreateGroup(c echo.Context) error {
	createGroupReq := new(requests.CreateGroupRequest)
	if err := c.Bind(createGroupReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := createGroupReq.BasicGroup.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}
	supervisorId := existUser.UserID

	existGroup := models.ClassSchedule{}
	classScheduleRepository := repositories.NewClassScheduleRepository(
		supervisorHandler.server.DB,
	)
	classScheduleRepository.GetClassScheduleByNumber(&existGroup, *createGroupReq.Number)

	classScheduleService := classschedule.NewClassScheduleService(supervisorHandler.server.DB)
	groupId, _ := classScheduleService.Create(createGroupReq, &supervisorId)

	var existLabExercises []models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseOrderByChapterLevel(&existLabExercises)

	var existGroupAssignmentExercise []models.GroupAssignmentExercise
	groupAssignmentExerciseRepo := repositories.NewGroupAssignmentExerciseRepository(
		supervisorHandler.server.DB,
	)
	groupAssignmentExerciseRepo.GetGroupAssignmnetExercisesByGroupID(
		&existGroupAssignmentExercise,
		groupId,
	)
	var exerciseList []uuid.UUID
	for _, exercise := range existGroupAssignmentExercise {
		exerciseList = append(exerciseList, exercise.ExerciseID)
	}
	groupAssignmentExerciseService := groupassignmentexercise.NewGroupAssignmentExerciseService(
		supervisorHandler.server.DB,
	)
	if len(existLabExercises) > len(existGroupAssignmentExercise) {
		for _, labExercise := range existLabExercises {
			if !utils.ContainsUUID(exerciseList, labExercise.ExerciseID) {
				groupAssignmentExerciseService.Create(groupId, labExercise.ExerciseID, true)
			}
		}
	}

	existClassSchedule := models.ClassSchedule{}
	classScheduleRepository.GetClassScheduleByGroupID(&existClassSchedule, groupId)

	var labClassInfos []models.LabClassInfo
	labClassInfoRepository := repositories.NewLabClassInfoRepository(supervisorHandler.server.DB)
	labClassInfoRepository.GetAllLabClassInfos(&labClassInfos)
	groupAssignmentChapterItemRepository := repositories.NewGroupAssignmentChapterItemRepository(
		supervisorHandler.server.DB,
	)
	groupAssignmentChapterItemService := groupassignmentchapteritem.NewGroupAssignmentChapterItemService(
		supervisorHandler.server.DB,
	)
	for _, labClassInfo := range labClassInfos {
		for itemIndex := 1; itemIndex <= 5; itemIndex++ {
			gaci := models.GroupAssignmentChapterItem{}
			groupAssignmentChapterItemRepository.GetGroupAssignmentChapterItemsByPK(
				&gaci,
				groupId,
				labClassInfo.ChapterID,
				int64(itemIndex),
			)
			if gaci.ItemID != int64(itemIndex) {
				groupAssignmentChapterItemService.Create(
					groupId,
					labClassInfo.ChapterID,
					int64(itemIndex),
					2,
					existClassSchedule.TimeStart,
					existClassSchedule.TimeEnd,
					constants.GroupAssignmentChapterItemStatus.Closed,
				)
			}
		}
	}

	groupChapterPermRepo := repositories.NewGroupChapterPermissionRepository(
		supervisorHandler.server.DB,
	)
	groupChapterPermService := groupchapterpermission.NewGroupChapterPermissionService(
		supervisorHandler.server.DB,
	)
	for _, labClassInfo := range labClassInfos {
		gcp := models.GroupChapterPermission{}
		groupChapterPermRepo.GetGroupChapterPermissionByPK(&gcp, groupId, labClassInfo.ChapterID)
		if gcp.ChapterID != labClassInfo.ChapterID {
			groupChapterPermService.Create(
				groupId,
				labClassInfo.ChapterID,
				*existClassSchedule.TimeStart,
				*existClassSchedule.TimeEnd,
				false,
				false,
				constants.GroupChapterPermStatus.Na,
			)
		}
	}

	classLabStaffService := classlabstaff.NewClassLabStaffService(supervisorHandler.server.DB)
	for _, item := range createGroupReq.Staffs {
		classLabStaffService.Create(groupId, item.StaffID)
	}

	return responses.MessageResponse(c, http.StatusOK, "Create Group Successful.")
}

// @Description Get All Available Group
// @ID supervisor-get-all-available-group
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param instructorId query string false "instructorId"
// @Param staffIds query []string false "staffIds" collectionFormat(multi)
// @Param year query string false "year"
// @Param semester query string false "semester"
// @Param day query string false "day"
// @Param page query string false "Page"
// @Param pageSize query string false "Page Size"
// @Success 200		{array}	responses.AvailableGroupsResponse
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/available_groups [get]
func (supervisorHandler *SupervisorHandler) GetAllAvailableGroups(c echo.Context) error {
	instructorId := c.QueryParam("instructorId")
	staffIds := c.QueryParams()["staffIds"]
	year := c.QueryParam("year")
	semester := c.QueryParam("semester")
	day := c.QueryParam("day")
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var existClassSchedules []models.ClassSchedule
	classSceduleR := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classSceduleR.GetAllClassSchedulesByQuery(
		&existClassSchedules,
		instructorId,
		staffIds,
		year,
		semester,
		day,
		page,
		pageSize,
	)

	var allClassSchedules []models.ClassSchedule
	classSceduleR.GetAllClassSchedules(&allClassSchedules)

	var allSupervisors []models.Supervisor
	supervisorR := repositories.NewSupervisorRepository(supervisorHandler.server.DB)
	supervisorR.GetAllSupervisors(&allSupervisors)
	response := responses.NewClassSchedulesResponse(
		existClassSchedules,
		allClassSchedules,
		allSupervisors,
	)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get My Groups
// @ID supervisor-get-my-groups
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param year query string false "Year"
// @Param page query string false "Page"
// @Param pageSize query string false "Page Size"
// @Success 200		{array}	responses.MyGroupResponse
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/my_groups [get]
func (supervisorHandler *SupervisorHandler) GetMyGroups(c echo.Context) error {
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")
	year := c.QueryParam("year")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var existClassSchedules []models.ClassSchedule
	classScheduleR := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleR.GetMyClassSchedulesByQuery(
		&existClassSchedules,
		existUser.UserID,
		year,
		page,
		pageSize,
	)

	var allClassSchedules []models.ClassSchedule
	classScheduleR.GetAllClassSchedules(&allClassSchedules)

	response := responses.NewMyClassSchedulesResponse(existClassSchedules, allClassSchedules)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Group Info By Group ID
// @ID supervisor-get-group-info-by-group-id
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Success 200		{object}	responses.ClassScheduleInfoResponse
// @Failure 403		{object}	responses.Error
// @Failure 404		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/group/{group_id} [get]
func (supervisorHandler *SupervisorHandler) GetGroupInfoByGroupID(c echo.Context) error {
	groupIdStr := c.Param("group_id")
	groupId, err := uuid.Parse(groupIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	existClassSchedule := models.ClassSchedule{}
	classScheduleR := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleR.GetClassSchedulePreloadByGroupID(&existClassSchedule, groupId)

	if *existClassSchedule.SupervisorID != existUser.UserID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	if existClassSchedule.GroupID != groupId {
		return responses.ErrorResponse(c, http.StatusNotFound, "Not found group.")
	}

	response := responses.NewClassScheduleInfoResponse(existClassSchedule)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Reset Student Password
// @ID supervisor-reset-student-password
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param stu_id path string true "Stu ID"
// @Success 200		{object}	responses.Data
// @Failure 403		{object}	responses.Error
// @Failure 404		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/reset_student_password/{stu_id} [put]
func (supervisorHandler *SupervisorHandler) ResetStudentPassword(c echo.Context) error {
	stuIdStr := c.Param("stu_id")
	stuId, err := uuid.Parse(stuIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	existStudent := models.User{}
	userRepository.GetUserByUserID(&existStudent, stuId)
	if existStudent.UserID != stuId {
		return responses.ErrorResponse(c, http.StatusNotFound, "User Not found.")
	}

	userService := user.NewUserService(supervisorHandler.server.DB)
	userService.ResetUserStudentPassword(&existStudent, existStudent.Student.KmitlID)

	return responses.MessageResponse(c, http.StatusOK, "Reset Student Password.")
}

// @Description Get My Group Info By Group ID
// @ID supervisor-get-my-group-info-by-group-id
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Success 200		{object}	responses.MyGroupInfoResponse
// @Failure 403		{object}	responses.Error
// @Failure 404		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/my_group_info/{group_id} [get]
func (supervisorHandler *SupervisorHandler) GetMyGroupInfo(c echo.Context) error {
	groupIdStr := c.Param("group_id")
	groupId, err := uuid.Parse(groupIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	existClassSchedule := models.ClassSchedule{}
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassSchedulePreloadByGroupID(&existClassSchedule, groupId)
	if existClassSchedule.GroupID != groupId {
		return responses.ErrorResponse(c, http.StatusNotFound, "Not Found Class Schedule.")
	}

	if existUser.UserID != *existClassSchedule.SupervisorID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var existDepts []models.Department
	deptRepo := repositories.NewDepartmentRepository(supervisorHandler.server.DB)
	deptRepo.GetAllDepts(&existDepts)

	var existStaffs []models.Supervisor
	supervisorRepo := repositories.NewSupervisorRepository(supervisorHandler.server.DB)
	supervisorRepo.GetAllSupervisors(&existStaffs)

	response := responses.NewMyClassScheduleInfoResponse(
		existClassSchedule,
		existDepts,
		existStaffs,
	)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Update My Group Info
// @ID supervisor-update-my-group-info
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.UpdateMyGroupRequest	true	"Update My Group Request"
// @Param group_id path string true "Group ID"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 404		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/my_group_info/{group_id} [put]
func (supervisorHandler *SupervisorHandler) UpdateMyGroupInfo(c echo.Context) error {
	updateMyGroupReq := new(requests.UpdateMyGroupRequest)
	if err := c.Bind(updateMyGroupReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := updateMyGroupReq.BasicGroup.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	groupIdStr := c.Param("group_id")
	groupId, err := uuid.Parse(groupIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	if !utils.IsRoleSupervisor(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	existClassSchedule := models.ClassSchedule{}
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&existClassSchedule, groupId)
	if existClassSchedule.GroupID != groupId {
		return responses.ErrorResponse(c, http.StatusNotFound, "Not Found Class Schedule.")
	}

	if existUser.UserID != *existClassSchedule.SupervisorID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	classScheduleService := classschedule.NewClassScheduleService(supervisorHandler.server.DB)
	classScheduleService.UpdateMyGroup(
		&existClassSchedule,
		updateMyGroupReq,
	)

	var existClassLabStaff []models.ClassLabStaff
	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	classLabStaffRepo.GetClassLabStaffByGroupID(&existClassLabStaff, groupId)

	classLabStaffService := classlabstaff.NewClassLabStaffService(supervisorHandler.server.DB)
	classLabStaffService.DeleteAll(&existClassLabStaff)
	for _, staff := range updateMyGroupReq.Staffs {
		classLabStaffService.Create(groupId, staff.StaffID)
	}

	return responses.MessageResponse(c, http.StatusOK, "Update Group Info.")
}

// @Description Create Exercise
// @ID supervisor-create-exercise
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.CreateLabExerciseRequest	true	"Creaet Lab Exercise Request"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/exercise [post]
func (supervisorHandler *SupervisorHandler) CreateExercise(c echo.Context) error {
	createLabExerciseReq := new(requests.CreateLabExerciseRequest)
	if err := c.Bind(createLabExerciseReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := createLabExerciseReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}
	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID

	existUser := models.User{}
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	userRepository.GetUserByUserID(&existUser, userId)
	if *existUser.Role != constants.Role.Supervisor {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}
	labExerciseService := labexercise.NewLabExerciseService(supervisorHandler.server.DB)
	exerciseId,err := labExerciseService.CreateWithoutSourceCode(createLabExerciseReq, &existUser.UserID, existUser.Username)
	if(err!= nil){
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Create Exercise Fail")
	}

	filename := fmt.Sprintf("exercise_"+exerciseId.String()+"*.py")
	tempFile,err := utils.CreateTempFile(filename,createLabExerciseReq.Sourcecode)
	if(err!=nil){
		return responses.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Create Temp File Fail %s",err))
	}
	defer os.Remove(tempFile.Name())
	
	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	uploadFileName,err := minioAction.UploadToMinio(
		tempFile,
		supervisorHandler.server.Config.Minio.BucketSupervisorCode,
		false,
	)
	if(err!=nil){
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.UpdateLabExerciseSourcecode(exerciseId.String(),uploadFileName)
	var labExerciseData models.LabExercise
	labExerciseRepo.GetLabExerciseByID(exerciseId.String(),&labExerciseData)
	
	return responses.Response(c,http.StatusOK,labExerciseData)
}

//
// func (supervisorHandler *SupervisorHandler) UpdateExerciseTestcases(c echo.Context) error {
// 	userJwt := c.Get("user").(*jwt.Token)
// 	claims := userJwt.Claims.(*token.JwtCustomClaims)
// 	userId := claims.UserID
//
// 	existUser := models.User{}
// 	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
// 	userRepository.GetUserByUserID(&existUser, userId)
// 	if *existUser.Role != constants.Role.Supervisor {
// 		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
// 	}
// 	return nil
// }

// @Description Save Exercise Testcase
// @ID supervisor-save-excercise-testcase
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.SaveExerciseTestcaseRequest	true	"Save Exercise Testcase"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/save_exercise_testcase [post]
func (supervisorHandler *SupervisorHandler) SaveExerciseTestcase(c echo.Context) error{
	saveExerciseTesetcaseReq := new(requests.SaveExerciseTestcaseRequest)
	if err:= c.Bind(saveExerciseTesetcaseReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := saveExerciseTesetcaseReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID

	var labExercise models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseByID(saveExerciseTesetcaseReq.ExerciseID,&labExercise)

	if(labExercise.CreatedBy != nil && *labExercise.CreatedBy != userId){
		return responses.ErrorResponse(c, http.StatusForbidden, "Not Have Permission To Do This!")
	}

	exerciseTestcaseRepo := repositories.NewExerciseTestcaseRepository(supervisorHandler.server.DB)
	if(len(saveExerciseTesetcaseReq.RemoveList)>0){
		exerciseTestcaseRepo.DeleteExerciseTetscaseID(saveExerciseTesetcaseReq.RemoveList)
	}

	for i := range saveExerciseTesetcaseReq.TestCaseList {
		saveExerciseTesetcaseReq.TestCaseList[i].IsReady = "no"
		if(saveExerciseTesetcaseReq.TestCaseList[i].TestcaseID == nil){
			fmt.Println("nil condition")
			exerciseUuid,_ := uuid.Parse(string(saveExerciseTesetcaseReq.TestCaseList[i].ExerciseID))
			convertSaveTestcase := models.ExerciseTestcase{
				ExerciseID: exerciseUuid,
				IsReady: saveExerciseTesetcaseReq.TestCaseList[i].IsReady,
				TestcaseContent: saveExerciseTesetcaseReq.TestCaseList[i].TestcaseContent,
				IsActive: &saveExerciseTesetcaseReq.TestCaseList[i].IsActive,
				IsShowStudent: &saveExerciseTesetcaseReq.TestCaseList[i].IsShowStudent,
				TestcaseNote: &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseNote,
				TestcaseOutput: &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseOutput,
				TestcaseError: &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseError,
			}
			saveId := exerciseTestcaseRepo.UpsertExerciseTestcaseID(convertSaveTestcase)
			saveIdString := saveId.String()
			if(saveId != nil){
				saveExerciseTesetcaseReq.TestCaseList[i].TestcaseID = &saveIdString
			}
		}else{
			fmt.Println("not nil")
			exerciseUuid,_ := uuid.Parse(string(saveExerciseTesetcaseReq.TestCaseList[i].ExerciseID))
			testcaseUuid,_ := uuid.Parse(string(*saveExerciseTesetcaseReq.TestCaseList[i].TestcaseID))
			convertSaveTestcase := models.ExerciseTestcase{
				ExerciseID: exerciseUuid,
				TestcaseID: &testcaseUuid,
				IsReady: saveExerciseTesetcaseReq.TestCaseList[i].IsReady,
				TestcaseContent: saveExerciseTesetcaseReq.TestCaseList[i].TestcaseContent,
				IsActive: &saveExerciseTesetcaseReq.TestCaseList[i].IsActive,
				IsShowStudent: &saveExerciseTesetcaseReq.TestCaseList[i].IsShowStudent,
				TestcaseNote: &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseNote,
				TestcaseOutput: &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseOutput,
				TestcaseError: &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseError,
			}
			exerciseTestcaseRepo.UpsertExerciseTestcaseID(convertSaveTestcase)
		}
	}

	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	minioFile,err := minioAction.GetFromMinio(supervisorHandler.server.Config.Minio.BucketSupervisorCode,*labExercise.Sourcecode)
	if(err!=nil){
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Get Data From Minio Object Fail")
	}
	defer minioFile.Close()

	var sourcecode strings.Builder
	_, err = io.Copy(&sourcecode, minioFile)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Copy Data From Minio Object Fail")
	}
	
	rabbit := rabbitmq_client.NewRabbitMQAction(supervisorHandler.server.RabitMQ,supervisorHandler.server.Config)
	message := requests.AddTestcaseRabitMessage{
		JobId: saveExerciseTesetcaseReq.JobID,
		JobType: "upsert-testcase",
		ExerciseId: saveExerciseTesetcaseReq.ExerciseID,
		TestcaseList: saveExerciseTesetcaseReq.TestCaseList,
		Sourcecode: sourcecode.String(),
	}
	err = rabbit.SendQueue(message)
	if(err!=nil){
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error While Send Queue RabbitMQ")
	}
	return responses.MessageResponse(c,http.StatusOK,"Testcases Are Being Run")
}