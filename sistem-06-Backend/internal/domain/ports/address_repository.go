package domain

import "context"

type AddressRepository interface {
	CreateAddress(ctx context.Context, address Address) error
}
