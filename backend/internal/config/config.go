package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret          string
	JWTAccessExpMinutes  int
	JWTRefreshExpDays    int

	AllowedOrigins string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	accessExp, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXP_MINUTES", "15"))
	refreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXP_DAYS", "7"))

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "financeku"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:          getEnv("JWT_SECRET", "change-me-in-production"),
		JWTAccessExpMinutes:  accessExp,
		JWTRefreshExpDays:    refreshExp,

		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
	}, nil
}

func (c *Config) DatabaseURL() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
