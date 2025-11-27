package model

type Address struct {
	ID         int    `json:"id"`
	Jalan      string `json:"jalan"`
	RT         string `json:"RT"`
	RW         string `json:"RW"`
	Kota       string `json:"Kota"`
	PostalCode string `json:"PostalCode"`
}

type AddressRequest struct {
	Jalan      string `json:"jalan" validate:"required"`
	RT         string `json:"RT" validate:"required,RT_RW"`
	RW         string `json:"RW" validate:"required,RT_RW"`
	Kota       string `json:"Kota" validate:"required"`
	PostalCode string `json:"PostalCode" validate:"required,postal_code"`
}
