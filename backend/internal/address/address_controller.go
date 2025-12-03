package address

import (
	"backend-sistem06.com/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	Log     *logrus.Logger
	UseCase *AddressUseCase
}

func NewAddressController(usecase *AddressUseCase, log *logrus.Logger) *AddressController {
	return &AddressController{
		Log:     log,
		UseCase: usecase,
	}
}

func (c *AddressController) Create(ctx *fiber.Ctx) error {
	request := new(AddressRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	res, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create address: %+v", err)
		return err
	}
	return ctx.JSON(pkg.WebResponse[*Address]{Data: res})
}
