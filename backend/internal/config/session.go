package config

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/spf13/viper"
)

func NewSession(viper *viper.Viper) *session.Store {
	store := postgres.New(postgres.Config{
		ConnectionURI: viper.GetString("database.url"),
		Table:         viper.GetString("session.table"),
		Reset:         false,
		GCInterval:    viper.GetDuration("session.gc_interval_minutes") * time.Minute,
	})

	return session.New(session.Config{
		Storage:        store,
		Expiration:     viper.GetDuration("session.expiration_hours") * time.Hour,
		CookieHTTPOnly: true,
		CookieSecure:   false, // true jika pakai HTTPS
		CookieSameSite: "Lax",
	})
}
