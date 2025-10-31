package repository

import (
	"database/sql"
	"fmt"

	"backend-sistem06.com/internal/entity"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) CreateUser(tx *sql.Tx, user *entity.UserEntity) error {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var lastInsertID int64
	err := r.DB.QueryRow(query, user.Name, user.Email, user.Password).Scan(&lastInsertID)
	if err != nil {
		r.Log.Errorf("failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.Log.Infof("user created successfully with ID: %d", lastInsertID)
	return nil
}

func (r *UserRepository) CountById(tx *sql.Tx, id int) (int, error) {
	query :=
		`
	SELECT COUNT(*) FROM users WHERE id = $1
	`
	var total int
	err := tx.QueryRow(query, id).Scan(&total)
	return total, err
}

func (r *UserRepository) CountByName(tx *sql.Tx, name string) (int, error) {
	query :=
		`
	SELECT COUNT(*) FROM users WHERE name = $1
	`
	var totalName int
	err := tx.QueryRow(query, name).Scan(&totalName)
	return totalName, err
}
