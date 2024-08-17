package configs

import (
	"os"
	"strconv"
)

type AddressBaseURL struct {
	BaseURL string
}

type Config struct {
	Environment            string
	Host                   string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = InitConfig()

func InitConfig() Config {
	env := getEnv("ENVIRONMENT", "dev")
	var host, port, dbUser, dbPassword, dbName string

	if env == "prod" {
		host = getEnv("PROD_DB_HOST", "c1i13pt05ja4ag.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com")
		port = getEnv("PROD_DB_PORT", "5432")
		dbUser = getEnv("PROD_DB_USER", "uatd3fq2t594ii")
		dbPassword = getEnv("PROD_DB_PASSWORD", "p4377b619ff8fd53aaf5c6ab44ae9bf0e72e7abd09ce2a4b1f1c5883d7b10f689")
		dbName = getEnv("PROD_DB_NAME", "d2jmp4epp3r1qp")
	} else {
		host = getEnv("DEV_DB_HOST", "localhost")
		port = getEnv("DEV_DB_PORT", "5432")
		dbUser = getEnv("DEV_DB_USER", "user_soldim_development")
		dbPassword = getEnv("DEV_DB_PASSWORD", "pwd_soldim_development")
		dbName = getEnv("DEV_DB_NAME", "soldim_development")
	}

	// return Config{
	// 	PublicHost:             os.Getenv("PUBLIC_HOST"),
	// 	Port:                   os.Getenv("DB_PORT"),
	// 	DBUser:                 os.Getenv("DB_USER"),
	// 	DBPassword:             os.Getenv("DB_PASSWORD"),
	// 	DBName:                 os.Getenv("DB_NAME"),
	// 	JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
	// 	JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
	// }

	return Config{
		Environment:            env,
		Host:                   host,
		Port:                   port,
		DBUser:                 dbUser,
		DBPassword:             dbPassword,
		DBName:                 dbName,
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
