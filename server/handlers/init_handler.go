package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	s "github.com/Project-IPCA/ipca-backend/server"
	"github.com/Project-IPCA/ipca-backend/services/department"
	labclassinfo "github.com/Project-IPCA/ipca-backend/services/lab_class_info"
	"github.com/Project-IPCA/ipca-backend/services/supervisor"
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
	deptNames := [6]string{"คอมพิวเตอร์", "ไฟฟ้า", "โยธา", "เคมี", "เครื่องกล", "อุตสาหการ"}
	departmentService := department.NewDepartmetService(initHanlder.server.DB)

	var wg sync.WaitGroup
	errChan := make(chan error, len(deptNames))

	for _, deptName := range deptNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			if err := departmentService.Create(name); err != nil {
				errChan <- fmt.Errorf("failed to create department %s: %v", name, err)
			}
		}(deptName)
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
		"oot1234",
		"oot1234",
		"Noppo",
		"Mummum",
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
func (initHandler *InitHandler) InitClassInfo (c echo.Context) error {
	labClassInfoService := labclassinfo.NewLabClassInfoService(initHandler.server.DB)
	labClassInfoService.Create(
		0,
		"Introduction",
		10,
		5,
	)
	return responses.MessageResponse(c, http.StatusOK, "Init Lab Class Info Success.")
}
