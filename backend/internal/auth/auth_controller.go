package auth

import (
	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/internal/usecase"
	"backend-sistem06.com/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log     *logrus.Logger
	UseCase *usecase.AuthUseCase
	Session *utils.SessionHandler
}

func NewAuthController(useCase *usecase.AuthUseCase, log *logrus.Logger, session *utils.SessionHandler) *AuthController {
	return &AuthController{
		Log:     log,
		UseCase: useCase,
		Session: session,
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	if err := c.Session.SetUserSession(ctx, response.ID, response.Email); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})

}
