package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
	"github.com/Project-IPCA/ipca-backend/services/user"
	userstudent "github.com/Project-IPCA/ipca-backend/services/user_student"
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
// @Router			/api/supervisor/add_students [post]
func (supervisorHandler *SupervisorHandler) AddStudents(c echo.Context) error {
	addStudentsReq := new(requests.AddStudentsTextRequest)

	if err := c.Bind(addStudentsReq); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}

	arrStudents := strings.Split(addStudentsReq.StudentsData, "\n")

	userService := user.NewUserService(supervisorHandler.server.DB)
	userStudentService := userstudent.NewUserStudentService(supervisorHandler.server.DB)
	userStudentRepository := repositories.NewUserStudentRepository(supervisorHandler.server.DB)

	for _, item := range arrStudents {
		data := strings.Split(item, " ")
		stuId := data[1]
		role := constants.Role.Student
		firstName := data[2]
		lastName := data[3]

		existUserStudent := models.UserStudent{}
		userStudentRepository.GetUserByStuID(&existUserStudent, stuId)

		if existUserStudent.StuStuID == stuId {
			return responses.ErrorResponse(
				c,
				http.StatusBadRequest,
				"User Student is Already Exist.",
			)
		}

		userId, _ := userService.Create(stuId, role)
		userStudentService.Create(userId, stuId, firstName, lastName)
	}

	return responses.MessageResponse(c, http.StatusCreated, "Add Student Successful")
}
