package http

import (
	"sistem-06-Backend/internal/dto"
	"sistem-06-Backend/internal/usecase"
	"sistem-06-Backend/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log     *logrus.Logger
	UseCase *usecase.AuthUseCase
	Session *pkg.SessionHandler
}

func NewAuthController(useCase *usecase.AuthUseCase, log *logrus.Logger, session *pkg.SessionHandler) *AuthController {
	return &AuthController{
		Log:     log,
		UseCase: useCase,
		Session: session,
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(dto.UserLoginRequest)
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

	return ctx.JSON(pkg.WebResponse[*dto.UserResponse]{Data: response})

}
