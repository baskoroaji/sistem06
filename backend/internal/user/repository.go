package User

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
	FindByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	FindByID(id int) (*entity.UserEntity, error)
	FindWithRoles(ctx context.Context, userid int) (*entity.UserWithRole, error)
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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	query := `SELECT id, name, email, password, created_at, updated_at 
	          FROM users WHERE email = $1`

	var user entity.UserEntity
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
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

// Khusus untuk load roles & permissions setelah login
func (r *UserRepository) FindWithRoles(ctx context.Context, userID int) (*entity.UserWithRole, error) {
	// Query user basic info
	userQuery := `SELECT id, name, email, password
	              FROM users WHERE id = $1`

	var user entity.UserWithRole
	err := r.DB.QueryRowContext(ctx, userQuery, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Query roles & permissions
	rolesQuery := `
		SELECT 
			r.id AS role_id,
			r.name AS role_name,
			p.name AS permission_name
		FROM user_roles ur
		JOIN roles r ON r.id = ur.roles_id
		LEFT JOIN roles_permissions rp ON rp.role_id = r.id
		LEFT JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1
		ORDER BY r.id, p.name
	`

	rows, err := r.DB.QueryContext(ctx, rolesQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleMap := make(map[int]*entity.RolesWithPermissions)
	permissionSet := make(map[int]map[string]bool)

	for rows.Next() {
		var (
			roleID   int
			roleName string
			permName sql.NullString
		)

		if err := rows.Scan(&roleID, &roleName, &permName); err != nil {
			return nil, err
		}

		// Initialize role if not exists
		if _, exists := roleMap[roleID]; !exists {
			roleMap[roleID] = &entity.RolesWithPermissions{
				ID:          roleID,
				Name:        roleName,
				Permissions: []string{},
			}
			permissionSet[roleID] = make(map[string]bool)
		}

		// Add permission if valid and unique
		if permName.Valid && permName.String != "" {
			if !permissionSet[roleID][permName.String] {
				roleMap[roleID].Permissions = append(roleMap[roleID].Permissions, permName.String)
				permissionSet[roleID][permName.String] = true
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice
	user.Roles = make([]entity.RolesWithPermissions, 0, len(roleMap))
	for _, role := range roleMap {
		user.Roles = append(user.Roles, *role)
	}

	return &user, nil
}
