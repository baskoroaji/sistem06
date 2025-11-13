package config

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewPostgres(config *viper.Viper, log *logrus.Logger) *sql.DB {
	URI := config.GetString("database.url")
	// username := config.GetString("database.username")
	// password := config.GetString("database.password")
	// port := config.GetInt("database.port")
	// host := config.GetString("database.host")
	// dbname := config.GetString("database.name")
	// dsn := fmt.Sprintf(
	// 	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	host,
	// 	port,
	// 	username,
	// 	password,
	// 	dbname,
	// )

	db, err := sql.Open("postgres", URI)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Connection pool configuration
	db.SetMaxIdleConns(config.GetInt("database.pool.idle"))
	db.SetMaxOpenConns(config.GetInt("database.pool.max"))
	db.SetConnMaxLifetime(
		time.Duration(config.GetInt("database.pool.lifetime")) * time.Second,
	)

	// Test koneksi ke DB
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}
