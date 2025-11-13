package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend-sistem06.com/internal/entity"
	"github.com/sirupsen/logrus"
)

type UserRepositoryInterface interface {
	CreateUser(tx *sql.Tx, user *entity.UserEntity) error
	CountById(tx *sql.Tx, id int) (int, error)
	CountByName(tx *sql.Tx, name string) (int, error)
	FindByEmail(ctx context.Context, email string) (*entity.UserWithRole, error)
	FindByID(id int) (*entity.UserEntity, error)
}
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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.UserWithRole, error) {
	query := `
			SELECT 
	u.id AS user_id,
	u.email,
	u.password,
	r.id AS role_id,
	r.name AS role_name,
	p.id AS permission_id,
	p.name AS permission_name
	FROM users u
	JOIN user_roles ur ON ur.user_id = u.id
	JOIN roles r ON r.id = ur.roles_id
	LEFT JOIN roles_permissions rp ON rp.role_id = r.id
	LEFT JOIN permissions p ON p.id = rp.permission_id
	WHERE u.email = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user *entity.UserWithRole
	roleMap := make(map[int64]*entity.Role)

	for rows.Next() {
		var (
			userID   int64
			email    string
			password string
			roleID   sql.NullInt64
			roleName sql.NullString
			permName sql.NullString
		)

		if err := rows.Scan(&userID, &email, &password, &roleID, &roleName, &permName); err != nil {
			return nil, err
		}

		if user == nil {
			user = &entity.UserWithRole{
				ID:       userID,
				Email:    email,
				Password: password,
				RoleName: []entity.Role{},
			}
		}

		if roleID.Valid {
			role, exists := roleMap[roleID.Int64]
			if !exists {
				role = &entity.Role{
					ID:          roleID.Int64,
					Name:        roleName.String,
					Permissions: []string{},
				}
				roleMap[roleID.Int64] = role
			}

			if permName.Valid && permName.String != "" {
				role.Permissions = append(role.Permissions, permName.String)
			}
		}
	}

	// pindahkan map ke slice
	for _, r := range roleMap {
		user.RoleName = append(user.RoleName, *r)
	}

	return user, nil

}

func (r *UserRepository) FindByID(id int) (*entity.UserEntity, error) {
	query :=
		`
	SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := r.DB.QueryRow(query, id)

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
