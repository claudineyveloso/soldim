package configs

import (
	"os"
	"strconv"
)

type AddressBaseURL struct {
	BaseURL string
}

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = InitConfig()

func InitConfig() Config {
	return Config{
		PublicHost:             os.Getenv("PUBLIC_HOST"),
		Port:                   os.Getenv("DB_PORT"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_NAME"),
		JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

func NewConfig(env string) *AddressBaseURL {
	var baseURL string
	switch env {
	case "dev":
		baseURL = "http://localhost:8080"
	case "prod":
		baseURL = "https://meuservidor.com" // Substitua pelo URL de produção
	default:
		baseURL = "http://localhost:8080" // Padrão para dev
	}

	return &AddressBaseURL{
		BaseURL: baseURL,
	}
}
