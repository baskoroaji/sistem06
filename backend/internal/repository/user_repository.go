package repository

import (
	"database/sql"
	"fmt"
	"time"

	"backend-sistem06.com/internal/entity"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewUserRepository(db *sql.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		DB:  db,
		Log: log,
	}
}

func (r *UserRepository) CreateUser(tx *sql.Tx, user *entity.UserEntity) error {
	now := time.Now().Unix()
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRow(query, user.Name, user.Email, user.Password, now, now)
	} else {
		row = r.DB.QueryRow(query, user.Name, user.Email, user.Password, now, now)
	}

	err := row.Scan(&user.ID)
	if err != nil {
		r.Log.Errorf("failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.Log.Infof("user created successfully with ID: %d", user.ID)
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

func (r *UserRepository) FindByEmail(email string) (*entity.UserEntity, error) {
	query :=
		`
	SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := r.DB.QueryRow(query, email)

	var user entity.UserEntity
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
