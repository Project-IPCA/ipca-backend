package main

import (
	"fmt"

	application "github.com/Project-IPCA/ipca-backend"
	"github.com/Project-IPCA/ipca-backend/config"
	"github.com/Project-IPCA/ipca-backend/docs"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
)

// @title IPCA API
// @version 1.0
// @description IPCA API
// @host petstore.swagger.io
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	cfg := config.NewConfig()
	if (cfg.Env.Enviroment == constants.EnviromentType.Develop) {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	}
	application.Start(cfg)
}
