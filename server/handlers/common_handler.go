package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

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

// @Description Update User Info
// @ID supervisor-update-user-info
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
	existUser := utils.GetUserClaims(c, *userRepository)

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

	response := responses.NewUserInfoResponse(existUser)
	return responses.Response(c, http.StatusOK, response)
}
