package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/rabbitmq_client"
	"github.com/Project-IPCA/ipca-backend/redis_client"
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

func (testHandler *TestHandler) TestRedis(c echo.Context) error {
	 redis := redis_client.NewRedisAction(testHandler.server.Redis)
	 err := redis.PublishMessage("oot","handsome")
	
	if err != nil {
		panic(err)
	}
	return responses.MessageResponse(c, http.StatusOK, "Test Redis ok")
}

func (testHandler *TestHandler) TestRabbitMQ(c echo.Context) error {
	rabbit := rabbitmq_client.NewRabbitMQAction(testHandler.server.RabitMQ,testHandler.server.Config)
	test := map[string]interface{}{
        "FirstName": "John",
        "LastName":  "Doe",
        "Age":       30,
    }
	err := rabbit.SendQueue(test)
   
   if err != nil {
	   panic(err)
   }
   return responses.MessageResponse(c, http.StatusOK, "Test Rabbit OK")
}