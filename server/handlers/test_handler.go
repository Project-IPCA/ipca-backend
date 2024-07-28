package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	s "github.com/Project-IPCA/ipca-backend/server"
)

type TestHandler struct {
	server *s.Server
}

func NewTestHandler(server *s.Server) *TestHandler {
	return &TestHandler{server: server}
}

// @Description Greeting
// @ID greeting
// @Tags Test
// @Accept json
// @Produce json
// @Success 200		{object}	responses.Data
// @Failure 404		{object}	responses.Error
// @Router			/api/greeting [get]
func (testHandler *TestHandler) Greeting(c echo.Context) error {
	return responses.MessageResponse(c, http.StatusOK, "Greeting OK")
}
