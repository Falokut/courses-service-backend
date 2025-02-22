package main

import (
	"courses-service/assembly"
	"log"

	"github.com/Falokut/go-kit/app"
	"github.com/Falokut/go-kit/shutdown"
)

// @title					courses-service
// @version					1.0.0
// @description				Сервис для записи на курсы
// @BasePath				/api/courses-service
//
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
//
//go:generate swag init --parseDependency
//go:generate rm -f docs/swagger.json docs/docs.go
func main() {
	app, err := app.New()
	if err != nil {
		log.Println("error while creating app: ", err.Error())
		return
	}
	logger := app.GetLogger()

	assembly, err := assembly.New(app.Context(), logger)
	if err != nil {
		logger.Fatal(app.Context(), err)
	}
	app.AddRunners(assembly.Runners()...)
	app.AddClosers(assembly.Closers()...)

	err = app.Run()
	if err != nil {
		app.Shutdown()
		logger.Fatal(app.Context(), err)
	}

	shutdown.On(func() {
		logger.Info(app.Context(), "starting shutdown")
		app.Shutdown()
		logger.Info(app.Context(), "shutdown completed")
	})
}
