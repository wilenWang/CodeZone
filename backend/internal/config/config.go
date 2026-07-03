package config

import "os"

type Config struct {
	Port           string
	MySQLDSN       string
	SessionSecret  string
	CORSOrigin     string
	DevSeed        bool
	EnableDevLogin bool
}

func Load() Config {
	return Config{
		Port:           env("PORT", "8080"),
		MySQLDSN:       env("MYSQL_DSN", "chat:chat@tcp(127.0.0.1:3306)/chat?parseTime=true&multiStatements=true"),
		SessionSecret:  env("SESSION_SECRET", "dev-session-secret-change-me"),
		CORSOrigin:     env("CORS_ORIGIN", "http://localhost:5173"),
		DevSeed:        env("DEV_SEED", "true") == "true",
		EnableDevLogin: env("ENABLE_DEV_LOGIN", "false") == "true",
	}
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
