package address

import (
	"backend-sistem06.com/internal/model"
)

func AddressToResponse(address *Address) *model.Address {
	return &model.Address{
		ID:         address.ID,
		Jalan:      address.Jalan,
		RT:         address.RT,
		RW:         address.RW,
		Kota:       address.Kota,
		PostalCode: address.PostalCode,
	}
}
