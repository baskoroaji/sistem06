package ports

import (
	"context"
	"sistem-06-Backend/internal/domain/entity"
)

type AddressRepository interface {
	CreateAddress(ctx context.Context, address *entity.Address) error
}
