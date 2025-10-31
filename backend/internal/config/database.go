package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func NewPostgres(config *viper.Viper) *sql.DB {
	username := config.GetString("database.username")
	password := config.GetString("database.password")
	port := config.GetInt("database.port")
	host := config.GetString("database.host")
	dbname := config.GetString("database.name")
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Connection pool configuration
	db.SetMaxIdleConns(config.GetInt("database.pool.idle"))
	db.SetMaxOpenConns(config.GetInt("database.pool.max"))
	db.SetConnMaxLifetime(
		time.Duration(config.GetInt("database.pool.lifetime")) * time.Second,
	)

	// Test koneksi ke DB
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
