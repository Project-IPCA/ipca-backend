package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	s "github.com/Project-IPCA/ipca-backend/server"
	"github.com/Project-IPCA/ipca-backend/services/department"
	labclassinfo "github.com/Project-IPCA/ipca-backend/services/lab_class_info"
	"github.com/Project-IPCA/ipca-backend/services/supervisor"
	"github.com/Project-IPCA/ipca-backend/services/ta"
	userservice "github.com/Project-IPCA/ipca-backend/services/user"
)

type InitHandler struct {
	server *s.Server
}

func NewInitHandler(server *s.Server) *InitHandler {
	return &InitHandler{server: server}
}

// @Description Init Department
// @ID init-department
// @Tags Init
// @Accept json
// @Produce json
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Router			/api/init/department [post]
func (initHanlder *InitHandler) InitDepartment(c echo.Context) error {
	deptNames := [6]requests.CreateDepartmentRequest{
		{
			Name:    "คอมพิวเตอร์",
			Name_EN: "Computer Engineering",
		},
		{
			Name:    "ไฟฟ้า",
			Name_EN: "Electrical Engineering",
		},
		{
			Name:    "โยธา",
			Name_EN: "Civil Engineering",
		},
		{
			Name:    "เคมี",
			Name_EN: "Chemical Engineering",
		},
		{
			Name:    "เครื่องกล",
			Name_EN: "Mechanical Engineering",
		},
		{
			Name:    "อุตสาหการ",
			Name_EN: "Industrial Engineering",
		}}
	departmentService := department.NewDepartmetService(initHanlder.server.DB)

	var wg sync.WaitGroup
	errChan := make(chan error, len(deptNames))

	for _, deptName := range deptNames {
		wg.Add(1)
		go func(name, name_en string) {
			defer wg.Done()
			if err := departmentService.Create(name, name_en); err != nil {
				errChan <- fmt.Errorf("failed to create department %s: %v", name, err)
			}
		}(deptName.Name, deptName.Name_EN)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return responses.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Failed to initialize departments",
			)
		}
	}

	return responses.MessageResponse(c, http.StatusOK, "Init department success.")
}

// @Description Init Supervisor
// @ID init-supervisor
// @Tags Init
// @Accept json
// @Produce json
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Router			/api/init/supervisor [post]
func (initHandler *InitHandler) InitSupervisor(c echo.Context) error {
	userService := userservice.NewUserService(initHandler.server.DB)
	supervisorService := supervisor.NewSupervisorService(initHandler.server.DB)

	userId, err := userService.CreateQuick(
		"test1",
		"test1",
		"test1",
		"test1",
		constants.Gender.Male,
		constants.Role.Supervisor,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Professor X is KING")
	}

	err = supervisorService.Create(userId, "คอมพิวเตอร์")
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Professor X is KING")
	}
	return responses.MessageResponse(c, http.StatusOK, "Init Supervisor Success.")
}

// @Description Init Lab Class Info
// @ID init-lab-class-info
// @Tags Init
// @Accept json
// @Produce json
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Router			/api/init/labclassinfo [post]
func (initHandler *InitHandler) InitClassInfo(c echo.Context) error {
	labClassInfoService := labclassinfo.NewLabClassInfoService(initHandler.server.DB)
	chapterName := [17]string{
		"Introduction",
		"Variables Expression Statement",
		"Conditional Execution",
		"'Loop while",
		"Loop for",
		"List",
		"String",
		"Function",
		"Dictionary",
		"Files",
		"Best Practice 1",
		"Best Practice 2",
		"Quiz #1 (chapter 01 - 03)",
		"Quiz #2 (chapter 01 - 06)",
		"Quiz #3 (chapter 01 - 09)",
		"Reserve #1",
		"Reserve #2",
	}

	for i := 0; i < len(chapterName); i++ {
		labClassInfoService.Create(
			i+1,
			chapterName[i],
			10,
			5,
		)
	}

	return responses.MessageResponse(c, http.StatusOK, "Init Lab Class Info Success.")
}

// @Description Init TA
// @ID init-ta
// @Tags Init
// @Accept json
// @Produce json
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Failure 500		{object}	responses.Error
// @Router			/api/init/ta [post]
func (initHandler *InitHandler) InitTA(c echo.Context) error {
	userService := userservice.NewUserService(initHandler.server.DB)
	taService := ta.NewTaService(initHandler.server.DB)

	userId, err := userService.CreateQuick(
		"ootTa",
		"ootTa",
		"TaOot",
		"Handsome",
		constants.Gender.Male,
		constants.Role.Ta,
	)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Professor X is KING")
	}

	err = taService.CreateTa(userId, nil, nil)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Professor X is KING")
	}
	return responses.MessageResponse(c, http.StatusOK, "Init Ta Success.")
}
