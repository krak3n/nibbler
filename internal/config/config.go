package config

import "os"

// DB Configuration
var (
	DBUser     = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName     = os.Getenv("DB_NAME")
	DBHost     = os.Getenv("DB_HOST")
	DBSSLMode  = os.Getenv("DB_SSL_MODE")
)
