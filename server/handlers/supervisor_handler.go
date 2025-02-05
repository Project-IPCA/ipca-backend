package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

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
	classlabstaff "github.com/Project-IPCA/ipca-backend/services/class_lab_staff"
	classschedule "github.com/Project-IPCA/ipca-backend/services/class_schedule"
	"github.com/Project-IPCA/ipca-backend/services/department"
	"github.com/Project-IPCA/ipca-backend/services/executive"
	exercisesubmission "github.com/Project-IPCA/ipca-backend/services/exercise_submission"
	exercisetestcase "github.com/Project-IPCA/ipca-backend/services/exercise_testcase"
	groupassignmentchapteritem "github.com/Project-IPCA/ipca-backend/services/group_assignment_chapter_item"
	groupassignmentexercise "github.com/Project-IPCA/ipca-backend/services/group_assignment_exercise"
	groupchapterpermission "github.com/Project-IPCA/ipca-backend/services/group_chapter_permission"
	groupchapterselecteditem "github.com/Project-IPCA/ipca-backend/services/group_chapter_selected_item"
	labexercise "github.com/Project-IPCA/ipca-backend/services/lab_exercise"
	rolepermission "github.com/Project-IPCA/ipca-backend/services/role_permission"
	"github.com/Project-IPCA/ipca-backend/services/student"
	studentassignmentchapteritem "github.com/Project-IPCA/ipca-backend/services/student_assignment_chapter_item"
	"github.com/Project-IPCA/ipca-backend/services/supervisor"
	"github.com/Project-IPCA/ipca-backend/services/ta"
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

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, addStudentsReq.GroupID)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(addStudentsReq.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(addStudentsReq.GroupID, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	arrStudents := strings.Split(addStudentsReq.StudentsData, "\n")

	userService := user.NewUserService(supervisorHandler.server.DB)
	studentService := student.NewStudentService(supervisorHandler.server.DB)
	studentRepository := repositories.NewStudentRepository(supervisorHandler.server.DB)

	var wg sync.WaitGroup
	errChan := make(chan error, len(arrStudents))

	for index, item := range arrStudents {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			data := strings.Fields(item)
			if len(data) == 0 || len(data) != 4 {
				errChan <- fmt.Errorf("Row %d Data Invalid", index+1)
				return
			}
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
	errList := make([]string, 0)
	for err := range errChan {
		if err != nil {
			errList = append(errList, err.Error())
		}
	}
	if len(errList) > 0 {
		errString := strings.Join(errList, "\n")
		return responses.ErrorResponse(c, http.StatusBadRequest, errString)
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
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)
		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var supervisorId uuid.UUID
	if createGroupReq.SupervisorId != nil && *createGroupReq.SupervisorId != uuid.Nil &&
		*existUser.Role != constants.Role.Supervisor {
		supervisorRepo := repositories.NewSupervisorRepository(supervisorHandler.server.DB)
		if supervisorRepo.CheckValidSuperID(*createGroupReq.SupervisorId) {
			supervisorId = *createGroupReq.SupervisorId
		} else {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Supervisor ID.")
		}
	} else {
		supervisorId = existUser.UserID
	}

	existGroup := models.ClassSchedule{}
	classScheduleRepository := repositories.NewClassScheduleRepository(
		supervisorHandler.server.DB,
	)
	classScheduleRepository.GetClassScheduleByNumber(&existGroup, *createGroupReq.Number)

	classScheduleService := classschedule.NewClassScheduleService(supervisorHandler.server.DB)
	groupId, err := classScheduleService.Create(createGroupReq, &supervisorId)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Error While Create Group.",
		)
	}

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
		fmt.Println(createGroupReq.Staffs)
		classLabStaffService.Create(groupId, item.StaffID)
	}

	return responses.MessageResponse(c, http.StatusOK, "Create Group Successful.")
}

// @Description Delete Group
// @ID supervisor-delete-group
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/group/{group_id} [delete]
func (supervisorHandler *SupervisorHandler) DeleteGroup(c echo.Context) error {
	groupIdStr := c.Param("group_id")
	groupId, err := uuid.Parse(groupIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	classscheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	var classSchedule models.ClassSchedule
	classscheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
	}

	var studentList []models.Student
	studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)
	studentRepo.GetStudentInGroupID(&studentList, groupId)

	userService := user.NewUserService(supervisorHandler.server.DB)
	var wg sync.WaitGroup
	errChan := make(chan error, len(studentList))

	for _, student := range studentList {
		wg.Add(1)
		go func(student models.Student) {
			defer wg.Done()
			var exerciseSubmission []models.ExerciseSubmission
			exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(
				supervisorHandler.server.DB,
			)
			exerciseSubmissionRepo.GetSubmissionByStudentID(student.StuID, &exerciseSubmission)

			minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
			for _, submission := range exerciseSubmission {
				minioAction.DeleteFileInMinio(
					supervisorHandler.server.Config.Minio.BucketStudentCode,
					submission.SourcecodeFilename,
				)
			}

			err = userService.Delete(student.User)
			if err != nil {
				errChan <- fmt.Errorf("Failed to delete student %s", student.KmitlID)
				return
			}
		}(student)
	}

	wg.Wait()
	close(errChan)

	errList := make([]string, 0)
	for err := range errChan {
		if err != nil {
			errList = append(errList, err.Error())
		}
	}
	if len(errList) > 0 {
		errString := strings.Join(errList, "\n")
		return responses.ErrorResponse(c, http.StatusInternalServerError, errString)
	}

	classScheduleService := classschedule.NewClassScheduleService(supervisorHandler.server.DB)
	err = classScheduleService.Delete(&classSchedule)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Delete Group.")
	}

	return responses.MessageResponse(c, http.StatusOK, "Delete Group Successful.")
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
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateAdminRole(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var existClassSchedules []models.ClassSchedule
	classSceduleR := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	totalClassSchedules := classSceduleR.GetAllClassSchedulesByQuery(
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

	var userAdmin []models.User
	userRepo := repositories.NewUserRepository(supervisorHandler.server.DB)
	userRepo.GetUserAdminRole(&userAdmin)
	response := responses.NewClassSchedulesResponse(
		existClassSchedules,
		allClassSchedules,
		userAdmin,
		page,
		pageSize,
		totalClassSchedules,
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
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateAdminRole(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var existClassSchedules []models.ClassSchedule
	classScheduleR := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	totalClassSchedules := classScheduleR.GetMyClassSchedulesByQuery(
		&existClassSchedules,
		existUser.UserID,
		year,
		page,
		pageSize,
	)

	var allClassSchedules []models.ClassSchedule
	classScheduleR.GetAllClassSchedules(&allClassSchedules)

	response := responses.NewMyClassSchedulesResponse(
		existClassSchedules,
		allClassSchedules,
		page,
		pageSize,
		int(totalClassSchedules),
	)
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

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassSchedulePreloadByGroupID(&classSchedule, groupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if *existUser.Role != constants.Role.Beyonder &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	response := responses.NewClassScheduleInfoResponse(classSchedule)
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

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	existStudent := models.User{}
	userRepository.GetUserByUserID(&existStudent, stuId)
	if existStudent.UserID != stuId {
		return responses.ErrorResponse(c, http.StatusNotFound, "User Not found.")
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, *existStudent.Student.GroupID)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(*existStudent.Student.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			*existStudent.Student.GroupID,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
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

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassSchedulePreloadByGroupID(&classSchedule, groupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	response := responses.NewMyClassScheduleInfoResponse(
		classSchedule,
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

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	classScheduleService := classschedule.NewClassScheduleService(supervisorHandler.server.DB)
	classScheduleService.UpdateMyGroup(
		&classSchedule,
		updateMyGroupReq,
	)

	var existClassLabStaff []models.ClassLabStaff
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

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.ExerciseAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	labExerciseService := labexercise.NewLabExerciseService(supervisorHandler.server.DB)
	exerciseId, err := labExerciseService.CreateWithoutSourceCode(
		createLabExerciseReq,
		&existUser.UserID,
		existUser.Username,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Create Exercise Fail")
	}

	filename := fmt.Sprintf("exercise_" + exerciseId.String() + "*.py")
	tempFile, err := utils.CreateTempFile(filename, createLabExerciseReq.Sourcecode)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Create Temp File Fail %s", err),
		)
	}
	defer os.Remove(tempFile.Name())

	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	uploadFileName, err := minioAction.UploadToMinio(
		tempFile,
		supervisorHandler.server.Config.Minio.BucketSupervisorCode,
		false,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	var labExerciseData models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseByID(exerciseId, &labExerciseData)
	labExerciseService.UpdateLabExerciseSourcecode(&labExerciseData, uploadFileName)

	return responses.Response(c, http.StatusOK, labExerciseData)
}

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
func (supervisorHandler *SupervisorHandler) SaveExerciseTestcase(c echo.Context) error {
	saveExerciseTesetcaseReq := new(requests.SaveExerciseTestcaseRequest)
	if err := c.Bind(saveExerciseTesetcaseReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := saveExerciseTesetcaseReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.ExerciseAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var labExercise models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseByID(saveExerciseTesetcaseReq.ExerciseID, &labExercise)

	exerciseTestcaseRepo := repositories.NewExerciseTestcaseRepository(supervisorHandler.server.DB)
	if len(saveExerciseTesetcaseReq.RemoveList) > 0 {
		exerciseTestcaseRepo.DeleteExerciseTetscaseID(saveExerciseTesetcaseReq.RemoveList)
	}

	for i := range saveExerciseTesetcaseReq.TestCaseList {
		saveExerciseTesetcaseReq.TestCaseList[i].IsReady = "no"
		if saveExerciseTesetcaseReq.TestCaseList[i].TestcaseID == nil {
			fmt.Println("nil condition")
			convertSaveTestcase := models.ExerciseTestcase{
				ExerciseID:      saveExerciseTesetcaseReq.TestCaseList[i].ExerciseID,
				IsReady:         saveExerciseTesetcaseReq.TestCaseList[i].IsReady,
				TestcaseContent: saveExerciseTesetcaseReq.TestCaseList[i].TestcaseContent,
				IsActive:        &saveExerciseTesetcaseReq.TestCaseList[i].IsActive,
				IsShowStudent:   &saveExerciseTesetcaseReq.TestCaseList[i].IsShowStudent,
				TestcaseNote:    &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseNote,
				TestcaseOutput:  &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseOutput,
				TestcaseError:   &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseError,
			}
			saveId := exerciseTestcaseRepo.UpsertExerciseTestcaseID(convertSaveTestcase)
			if saveId != nil {
				saveExerciseTesetcaseReq.TestCaseList[i].TestcaseID = saveId
			}
		} else {
			fmt.Println("not nil")
			convertSaveTestcase := models.ExerciseTestcase{
				ExerciseID:      saveExerciseTesetcaseReq.TestCaseList[i].ExerciseID,
				TestcaseID:      saveExerciseTesetcaseReq.TestCaseList[i].TestcaseID,
				IsReady:         saveExerciseTesetcaseReq.TestCaseList[i].IsReady,
				TestcaseContent: saveExerciseTesetcaseReq.TestCaseList[i].TestcaseContent,
				IsActive:        &saveExerciseTesetcaseReq.TestCaseList[i].IsActive,
				IsShowStudent:   &saveExerciseTesetcaseReq.TestCaseList[i].IsShowStudent,
				TestcaseNote:    &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseNote,
				TestcaseOutput:  &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseOutput,
				TestcaseError:   &saveExerciseTesetcaseReq.TestCaseList[i].TestcaseError,
			}
			exerciseTestcaseRepo.UpsertExerciseTestcaseID(convertSaveTestcase)
		}
	}

	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	minioFile, err := minioAction.GetFromMinio(
		supervisorHandler.server.Config.Minio.BucketSupervisorCode,
		*labExercise.Sourcecode,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Get Data From Minio Object Fail",
		)
	}
	defer minioFile.Close()

	var sourcecode strings.Builder
	_, err = io.Copy(&sourcecode, minioFile)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Copy Data From Minio Object Fail",
		)
	}

	rabbit := rabbitmq_client.NewRabbitMQAction(
		supervisorHandler.server.RabitMQ,
		supervisorHandler.server.Config,
	)
	message := requests.AddTestcaseRabitMessage{
		JobId:        saveExerciseTesetcaseReq.JobID,
		JobType:      "upsert-testcase",
		ExerciseId:   saveExerciseTesetcaseReq.ExerciseID,
		TestcaseList: saveExerciseTesetcaseReq.TestCaseList,
		Sourcecode:   sourcecode.String(),
	}
	err = rabbit.SendQueue(message)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Error While Send Queue RabbitMQ",
		)
	}
	return responses.MessageResponse(c, http.StatusOK, "Testcases Are Being Run")
}

// @Description Update Group Assigned Chapter Item
// @ID supervisor-update-group-assigned-chapter-item
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.UpdateGroupAssignedChapterItemRequest	true	"Update Group Assigned Chapter Item"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/update_group_assigned_chapter_item [post]
func (supervisorHandler *SupervisorHandler) UpdateGroupAssignedChapterItem(c echo.Context) error {
	updateGroupAssignedChapterItemReq := new(requests.UpdateGroupAssignedChapterItemRequest)
	if err := c.Bind(updateGroupAssignedChapterItemReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := updateGroupAssignedChapterItemReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule,
		updateGroupAssignedChapterItemReq.GroupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(updateGroupAssignedChapterItemReq.GroupId, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			updateGroupAssignedChapterItemReq.GroupId,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var groupChapterSelectedItem []models.GroupChapterSelectedItem
	groupChapterSelectedItemRepo := repositories.NewGroupChapterSelectedItemRepository(
		supervisorHandler.server.DB,
	)
	groupChapterSelectedItemRepo.GetSelectedItemByGroupChapterItemId(
		&groupChapterSelectedItem,
		updateGroupAssignedChapterItemReq.GroupId,
		updateGroupAssignedChapterItemReq.ChapterId,
		updateGroupAssignedChapterItemReq.ItemId,
	)

	existSelectItem := make(map[uuid.UUID]bool)
	for _, existItem := range groupChapterSelectedItem {
		existSelectItem[existItem.ExerciseID] = true
	}

	groupChapterSelectedItemService := groupchapterselecteditem.NewGroupChapterSelectedItemService(
		supervisorHandler.server.DB,
	)

	for _, selectItem := range updateGroupAssignedChapterItemReq.SelectedItem {
		if !existSelectItem[selectItem] {
			err := groupChapterSelectedItemService.Create(
				updateGroupAssignedChapterItemReq.GroupId,
				updateGroupAssignedChapterItemReq.ChapterId,
				int64(updateGroupAssignedChapterItemReq.ItemId),
				selectItem)
			if err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Add Selected Exercise Fail",
				)
			}
		}
	}

	newExistSelectItem := make(map[uuid.UUID]bool)
	for _, existItem := range updateGroupAssignedChapterItemReq.SelectedItem {
		newExistSelectItem[existItem] = true
	}

	for _, selectItem := range groupChapterSelectedItem {
		if !newExistSelectItem[selectItem.ExerciseID] {
			err := groupChapterSelectedItemService.Delete(&selectItem)
			if err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Delete Selected Exercise Fail",
				)
			}
		}
	}

	return responses.MessageResponse(
		c,
		http.StatusOK,
		"Updated All AssignedChapterItem Successfully'",
	)
}

// @Description Update All Group Assigned Chapter Item
// @ID supervisor-update-all-group-assigned-chapter-item
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.UpdateAllGroupAssignedChapterItemRequest	true	"Update All Group Assigned Chapter Item"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/update_all_group_assigned_chapter_item [post]
func (supervisorHandler *SupervisorHandler) UpdateAllGroupAssignedChapterItem(
	c echo.Context,
) error {
	updateAllGroupAssignedChapterItemReq := new(requests.UpdateAllGroupAssignedChapterItemRequest)
	if err := c.Bind(updateAllGroupAssignedChapterItemReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := updateAllGroupAssignedChapterItemReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule,
		updateAllGroupAssignedChapterItemReq.GroupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(updateAllGroupAssignedChapterItemReq.GroupId, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			updateAllGroupAssignedChapterItemReq.GroupId,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	for _, data := range updateAllGroupAssignedChapterItemReq.UpdatePool {
		var groupChapterSelectedItem []models.GroupChapterSelectedItem
		groupChapterSelectedItemRepo := repositories.NewGroupChapterSelectedItemRepository(
			supervisorHandler.server.DB,
		)
		groupChapterSelectedItemRepo.GetSelectedItemByGroupChapterItemId(
			&groupChapterSelectedItem,
			updateAllGroupAssignedChapterItemReq.GroupId,
			updateAllGroupAssignedChapterItemReq.ChapterId,
			data.ItemId,
		)

		existSelectItem := make(map[uuid.UUID]bool)
		for _, existItem := range groupChapterSelectedItem {
			existSelectItem[existItem.ExerciseID] = true
		}

		groupChapterSelectedItemService := groupchapterselecteditem.NewGroupChapterSelectedItemService(
			supervisorHandler.server.DB,
		)

		for _, selectItem := range data.SelectedItem {
			if !existSelectItem[selectItem] {
				err := groupChapterSelectedItemService.Create(
					updateAllGroupAssignedChapterItemReq.GroupId,
					updateAllGroupAssignedChapterItemReq.ChapterId,
					int64(data.ItemId),
					selectItem)
				if err != nil {
					return responses.ErrorResponse(
						c,
						http.StatusInternalServerError,
						"Add Selected Exercise Fail",
					)
				}
			}
		}

		newExistSelectItem := make(map[uuid.UUID]bool)
		for _, existItem := range data.SelectedItem {
			newExistSelectItem[existItem] = true
		}

		for _, selectItem := range groupChapterSelectedItem {
			if !newExistSelectItem[selectItem.ExerciseID] {
				err := groupChapterSelectedItemService.Delete(&selectItem)
				if err != nil {
					return responses.ErrorResponse(
						c,
						http.StatusInternalServerError,
						"Delete Selected Exercise Fail",
					)
				}
			}
		}
	}

	return responses.MessageResponse(
		c,
		http.StatusOK,
		"Updated All AssignedChapterItem Successfully'",
	)
}

// @Description Get Lab Chapter Info
// @ID supervisor-get-lab-chapter-info
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id query string false "group_id"
// @Param chapter_idx query string false "chapter_idx"
// @Success 200		{array}	responses.GetLabChapterInfoResponse
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/get_lab_chapter_info [get]
func (supervisorHandler *SupervisorHandler) GetLabChapterInfo(c echo.Context) error {
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.ExerciseAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	groupId := c.QueryParam("group_id")
	chapterIdx := c.QueryParam("chapter_idx")
	groupUuid, err := uuid.Parse(groupId)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Parse Group_id")
	}
	chapterIdxInt, err := strconv.Atoi(chapterIdx)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Can't Convert Chapter_idx To Int",
		)
	}

	var labClassInfo models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(supervisorHandler.server.DB)
	labClassInfoRepo.GetLabClassInfoByChapterIndex(&labClassInfo, chapterIdxInt)

	var groupChapterSelectItem []models.GroupChapterSelectedItem
	groupChapterSelectItemRepo := repositories.NewGroupChapterSelectedItemRepository(
		supervisorHandler.server.DB,
	)
	groupChapterSelectItemRepo.GetSelectedItemByGroupChapterId(
		&groupChapterSelectItem,
		groupUuid,
		labClassInfo.ChapterID,
	)

	var exerciseList []models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseByChapterID(&exerciseList, labClassInfo.ChapterID)

	response := responses.NewGetLabChapterInfoResponse(
		labClassInfo,
		groupUuid,
		groupChapterSelectItem,
		exerciseList,
	)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Student Group List
// @ID supervisor-get-student-group-list
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id query string false "group_id"
// @Param page query string false "Page"
// @Param pageSize query string false "Page Size"
// @Success 200		{object}	responses.GetStudentWithAssigmentScoreResponse
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/get_student_group_list [get]
func (supervisorHandler *SupervisorHandler) GetStudentGroupList(c echo.Context) error {
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")
	groupId := c.QueryParam("group_id")
	groupUuid, err := uuid.Parse(groupId)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Parse Group_id")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupUuid)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if *existUser.Role != constants.Role.Beyonder &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var labClassInfo []models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(supervisorHandler.server.DB)
	labClassInfoRepo.GetAllLabClassInfos(&labClassInfo)

	var student []models.Student
	studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)
	totalStudents := studentRepo.GetStudentsAndAssignmentScoreByGroupID(
		&student,
		groupUuid,
		page,
		pageSize,
	)

	response := responses.NewGetStudentWithAssigmentScoreByGroupID(
		labClassInfo,
		student,
		groupUuid,
		page,
		pageSize,
		totalStudents,
	)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Set Chapter Permission
// @ID supervisor-set-chapter-permission
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.SetChapterPemissionRequest	true	"Set Chapter Permission"
// @Success 200		{object}	responses.SetChapterPermissionResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/set_chapter_permission [post]
func (supervisorHandler *SupervisorHandler) SetChapterPemission(c echo.Context) error {
	setPermissionReq := new(requests.SetChapterPemissionRequest)
	if err := c.Bind(setPermissionReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := setPermissionReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, setPermissionReq.GroupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(setPermissionReq.GroupId, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(setPermissionReq.GroupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if !utils.ContainsString(constants.ActionTypeList, setPermissionReq.Permission.Type) {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Type")
	}

	var groupChapterPermission models.GroupChapterPermission
	groupChapterPermissionRepo := repositories.NewGroupChapterPermissionRepository(
		supervisorHandler.server.DB,
	)
	groupChapterPermissionRepo.GetGroupChapterPermissionByPK(
		&groupChapterPermission,
		setPermissionReq.GroupId,
		setPermissionReq.ChapterId,
	)

	if setPermissionReq.Permission.Prefix == constants.PermissionPrefix.Submit {
		groupChapterPermission.AllowSubmitType = setPermissionReq.Permission.Type
		if *setPermissionReq.Sync {
			groupChapterPermission.AllowAccessType = setPermissionReq.Permission.Type
		}
		if setPermissionReq.Permission.Type == constants.AccessType.Always ||
			setPermissionReq.Permission.Type == constants.AccessType.Deny {
			groupChapterPermission.SubmitTimeStart = nil
			groupChapterPermission.SubmitTimeEnd = nil
			if *setPermissionReq.Sync {
				groupChapterPermission.AccessTimeStart = nil
				groupChapterPermission.AccessTimeEnd = nil
			}
		} else if setPermissionReq.Permission.Type == constants.AccessType.TimerPaused {
			if setPermissionReq.Permission.TimeStart == nil {
				return responses.ErrorResponse(c, http.StatusBadRequest, "Missing Time Satrt")
			}
			groupChapterPermission.SubmitTimeStart = setPermissionReq.Permission.TimeStart
			if *setPermissionReq.Sync {
				groupChapterPermission.AccessTimeStart = setPermissionReq.Permission.TimeStart
			}
		} else {
			if setPermissionReq.Permission.TimeStart == nil || setPermissionReq.Permission.TimeEnd == nil {
				return responses.ErrorResponse(c, http.StatusBadRequest, "Missing Time Satrt Or Time End")
			}
			groupChapterPermission.SubmitTimeStart = setPermissionReq.Permission.TimeStart
			groupChapterPermission.SubmitTimeEnd = setPermissionReq.Permission.TimeEnd
			if *setPermissionReq.Sync {
				groupChapterPermission.AccessTimeStart = setPermissionReq.Permission.TimeStart
				groupChapterPermission.AccessTimeEnd = setPermissionReq.Permission.TimeEnd
			}
		}
	} else {
		groupChapterPermission.AllowAccessType = setPermissionReq.Permission.Type
		if *setPermissionReq.Sync {
			groupChapterPermission.AllowSubmitType = setPermissionReq.Permission.Type
		}
		if setPermissionReq.Permission.Type == constants.AccessType.Always || setPermissionReq.Permission.Type == constants.AccessType.Deny {
			groupChapterPermission.AccessTimeStart = nil
			groupChapterPermission.AccessTimeEnd = nil
			if *setPermissionReq.Sync {
				groupChapterPermission.SubmitTimeStart = nil
				groupChapterPermission.SubmitTimeEnd = nil
			}
		} else if setPermissionReq.Permission.Type == constants.AccessType.TimerPaused {
			if setPermissionReq.Permission.TimeStart == nil {
				return responses.ErrorResponse(c, http.StatusBadRequest, "Missing Time Satrt")
			}
			groupChapterPermission.AccessTimeStart = setPermissionReq.Permission.TimeStart
			if *setPermissionReq.Sync {
				groupChapterPermission.SubmitTimeStart = setPermissionReq.Permission.TimeStart
			}
		} else {
			if setPermissionReq.Permission.TimeStart == nil || setPermissionReq.Permission.TimeEnd == nil {
				return responses.ErrorResponse(c, http.StatusBadRequest, "Missing Time Satrt Or Time End")
			}
			groupChapterPermission.AccessTimeStart = setPermissionReq.Permission.TimeStart
			groupChapterPermission.AccessTimeEnd = setPermissionReq.Permission.TimeEnd
			if *setPermissionReq.Sync {
				groupChapterPermission.SubmitTimeStart = setPermissionReq.Permission.TimeStart
				groupChapterPermission.SubmitTimeEnd = setPermissionReq.Permission.TimeEnd
			}
		}
	}

	groupChapterPermissionService := groupchapterpermission.NewGroupChapterPermissionService(
		supervisorHandler.server.DB,
	)
	if groupChapterPermission.AccessTimeStart != nil {
		bufferTimeStart := groupChapterPermission.AccessTimeStart.Add(-5 * time.Second)
		groupChapterPermission.AccessTimeStart = &bufferTimeStart
	}
	if groupChapterPermission.SubmitTimeStart != nil {
		bufferTimeStart := groupChapterPermission.SubmitTimeStart.Add(-5 * time.Second)
		groupChapterPermission.SubmitTimeStart = &bufferTimeStart
	}
	err = groupChapterPermissionService.UpdateByModel(&groupChapterPermission)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	redis := redis_client.NewRedisAction(supervisorHandler.server.Redis)
	redisCnl := fmt.Sprintf(
		"%s:%s",
		constants.RedisChannel.GroupPermission,
		setPermissionReq.GroupId,
	)
	if err := redis.PublishMessage(redisCnl, "permission-change"); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
	}

	response := responses.NewSetChapterPermissionResponse(groupChapterPermission)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Set Allow Group Login
// @ID supervisor-set-allow-group-login
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.SetAllowGroupLoginRequest	true	"Set Allow Group Login"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/set_allow_group_login [post]
func (supervisorHandler *SupervisorHandler) SetAllowGroupLogin(c echo.Context) error {
	setAllowGroupLoginReq := new(requests.SetAllowGroupLoginRequest)
	if err := c.Bind(setAllowGroupLoginReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := setAllowGroupLoginReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassSchedulePreloadByGroupID(
		&classSchedule,
		setAllowGroupLoginReq.GroupID,
	)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(setAllowGroupLoginReq.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			setAllowGroupLoginReq.GroupID,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	classScheduleService := classschedule.NewClassScheduleService(supervisorHandler.server.DB)
	classScheduleService.UpdateAllowLogin(&classSchedule, setAllowGroupLoginReq.AllowLogin)

	return responses.MessageResponse(c, http.StatusOK, "Setting Updated Successfully")
}

// @Description Set Allow Group Upload Picture
// @ID supervisor-set-allow-group-upload-picture
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.SetAllowGroupUploadPictureRequest	true	"Set Allow Group Upload Picture"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/set_allow_group_upload_picture [post]
func (supervisorHandler *SupervisorHandler) SetAllowGroupUploadPicture(c echo.Context) error {
	setAllowGroupUploadPictureRequest := new(requests.SetAllowGroupUploadPictureRequest)
	if err := c.Bind(setAllowGroupUploadPictureRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := setAllowGroupUploadPictureRequest.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassSchedulePreloadByGroupID(
		&classSchedule,
		setAllowGroupUploadPictureRequest.GroupID,
	)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(setAllowGroupUploadPictureRequest.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			setAllowGroupUploadPictureRequest.GroupID,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	classScheduleService := classschedule.NewClassScheduleService(supervisorHandler.server.DB)
	classScheduleService.UpdateAllowUploadPicture(
		&classSchedule,
		setAllowGroupUploadPictureRequest.AllowUploadPicture,
	)

	return responses.MessageResponse(c, http.StatusOK, "Setting Updated Successfully")
}

// @Description Delete Exercise
// @ID supervisor-delete-exercise
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param exercise_id path string true "Exercise ID"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/exercise/{exercise_id} [delete]
func (supervisorHandler *SupervisorHandler) DeleteExercise(c echo.Context) error {
	exerciseIdStr := c.Param("exercise_id")
	exerciseId, err := uuid.Parse(exerciseIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.ExerciseAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var labExercise models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseByID(exerciseId, &labExercise)

	var exerciseSubmission []models.ExerciseSubmission
	exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(
		supervisorHandler.server.DB,
	)
	exerciseSubmissionRepo.GetSubmissionByExerciseID(exerciseId, &exerciseSubmission)

	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	for _, submission := range exerciseSubmission {
		minioAction.DeleteFileInMinio(
			supervisorHandler.server.Config.Minio.BucketStudentCode,
			submission.SourcecodeFilename,
		)
	}
	minioAction.DeleteFileInMinio(
		supervisorHandler.server.Config.Minio.BucketSupervisorCode,
		*labExercise.Sourcecode,
	)

	labExerciseService := labexercise.NewLabExerciseService(supervisorHandler.server.DB)
	labExerciseService.Delete(&labExercise)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Exercise.")
	}

	return responses.MessageResponse(c, http.StatusOK, "Delete Exercise Done")
}

// @Description Update Student Can Submit
// @ID supervisor-update-student-can-submit
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param student_id path string true "Student ID"
// @Param params body	requests.CanSubmitRequest	true	"Can Submit Request"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/student_can_submit/{student_id} [put]
func (supervisorHandler *SupervisorHandler) UpdateStudentCanSubmit(c echo.Context) error {
	studentIdStr := c.Param("student_id")
	studentId, err := uuid.Parse(studentIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	canSubmitReq := new(requests.CanSubmitRequest)

	if err := c.Bind(canSubmitReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}

	if err := canSubmitReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	existStudent := models.Student{}
	studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)
	studentRepo.GetStudentByStuID(&existStudent, studentId)
	if existStudent.StuID != studentId {
		return responses.ErrorResponse(c, http.StatusNotFound, "User Not found.")
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, *existStudent.GroupID)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(*existStudent.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(*existStudent.GroupID, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	studentService := student.NewStudentService(supervisorHandler.server.DB)
	studentService.UpdateCanSubmit(&existStudent, canSubmitReq.CanSubmit)

	redis := redis_client.NewRedisAction(supervisorHandler.server.Redis)
	redisCnl := fmt.Sprintf(
		"%s:%s",
		constants.RedisChannel.UserEvent,
		existStudent.StuID,
	)
	redisMsg := redis.NewMessage("can-submit", &existStudent.StuID)
	if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
	}

	return responses.MessageResponse(c, http.StatusOK, "Updated Student Can Submit")
}

// @Description Get Student Info
// @ID supervisor-get-student-info
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param student_id path string true "Student ID"
// @Success 200		{object}	responses.UserStudentInfoResponse
// @Failure 404		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/student_info/{student_id} [get]
func (supervisorHandler *SupervisorHandler) GetStudentInfo(c echo.Context) error {
	studentIdStr := c.Param("student_id")
	studentId, err := uuid.Parse(studentIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	existStudent := models.User{}
	userRepository.GetUserStudentAndGroupByUserID(&existStudent, studentId)
	if existStudent.UserID != studentId {
		return responses.ErrorResponse(c, http.StatusNotFound, "User Not found.")
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, *existStudent.Student.GroupID)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(*existStudent.Student.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			*existStudent.Student.GroupID,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	response := responses.NewUserStudentInfoResponse(existStudent)

	if existStudent.UserID != studentId {
		return responses.ErrorResponse(c, http.StatusNotFound, "Not found student.")
	}

	return responses.Response(c, http.StatusOK, response)
}

// @Description Logout All Student In Group
// @ID supervisor-logout-all-student-in-group
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Success 200		{object}	responses.LogoutAllStudentResponse
// @Failure 400		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/logout_all_student/{group_id} [put]
func (supervisorHandler *SupervisorHandler) LogoutAllStudentInGroup(c echo.Context) error {
	groupIdStr := c.Param("group_id")
	groupId, err := uuid.Parse(groupIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassSchedulePreloadByGroupID(&classSchedule, groupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.GroupAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var studentList []models.Student
	studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)
	studentRepo.GetStudentInGroupID(&studentList, groupId)

	userService := user.NewUserService(supervisorHandler.server.DB)

	var mu sync.Mutex
	count := 0
	stuLogout := make([]models.Student, 0)
	wg := sync.WaitGroup{}
	for _, stu := range studentList {
		if stu.User.IsOnline {
			wg.Add(1)
			go func(student models.Student) {
				defer wg.Done()

				userService.UpdateIsOnline(stu.User, false)

				mu.Lock()
				count++
				stuLogout = append(stuLogout, student)
				mu.Unlock()
			}(stu)
		}
	}
	wg.Wait()

	ip, port, userAgent := utils.GetNetworkRequest(c)

	activityLogService := activitylog.NewActivityLogService(supervisorHandler.server.DB)
	newLog, err := activityLogService.Create(
		&groupId,
		existUser.Username,
		ip,
		&port,
		&userAgent,
		constants.LogPage.ManageStudent,
		constants.LogAction.LogoutAllStudent,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Insert Log.")
	}

	redis := redis_client.NewRedisAction(supervisorHandler.server.Redis)
	redisCnl := fmt.Sprintf(
		"%s:%s",
		constants.RedisChannel.OnlineStudent,
		groupId,
	)
	var nilUUID uuid.UUID
	redisMsg := redis.NewMessage("logout-all", &nilUUID)
	if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
	}

	redisCnl = fmt.Sprintf(
		"%s:%s",
		constants.RedisChannel.Log,
		groupId,
	)
	if err := redis.PublishMessage(redisCnl, newLog); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
	}

	redisCnl = fmt.Sprintf(
		"%s:%s",
		constants.RedisChannel.GroupPermission,
		groupId,
	)
	if err := redis.PublishMessage(redisCnl, "logout-all"); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
	}

	response := responses.NewLogoutAllStudentResponse(count, stuLogout)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Delete Student
// @ID supervisor-delete-student
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param stu_id path string true "Student ID"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/student/{stu_id} [delete]
func (supervisorHandler *SupervisorHandler) DeleteStudent(c echo.Context) error {
	studentIdStr := c.Param("stu_id")
	studentId, err := uuid.Parse(studentIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var userInfo models.User
	userRepository.GetUserByUserID(&userInfo, studentId)
	if *userInfo.Role != constants.Role.Student {
		return responses.ErrorResponse(
			c,
			http.StatusForbidden,
			"This User Isn't Student Can't delete.",
		)
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, *userInfo.Student.GroupID)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(*userInfo.Student.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(*userInfo.Student.GroupID, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var exerciseSubmission []models.ExerciseSubmission
	exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(
		supervisorHandler.server.DB,
	)
	exerciseSubmissionRepo.GetSubmissionByStudentID(studentId, &exerciseSubmission)

	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	for _, submission := range exerciseSubmission {
		minioAction.DeleteFileInMinio(
			supervisorHandler.server.Config.Minio.BucketStudentCode,
			submission.SourcecodeFilename,
		)
	}

	userService := user.NewUserService(supervisorHandler.server.DB)
	err = userService.Delete(&userInfo)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Delete Student.")
	}

	return responses.MessageResponse(c, http.StatusOK, "Delete Student Done.")
}

// @Description Get Exercise Data
// @ID supervisor-get-exercise-data
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param exercise_id path string true "Exercise ID"
// @Success 200		{object}	responses.GetExerciseDataResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/get_exercise_data/{exercise_id} [get]
func (supervisorHandler *SupervisorHandler) GetExerciseData(c echo.Context) error {
	exerciseIdStr := c.Param("exercise_id")
	exerciseId, err := uuid.Parse(exerciseIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.ExerciseAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var labExerciseData models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseByID(exerciseId, &labExerciseData)

	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	minioFile, err := minioAction.GetFromMinio(
		supervisorHandler.server.Config.Minio.BucketSupervisorCode,
		*labExerciseData.Sourcecode,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Get Data From Minio Object Fail",
		)
	}
	defer minioFile.Close()

	var sourcecode strings.Builder
	_, err = io.Copy(&sourcecode, minioFile)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Copy Data From Minio Object Fail",
		)
	}

	response := responses.NewGetExerciseDataResponse(labExerciseData, sourcecode.String())
	return responses.Response(c, http.StatusOK, response)
}

// @Description Cancle Student Submission
// @ID supervisor-cancle-student-submission
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param submission_id path string true "Submission ID"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/cancle_student_submission/{submission_id} [put]
func (supervisorHandler *SupervisorHandler) CancleStduentSubmission(c echo.Context) error {
	submissionIdStr := c.Param("submission_id")
	submissionId, err := uuid.Parse(submissionIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	var exerciseSubmissionData models.ExerciseSubmission
	exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(
		supervisorHandler.server.DB,
	)
	err = exerciseSubmissionRepo.GetSubmissionByID(submissionId, &exerciseSubmissionData)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Not Found Submission.")
		}
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(
		&classSchedule,
		*exerciseSubmissionData.Student.GroupID,
	)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(*exerciseSubmissionData.Student.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			*exerciseSubmissionData.Student.GroupID,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	exerciseSubmissionService := exercisesubmission.NewExerciseSubmissionService(
		supervisorHandler.server.DB,
	)
	exerciseSubmissionService.CancleSubmission(&exerciseSubmissionData)

	studentAssignChapterItemService := studentassignmentchapteritem.NewStudentAssignmentChapterItem(
		supervisorHandler.server.DB,
	)
	studentAssignChapterItemService.ResetMarking(
		exerciseSubmissionData.StuID,
		*exerciseSubmissionData.LabExercise.ChapterID,
		*exerciseSubmissionData.LabExercise.Level,
	)

	ip, port, userAgent := utils.GetNetworkRequest(c)

	activityLogService := activitylog.NewActivityLogService(supervisorHandler.server.DB)
	newLog, err := activityLogService.Create(
		exerciseSubmissionData.Student.GroupID,
		existUser.Username,
		ip,
		&port,
		&userAgent,
		constants.LogPage.ManageStudent,
		fmt.Sprintf(
			"reject submission #%s stu_id:%s chapter:%s item:%s",
			submissionId,
			exerciseSubmissionData.StuID,
			exerciseSubmissionData.LabExercise.ChapterID,
			*exerciseSubmissionData.LabExercise.Level,
		),
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Insert Log.")
	}

	redis := redis_client.NewRedisAction(supervisorHandler.server.Redis)
	redisCnl := fmt.Sprintf(
		"%s:%s",
		constants.RedisChannel.Log,
		exerciseSubmissionData.Student.GroupID,
	)

	if err := redis.PublishMessage(redisCnl, newLog); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
	}

	redisCnl = fmt.Sprintf(
		"%s:%s",
		constants.RedisChannel.UserEvent,
		exerciseSubmissionData.StuID,
	)
	redisMsg := redis.NewMessage("reject-submission", &exerciseSubmissionData.StuID)
	if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
	}

	return responses.MessageResponse(c, http.StatusOK, "Submission Canceled Successfully.")
}

// @Description Get Student Chapter List
// @ID supervisor-get-student-chapter-list
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param studentId query string true "Student ID"
// @Success 200		{array}	responses.GetChapterListResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/student_chapter_list [get]
func (supervisorHandler *SupervisorHandler) GetStudentChapterList(c echo.Context) error {
	studentIdStr := c.QueryParam("studentId")
	studentId, err := uuid.Parse(studentIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Parse Group_id")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	existStudent := models.User{}
	userRepository.GetUserByUserID(&existStudent, studentId)
	if existStudent.UserID != studentId {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, *existStudent.Student.GroupID)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(*existStudent.Student.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			*existStudent.Student.GroupID,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var labClassInfos []models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(supervisorHandler.server.DB)
	labClassInfoRepo.GetAllLabClassInfos(&labClassInfos)

	var groupChapterPermission []models.GroupChapterPermission
	groupChapterPermissionRepo := repositories.NewGroupChapterPermissionRepository(
		supervisorHandler.server.DB,
	)
	groupChapterPermissionRepo.GetGroupChapterPermissionByGroupID(
		&groupChapterPermission,
		*existStudent.Student.GroupID,
	)

	var allGroupChapterItems []models.GroupAssignmentChapterItem
	groupChapterItemRepo := repositories.NewGroupAssignmentChapterItemRepository(
		supervisorHandler.server.DB,
	)
	groupChapterItemRepo.GetAllGroupAssignmentChapterItemsByGroupId(
		&allGroupChapterItems,
		*existStudent.Student.GroupID,
	)

	var allStudentAssignChapterItems []models.StudentAssignmentChapterItem
	studentAssignItemRepo := repositories.NewStudentAssignChapterItemRepository(
		supervisorHandler.server.DB,
	)
	studentAssignItemRepo.GetAllStudentAssignChapterWithSubmission(
		&allStudentAssignChapterItems,
		studentId,
	)

	response := responses.NewGetChapterListResponse(
		groupChapterPermission,
		allGroupChapterItems,
		allStudentAssignChapterItems,
		true,
	)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Assgined Student Exercise
// @ID supervisor-get-assigned-student-exercise
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param stu_id query string true "Student_ID"
// @Param chapter_idx query string true "Chapter Index"
// @Param item_id query string false "Item ID"
// @Success 200		{object}	responses.GetAssginStudentExerciseResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/assigned_student_exercise [get]
func (supervisorHandler *SupervisorHandler) GetAssginStudentExercise(c echo.Context) error {
	stuIdStr := c.QueryParam("stu_id")
	chapterIdxStr := c.QueryParam("chapter_idx")
	itemIdStr := c.QueryParam("item_id")

	if stuIdStr == "" || chapterIdxStr == "" || itemIdStr == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	stuId, err := uuid.Parse(stuIdStr)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Can't Convert Student ID",
		)
	}

	chapterInt, err := strconv.Atoi(chapterIdxStr)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Can't Convert Chapter Index",
		)
	}

	itemInt, err := strconv.Atoi(itemIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Convert Item ID")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	existStudent := models.User{}
	userRepository.GetUserByUserID(&existStudent, stuId)
	if existStudent.UserID != stuId {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, *existStudent.Student.GroupID)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(*existStudent.Student.GroupID, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.StudentAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(
			*existStudent.Student.GroupID,
			existUser.UserID,
		) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var labClassInfo models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(supervisorHandler.server.DB)
	labClassInfoRepo.GetLabClassInfoByChapterIndex(&labClassInfo, chapterInt)

	if int64(chapterInt) > labClassInfoRepo.GetCount() || chapterInt < 0 {
		return responses.ErrorResponse(c, http.StatusForbidden, "Chapter Index Out of Range.")
	}

	if itemInt > labClassInfo.NoItems || itemInt < 0 {
		return responses.ErrorResponse(c, http.StatusForbidden, "Item ID Out of Range.")
	}

	var studentAssignChapterItems models.StudentAssignmentChapterItem
	studentAssignItemRepo := repositories.NewStudentAssignChapterItemRepository(
		supervisorHandler.server.DB,
	)
	err = studentAssignItemRepo.GetStudentAssignChapterItem(
		&studentAssignChapterItems,
		stuId,
		labClassInfo.ChapterID,
		itemInt,
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			var studentData models.Student
			studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)
			studentRepo.GetStudentByStuID(&studentData, stuId)

			var groupChapterPermission models.GroupChapterPermission
			groupChapterPermissionRepo := repositories.NewGroupChapterPermissionRepository(
				supervisorHandler.server.DB,
			)
			groupChapterPermissionRepo.GetGroupChapterPermissionByPK(
				&groupChapterPermission,
				*studentData.GroupID,
				labClassInfo.ChapterID,
			)
			studentAssignChapterItemService := studentassignmentchapteritem.NewStudentAssignmentChapterItem(
				supervisorHandler.server.DB,
			)
			createData, err := studentAssignChapterItemService.Create(
				stuId,
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
		return responses.MessageResponse(c, http.StatusOK, "No Exercise Assigned To This Student.")
	}

	var labExercise models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	labExerciseRepo.GetLabExerciseByID(*studentAssignChapterItems.ExerciseID, &labExercise)

	response := responses.NewGetAssginStudentExerciseResponse(labExercise)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Update Exercise
// @ID supervisor-update-exercise
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.UpdateLabExerciseRequest	true	"Update Exercise"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/exercise [put]
func (supervisorHandler *SupervisorHandler) UpdateExercise(c echo.Context) error {
	updateLabExerciseReq := new(requests.UpdateLabExerciseRequest)
	if err := c.Bind(updateLabExerciseReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := updateLabExerciseReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.ExerciseAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var labExerciseData models.LabExercise
	labExerciseRepo := repositories.NewLabExerciseRepository(supervisorHandler.server.DB)
	err = labExerciseRepo.GetLabExerciseByID(*updateLabExerciseReq.ExerciseID, &labExerciseData)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.ErrorResponse(
				c,
				http.StatusBadRequest,
				"Not Found Exercise.",
			)
		}
	}

	labExerciseService := labexercise.NewLabExerciseService(supervisorHandler.server.DB)
	err = labExerciseService.UpdateLabExercise(&labExerciseData, *updateLabExerciseReq)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	tempFile, err := utils.CreateTempFile(
		fmt.Sprintf(*labExerciseData.Sourcecode+"*.py"),
		updateLabExerciseReq.Sourcecode,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Create Temp File Fail %s", err),
		)
	}
	defer os.Remove(tempFile.Name())

	minioAction := minioclient.NewMinioAction(supervisorHandler.server.Minio)
	err = minioAction.DeleteFileInMinio(
		supervisorHandler.server.Config.Minio.BucketSupervisorCode,
		*labExerciseData.Sourcecode,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	filename, err := minioAction.UploadToMinio(
		tempFile,
		supervisorHandler.server.Config.Minio.BucketSupervisorCode,
		false,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	labExerciseService.UpdateLabExerciseSourcecode(&labExerciseData, filename)

	exerciseTestcaseService := exercisetestcase.NewExerciseTestcaseService(
		supervisorHandler.server.DB,
	)
	exerciseTestcaseService.UpdateTestcaseIsReadyByExerciseID(
		*updateLabExerciseReq.ExerciseID,
		constants.TestcaseIsReadyType.No,
	)

	var exerciseTestcase []models.ExerciseTestcase
	exerciseTestcaseRepo := repositories.NewExerciseTestcaseRepository(supervisorHandler.server.DB)
	exerciseTestcaseRepo.GetTestcasesByExerciseID(labExerciseData.ExerciseID, &exerciseTestcase)

	testcaseList := make([]requests.ExerciseTestcaseReq, 0)
	for _, testcase := range exerciseTestcase {
		testcaseList = append(testcaseList, requests.ExerciseTestcaseReq{
			TestcaseID:      testcase.TestcaseID,
			ExerciseID:      testcase.ExerciseID,
			IsReady:         testcase.IsReady,
			TestcaseContent: testcase.TestcaseContent,
			IsActive:        *testcase.IsActive,
			IsShowStudent:   *testcase.IsShowStudent,
			TestcaseNote:    *testcase.TestcaseNote,
			TestcaseOutput:  *testcase.TestcaseOutput,
			TestcaseError:   *testcase.TestcaseError,
		})
	}

	rabbit := rabbitmq_client.NewRabbitMQAction(
		supervisorHandler.server.RabitMQ,
		supervisorHandler.server.Config,
	)
	message := requests.AddTestcaseRabitMessage{
		JobId:        *updateLabExerciseReq.JobID,
		JobType:      "upsert-testcase",
		ExerciseId:   *updateLabExerciseReq.ExerciseID,
		TestcaseList: testcaseList,
		Sourcecode:   updateLabExerciseReq.Sourcecode,
	}
	err = rabbit.SendQueue(message)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Error While Send Queue RabbitMQ",
		)
	}

	return responses.MessageResponse(
		c,
		http.StatusOK,
		"Update Exercise Successfully Wait For Testcase Output",
	)
}

// @Description Create Admin
// @ID supervisor-create-admin
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.CreateAdminRequest	true	"Create Admin"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/admin [post]
func (supervisorHandler *SupervisorHandler) CreateAdmin(c echo.Context) error {
	createAdminReq := new(requests.CreateAdminRequest)
	if err := c.Bind(createAdminReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := createAdminReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
	}

	if !utils.ContainsString(constants.AdminRoleList, createAdminReq.Role) {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Role.")
	}

	userService := user.NewUserService(supervisorHandler.server.DB)
	userId, err := userService.Create(
		createAdminReq.Username,
		createAdminReq.Username,
		createAdminReq.Firstname,
		createAdminReq.Lastname,
		createAdminReq.Gender,
		createAdminReq.Role,
		createAdminReq.DeptID,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Create User.")
	}

	switch createAdminReq.Role {
	case constants.Role.Supervisor:
		{
			supervisorService := supervisor.NewSupervisorService(supervisorHandler.server.DB)
			err = supervisorService.Create(userId, "คอมพิวเตอร์")
			if err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Can't Create Supervisor.",
				)
			}
			break
		}
	case constants.Role.Ta:
		{
			taService := ta.NewTaService(supervisorHandler.server.DB)
			err = taService.CreateTa(userId, nil, nil)
			if err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Can't Create TA.",
				)
			}
			break
		}
	case constants.Role.Executive:
		{
			executiveService := executive.NewExecutiveService(supervisorHandler.server.DB)
			err = executiveService.Create(userId)
			if err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Can't Create Executive.",
				)
			}
			break
		}
	default:
		{
			return responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid Role.")
		}
	}

	return responses.MessageResponse(c, http.StatusOK, "Create Admin Success.")
}

// @Description Delete Admin
// @ID supervisor-delete-admin
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param admin_id path string true "Admin ID"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/admin/{admin_id}  [delete]
func (supervisorHandler *SupervisorHandler) DeleteAdmin(c echo.Context) error {
	adminIdStr := c.Param("admin_id")
	adminId, err := uuid.Parse(adminIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var deleteAdmin models.User
	userRepository.GetUserByUserID(&deleteAdmin, adminId)
	if !utils.ContainsString(constants.AdminRoleList, *deleteAdmin.Role) {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Admin Role.")
	}

	if *deleteAdmin.Role == constants.Role.Supervisor &&
		*existUser.Role != constants.Role.Beyonder {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	userService := user.NewUserService(supervisorHandler.server.DB)
	err = userService.DeleteAdmin(&deleteAdmin)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Delete Admin Fail.")
	}

	return responses.MessageResponse(c, http.StatusOK, "Delete Admin Success.")
}

// @Description Create Department
// @ID supervisor-create-department
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.CreateDepartmentRequest	true	"Create Department"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/department [post]
func (supervisorHandler *SupervisorHandler) CreateDepartment(c echo.Context) error {
	createDepartmentReq := new(requests.CreateDepartmentRequest)
	if err := c.Bind(createDepartmentReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := createDepartmentReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
	}

	departmentService := department.NewDepartmetService(supervisorHandler.server.DB)
	err = departmentService.Create(createDepartmentReq.Name, createDepartmentReq.Name_EN)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Fail To Create Department.",
		)
	}

	return responses.MessageResponse(c, http.StatusOK, "Create Department Success.")
}

// @Description Set Role Permission
// @ID supervisor-set-role-permission
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param params body	requests.SetRolePermissionRequest	true	"Set Role Permission"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/set_role_permission [post]
func (supervisorHandler *SupervisorHandler) SetRolePermission(c echo.Context) error {
	setRolePermissionReq := new(requests.SetRolePermissionRequest)
	if err := c.Bind(setRolePermissionReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}
	if err := setRolePermissionReq.Validate(); err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid Request",
		)
	}

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
	rolePermissionService := rolepermission.NewRolePermissionService(supervisorHandler.server.DB)

	var wg sync.WaitGroup
	errChan := make(chan error, len(setRolePermissionReq.Data))
	for _, permissionData := range setRolePermissionReq.Data {
		wg.Add(1)
		go func(role string, permission []string) {
			defer wg.Done()
			var rolePermission []models.RolePermission
			rolePermissionRepo.GetPermissionByRole(&rolePermission, role)

			exitsPermission := make(map[string]bool)
			for _, permission := range rolePermission {
				exitsPermission[permission.Permission] = true
			}
			for _, newPermission := range permission {
				if !exitsPermission[newPermission] {
					err := rolePermissionService.Create(role, newPermission)
					if err != nil {
						errChan <- fmt.Errorf("error : %v", err.Error())
						return
					}
				}
			}

			newExitsPermission := make(map[string]bool)
			for _, exitsPermission := range permission {
				newExitsPermission[exitsPermission] = true
			}
			for _, oldPermission := range rolePermission {
				if !newExitsPermission[oldPermission.Permission] {
					err := rolePermissionService.Delete(&oldPermission)
					if err != nil {
						errChan <- fmt.Errorf("error : %v", err.Error())
						return
					}
				}
			}
		}(permissionData.Role, permissionData.Permission)
	}
	wg.Wait()
	close(errChan)
	errList := make([]string, 0)
	for err := range errChan {
		if err != nil {
			errList = append(errList, err.Error())
		}
	}
	if len(errList) > 0 {
		errString := strings.Join(errList, "\n")
		return responses.ErrorResponse(c, http.StatusBadRequest, errString)
	}

	return responses.MessageResponse(c, http.StatusOK, "Set Role Permission Successful.")
}

// @Description Get Role Permission
// @ID supervisor-get-role-permission
// @Tags Supervisor
// @Accept json
// @Produce json
// @Success 200		{object}	responses.GetRolePermissionResponse
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/role_permission [get]
func (supervisorHandler *SupervisorHandler) GetRolePermission(c echo.Context) error {
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateAdminRole(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var rolePermission []models.RolePermission
	rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
	rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

	response := responses.NewGetRolePermissionResponse(rolePermission, existUser)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get All Role Permission
// @ID supervisor-get-all-role-permission
// @Tags Supervisor
// @Accept json
// @Produce json
// @Success 200		{object}	responses.GetAllRolePermissionResponse
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/all_role_permission [get]
func (supervisorHandler *SupervisorHandler) GetAllRolePermission(c echo.Context) error {
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission")
	}

	var rolePermission []models.RolePermission
	rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
	rolePermissionRepo.GetAllPermissionRole(&rolePermission)

	response := responses.NewGetAllRolePermissionResponse(rolePermission)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Average Chapter Score
// @ID supervisor-get-average-Chapter-score
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id query string false "Group_ID"
// @Param year query string false "Year"
// @Success 200		{array}		float64
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/stats/score/chapter  [get]
func (supervisorHandler *SupervisorHandler) GetAverageChapterScore(c echo.Context) error {
	groupIdStr := c.QueryParam("groupId")
	year := c.QueryParam("year")
	groupId, _ := uuid.Parse(groupIdStr)

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)
		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if groupId != uuid.Nil {
		var classSchedule models.ClassSchedule
		classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
		classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupId)

		if classSchedule.GroupID == uuid.Nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
		}

		if *existUser.Role == constants.Role.Supervisor &&
			*classSchedule.SupervisorID != existUser.UserID &&
			!classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		} else if *existUser.Role != constants.Role.Beyonder && !classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var allLabClassInfo []models.LabClassInfo
	labClassInfoRepo := repositories.NewLabClassInfoRepository(supervisorHandler.server.DB)
	labClassInfoRepo.GetAllLabClassInfos(&allLabClassInfo)

	studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)
	studentCount := studentRepo.GetStudentGroupOrYearCount(groupId, year)
	studentAssignItemsRepo := repositories.NewStudentAssignChapterItemRepository(
		supervisorHandler.server.DB,
	)

	averrageScore := make([]float64, len(allLabClassInfo))

	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(allLabClassInfo))

	for _, labClassInfo := range allLabClassInfo {
		wg.Add(1)
		go func(labClassInfo models.LabClassInfo) {
			defer wg.Done()

			var studentAssignItems []models.StudentAssignmentChapterItem
			err := studentAssignItemsRepo.GetStudentChapterByChapterIDAndGroupOrYear(
				&studentAssignItems,
				groupId,
				labClassInfo.ChapterID,
				year,
			)
			if err != nil {
				errChan <- fmt.Errorf("error : %v", err.Error())
				return
			}

			totalScore := 0
			for _, assignItem := range studentAssignItems {
				totalScore = totalScore + assignItem.Marking
			}

			average := float64(totalScore)
			if studentCount != 0 {
				average = float64(totalScore) / float64(studentCount)
			}

			mu.Lock()
			averrageScore[labClassInfo.ChapterIndex-1] = average
			mu.Unlock()
		}(labClassInfo)
	}
	wg.Wait()
	close(errChan)
	errList := make([]string, 0)
	for err := range errChan {
		if err != nil {
			errList = append(errList, err.Error())
		}
	}
	if len(errList) > 0 {
		errString := strings.Join(errList, "\n")
		return responses.ErrorResponse(c, http.StatusInternalServerError, errString)
	}

	return responses.Response(c, http.StatusOK, averrageScore)
}

// @Description Get Total Student (Have to send only group_id or year)
// @ID supervisor-get-total-student
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id query string false "Group_ID"
// @Param status query string false "Status"
// @Param year query string false "Year"
// @Success 200		{object}	responses.TotalStudentResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/students/total [get]
func (supervisorHandler *SupervisorHandler) GetTotalStudent(c echo.Context) error {
	groupId := c.QueryParam("groupId")
	status := c.QueryParam("status")
	year := c.QueryParam("year")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var totalStudent int64
	studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)

	if groupId != "" {
		groupUuid, err := uuid.Parse(groupId)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid Group ID.")
		}

		var classSchedule models.ClassSchedule
		classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
		classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupUuid)

		if classSchedule.GroupID == uuid.Nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
		}

		classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
		if !utils.ValidateSupervisorAndBeyonder(existUser) {
			if !(classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID)) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}

		if *existUser.Role == constants.Role.Supervisor &&
			*classSchedule.SupervisorID != existUser.UserID {
			if !classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}
	}
	totalStudent = studentRepo.GetTotalStudent(groupId, year, status)

	response := responses.NewTotalStudentResponse(totalStudent)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Total Staff
// @ID supervisor-get-total-staff
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id query string false "Group_ID"
// @Success 200		{object}	responses.TotalStaffResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/staffs/total [get]
func (supervisorHandler *SupervisorHandler) GetTotalStaff(c echo.Context) error {
	groupId := c.QueryParam("groupId")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if groupId != "" {
		groupUuid, err := uuid.Parse(groupId)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid Group ID.")
		}

		var classSchedule models.ClassSchedule
		classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
		classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupUuid)

		if classSchedule.GroupID == uuid.Nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
		}

		classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
		if !utils.ValidateSupervisorAndBeyonder(existUser) {
			if !(classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID)) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}

		if *existUser.Role == constants.Role.Supervisor &&
			*classSchedule.SupervisorID != existUser.UserID {
			if !classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}
	} 

	userRepo := repositories.NewUserRepository(supervisorHandler.server.DB)
	totalStaff := userRepo.GetTotalAdmin(groupId)
	
	response := responses.NewTotalStaffResponse(totalStaff)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Total Submissions
// @ID supervisor-get-total-submissions
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id query string false "Group_ID"
// @Param year query string false "Year"
// @Success 200		{object}	responses.TotalSubmissionsResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/submissions/total [get]
func (supervisorHandler *SupervisorHandler) GetTotalExerciseSubmissions(c echo.Context) error {
	groupId := c.QueryParam("group_id")
	year := c.QueryParam("year")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var totalSubmissions int64
	exerciseSubmissionRepo := repositories.NewExerciseSubmissionRepository(
		supervisorHandler.server.DB,
	)

	if groupId != "" {

		groupUuid, err := uuid.Parse(groupId)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid Group ID.")
		}

		var classSchedule models.ClassSchedule
		classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
		classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupUuid)

		if classSchedule.GroupID == uuid.Nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
		}

		classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
		if !utils.ValidateSupervisorAndBeyonder(existUser) {
			if !(classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID)) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}

		if *existUser.Role == constants.Role.Supervisor &&
			*classSchedule.SupervisorID != existUser.UserID {
			if !classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}
	}

	totalSubmissions = exerciseSubmissionRepo.GetTotalSubmissions(groupId, year)

	response := responses.NewTotalSubmissionsResponse(totalSubmissions)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Total Groups
// @ID supervisor-get-total-groups
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param year query string false "Year"
// @Success 200		{object}	responses.TotalGroupsResponse
// @Failure 403		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/groups/total [get]
func (supervisorHandler *SupervisorHandler) GetTotalGroup(c echo.Context) error {
	year := c.QueryParam("year")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var totalGroup int64
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	totalGroup = classScheduleRepo.GetTotalGroup(year)

	response := responses.NewTotalGroupsResponse(totalGroup)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Total Groups
// @ID supervisor-get-score-ranking
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id path string true "Group ID"
// @Success 200		{array}		responses.StudentRankingResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/score_ranking/{group_id} [get]
func (supervisorHandler *SupervisorHandler) GetScoreRankingByGroup(c echo.Context) error {
	groupIdStr := c.Param("group_id")
	groupId, err := uuid.Parse(groupIdStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request Param")
	}

	classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !(classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) && utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin)) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var classSchedule models.ClassSchedule
	classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
	classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupId)

	if classSchedule.GroupID == uuid.Nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
	}

	if *existUser.Role == constants.Role.Supervisor &&
		*classSchedule.SupervisorID != existUser.UserID {
		if !classLabStaffRepo.CheckStaffValidInClass(groupId, existUser.UserID) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var students []models.StudentWithAggregate
	studentRepo := repositories.NewStudentRepository(supervisorHandler.server.DB)
	err = studentRepo.GetStudentGroupRanking(&students, groupId)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Fail When Get Ranking.")
	}

	response := responses.NewStudentRankingResponse(students)

	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Submissions Over Time
// @ID supervisor-get-submissions-over-time
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param group_id query string false "Group_ID"
// @Param year query string false "Year"
// @Success 200		{object}	responses.StatsSubmissionsResponse
// @Failure 400		{object}	responses.Error
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/stats/submission/time [get]
func (supervisorHandler *SupervisorHandler) GetSubmissionsOverTime(c echo.Context) error {
	groupId := c.QueryParam("group_id")
	year := c.QueryParam("year")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	if groupId != "" {

		groupUuid, err := uuid.Parse(groupId)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid Group ID.")
		}

		var classSchedule models.ClassSchedule
		classScheduleRepo := repositories.NewClassScheduleRepository(supervisorHandler.server.DB)
		classScheduleRepo.GetClassScheduleByGroupID(&classSchedule, groupUuid)

		if classSchedule.GroupID == uuid.Nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Group.")
		}

		classLabStaffRepo := repositories.NewClassLabStaffRepository(supervisorHandler.server.DB)
		if !utils.ValidateSupervisorAndBeyonder(existUser) {
			if !(classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID)) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}

		if *existUser.Role == constants.Role.Supervisor &&
			*classSchedule.SupervisorID != existUser.UserID {
			if !classLabStaffRepo.CheckStaffValidInClass(groupUuid, existUser.UserID) {
				return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
			}
		}
	}

	today := time.Now()
	startDate := today.AddDate(0, 0, -6)

	var submissionsList []int64
	var dateList []string

	submissionRepo := repositories.NewExerciseSubmissionRepository(supervisorHandler.server.DB)
	for d := startDate; !d.After(today); d = d.AddDate(0, 0, 1) {
		var submissions int64
		submissionRepo.GetSubmissionsByDate(
			&submissions,
			groupId,
			year,
			d,
		)
		dateList = append(dateList, d.Format("2006-01-02"))
		submissionsList = append(submissionsList, submissions)
	}

	response := responses.NewStatsSubmissionsResponse(submissionsList, dateList)
	return responses.Response(c, http.StatusOK, response)
}

// @Description Get Average Department Score
// @ID supervisor-get-average-departments-score
// @Tags Supervisor
// @Accept json
// @Produce json
// @Param year query string false "Year"
// @Success 200		{array}		responses.AverageDeptScoreResponse
// @Failure 403		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/supervisor/average_dept_score [get]
func (supervisorHandler *SupervisorHandler) GetAverageDeptScore(c echo.Context) error {
	year := c.QueryParam("year")

	userRepository := repositories.NewUserRepository(supervisorHandler.server.DB)
	existUser, err := utils.GetUserClaims(c, *userRepository)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, err.Error())
	}

	if !utils.ValidateSupervisorAndBeyonder(existUser) {
		var rolePermission []models.RolePermission
		rolePermissionRepo := repositories.NewRolePermissionRepository(supervisorHandler.server.DB)
		rolePermissionRepo.GetPermissionByRole(&rolePermission, *existUser.Role)

		if !utils.ValidateRolePermission(rolePermission, constants.PermissionType.DashboardAdmin) {
			return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Permission.")
		}
	}

	var department []models.DepartmentWithAggregate
	departmentRepo := repositories.NewDepartmentRepository(supervisorHandler.server.DB)
	err = departmentRepo.GetAllDeptsWithTotalMarks(&department, year)

	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error While Qurey.")
	}

	response := responses.NewAverageDeptScoreResponse(department)

	return responses.Response(c, http.StatusOK, response)
}
