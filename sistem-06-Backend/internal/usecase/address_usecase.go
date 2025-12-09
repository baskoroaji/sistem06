package usecase

import (
	"context"
	"database/sql"

	"backend-sistem06.com/internal/entity"
	"backend-sistem06.com/internal/model"
	"backend-sistem06.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressUseCase struct {
	DB                *sql.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	AddressRepository repository.AddressRepositoryInterface
}

func NewAddressUseCase(db *sql.DB, log *logrus.Logger, validate *validator.Validate, addressRepository repository.AddressRepositoryInterface) *AddressUseCase {
	return &AddressUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		AddressRepository: addressRepository,
	}
}

func (c *AddressUseCase) Create(ctx context.Context, request *model.AddressRequest) (*model.Address, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to begin transaction: %+v", err)

		if err == context.Canceled || err == context.DeadlineExceeded {
			return nil, fiber.NewError(fiber.StatusRequestTimeout, "request timeout or canceled")
		}

		return nil, fiber.ErrInternalServerError
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		validationErrors := utils.ValidationError(err)
		c.Log.Warnf("Validation failed: %+v", validationErrors)

		// Return response error yang bisa dibaca frontend
		return nil, fiber.NewError(fiber.StatusBadRequest, utils.FormatValidationErrors(validationErrors))
	}

	address := &entity.Address{
		Jalan:      request.Jalan,
		RT:         request.RT,
		RW:         request.RW,
		Kota:       request.Kota,
		PostalCode: request.PostalCode,
	}

	if err := c.AddressRepository.CreateAddress(tx, address); err != nil {
		c.Log.Warnf("Database insert error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit(); err != nil {
		c.Log.Warnf("failed commit transcation: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.AddressToResponse(address), nil
}
