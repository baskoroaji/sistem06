package address

import (
	"database/sql"
	"fmt"
	"time"

	"backend-sistem06.com/internal/entity"
	"github.com/sirupsen/logrus"
)

type AddressRepositoryInterface interface {
	CreateAddress(tx *sql.Tx, address *entity.Address) error
}

type AddressRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewAddressRepository(db *sql.DB, log *logrus.Logger) *AddressRepository {
	return &AddressRepository{
		DB:  db,
		Log: log,
	}
}

func (r *AddressRepository) CreateAddress(tx *sql.Tx, address *entity.Address) error {
	now := time.Now().Unix()
	query :=
		`
	INSERT INTO address (jalan, rt, rw, kota, postal_code, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRow(query, address.Jalan, address.RT, address.RW, address.Kota, address.PostalCode, now, now)
	} else {
		row = r.DB.QueryRow(query, address.Jalan, address.RT, address.RW, address.Kota, address.PostalCode, now, now)
	}

	err := row.Scan(&address.ID)
	if err != nil {
		r.Log.Errorf("failed to create address: %v", err)
		return fmt.Errorf("failed to create address: %v", err)
	}
	r.Log.Infof("address added successfully with ID: %d", address.ID)
	return nil
}

func (r *AddressRepository) UpdateAddress(tx *sql.Tx, addressID int, change map[string]interface{}) error {
	if len(change) == 0 {
		return fmt.Errorf("no changes in field")
	}

	//Maybe i will refactor and use sqlx
}
