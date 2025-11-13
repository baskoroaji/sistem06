package config

import (
	"database/sql"

	"backend-sistem06.com/internal/http"
	"backend-sistem06.com/internal/http/route"
	"backend-sistem06.com/internal/repository"
	"backend-sistem06.com/internal/usecase"
	"backend-sistem06.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	App       *fiber.App
	Config    *viper.Viper
	Log       *logrus.Logger
	DB        *sql.DB
	Validator *validator.Validate
	session   *session.Store
}

func Bootstrap(config *BootstrapConfig) {
	helloContoller := http.NewHelloController()

	userRepository := repository.NewUserRepository(config.DB, config.Log)
	sessionHandler := utils.NewSessionHandler(config.session, config.Log)
	// authRepository := repository.NewTokenRepository(config.DB, config.Log)

	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validator, userRepository)
	authUseCase := usecase.NewAuthUseCase(config.DB, config.Log, config.Validator, userRepository, config.session)

	userController := http.NewUserController(userUseCase, config.Log)
	authController := http.NewAuthController(authUseCase, config.Log, sessionHandler)

	routeConfig := route.RouteConfig{
		App:             config.App,
		HelloController: helloContoller,
		UserController:  userController,
		AuthController:  authController,
	}
	routeConfig.Setup()

}
