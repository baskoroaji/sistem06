package repository

import (
	"database/sql"
	"fmt"
	"time"

	"backend-sistem06.com/internal/entity"
	"github.com/sirupsen/logrus"
)

type TokenRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewTokenRepository(db *sql.DB, log *logrus.Logger) *TokenRepository {
	return &TokenRepository{
		DB:  db,
		Log: log,
	}
}

func (r *TokenRepository) CreateToken(tx *sql.Tx, token *entity.PersonalAccessToken) error {
	now := time.Now().Unix()
	query :=
		`
	INSERT INTO personal_access_token (user_id, personal_access_token, created_at, expired_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(query, token.UserID, token.Token, now, now)
	} else {
		row = r.DB.QueryRow(query, token.UserID, token.Token, now, now)
	}

	err := row.Scan(&token.ID)
	if err != nil {
		r.Log.Errorf("failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.Log.Infof("user created successfully with ID: %d", token.ID)
	return nil

}

func (r *TokenRepository) FindTokenById(id int) (*entity.PersonalAccessToken, error) {

	query :=
		`
	SELECT user_id, personal_access_token, created_at, expired_at FROM personal_access_token
	WHERE id = $1
	`
	row := r.DB.QueryRow(query, id)

	var token entity.PersonalAccessToken
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.CreatedAt, &token.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Token not found")
		}
		return nil, err
	}

	return &token, nil
}
