package config

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewSession(viper *viper.Viper, db *sql.DB, log *logrus.Logger) *session.Store {
	// connectionURI := viper.GetString("database.url")
	tableName := viper.GetString("session.table")

	if tableName == "" {
		tableName = "sessions"
		log.Warn("session.table not configured, using default: sessions")
	}

	gcInterval := viper.GetDuration("session.gc_interval_minutes")
	if gcInterval == 0 {
		gcInterval = 10
	}

	expiration := viper.GetDuration("session.expiration_hours")
	if expiration == 0 {
		expiration = 24
	}

	log.Infof("Initializing session storage (table: %s)", tableName)

	storage := postgres.New(postgres.Config{
		Host:       viper.GetString("database.host"),
		Username:   viper.GetString("database.username"),
		Password:   viper.GetString("database.password"),
		Port:       viper.GetInt("database.port"),
		Database:   viper.GetString("database.name"),
		Reset:      false,
		GCInterval: gcInterval * time.Minute,
	})

	store := session.New(session.Config{
		Storage:        storage,
		Expiration:     expiration * time.Hour,
		KeyLookup:      "cookie:session_id",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	log.Info("Session store initialized successfully")
	return store
}
