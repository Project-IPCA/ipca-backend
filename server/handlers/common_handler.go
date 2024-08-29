package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/pkg/utils"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
)

type CommonHandler struct {
	server *s.Server
}

func NewCommonHandler(server *s.Server) *CommonHandler {
	return &CommonHandler{server: server}
}

// @Description Get User Info
// @ID supervisor-get-user-info
// @Tags Common
// @Accept json
// @Produce json
// @Success 200		{object}	responses.UserInfoResponse
// @Security BearerAuth
// @Router			/api/common/user_info [get]
func (commonHandler *CommonHandler) GetUserInfo(c echo.Context) error {
	userRepository := repositories.NewUserRepository(commonHandler.server.DB)
	existUser := utils.GetUserClaims(c, *userRepository)
	response := responses.NewUserInfoResponse(existUser)
	return responses.Response(c, http.StatusOK, response)
}
