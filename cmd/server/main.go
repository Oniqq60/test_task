package main

import (
	"flag"
	"os"

	api "test_task/gen/task"
	"test_task/handlers"
	"test_task/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoMiddleware "github.com/oapi-codegen/echo-middleware"
)

func main() {
	port := flag.String("port", "8080", "HTTP port")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		os.Exit(1)
	}
	swagger.Servers = nil

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoMiddleware.OapiRequestValidator(swagger))

	tm := service.NewTaskManager(100)
	api.RegisterHandlers(e, handlers.NewTaskAPI(tm))

	e.Logger.Fatal(e.Start(":" + *port))
}
