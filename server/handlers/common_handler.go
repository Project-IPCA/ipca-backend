package handlers

import (
	"fmt"
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
// @ID common-get-user-info
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

// @Description Get Keyword List
// @ID common-get-keyword-list
// @Tags Common
// @Accept json
// @Produce json
// @Param params body	requests.GetKeywordListRequest	true	"Get Keyword List Request"
// @Success 200		{object}	responses.Data
// @Failure 400		{object}	responses.Error
// @Router			/api/common/get_keyword_list [post]
func (commonHandler *CommonHandler) GetKeywordList(c echo.Context) error{
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
	keywordList,err := utils.GetKeywordFromCode(getKeywordListRequest.Sourcecode)
	if(err!=nil){
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Error While Running Sourcecode %s", err),
		)
	}
	return responses.Response(c,http.StatusOK,keywordList)
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
func (commonHandler *CommonHandler) KeywordCheck(c echo.Context) error{
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
	checkKeyword,err := utils.KeywordCheck(checkKeywordRequest.Sourcecode,checkKeywordRequest.ExerciseKeywordList)
	if(err!=nil){
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("Error While Running Sourcecode %s", err),
		)
	}
	return responses.Response(c,http.StatusOK,checkKeyword)
}