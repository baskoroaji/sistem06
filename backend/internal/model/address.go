package model

type Address struct {
	ID         int    `json:"id"`
	Jalan      string `json:"jalan"`
	RT         int    `json:"RT"`
	RW         int    `json:"RW"`
	Kota       string `json:"Kota"`
	PostalCode int    `json:"PostalCode"`
}

type AddressRequest struct {
	Jalan      string `json:"jalan" validate:"required"`
	RT         string `json:"RT" validate:"required,RT_RW"`
	RW         string `json:"RW" validate:"required,RT_RW"`
	Kota       string `json:"Kota" validate:"required"`
	PostalCode string `json:"PostalCode" validate:"required,postal_code"`
}
