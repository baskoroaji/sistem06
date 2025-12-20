package repository

import (
	"sistem-06-Backend/internal/database/sqlc"

	"github.com/sirupsen/logrus"
)

type UserRepositoryImpl struct {
	q *sqlc.Queries
	log *logrus.Logger
}

func NewUserRepository (q *sqlc.Querier, log *logrus.Logger) *
