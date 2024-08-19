package configs

import (
	"log"
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
	log.Println(`ENVIRONMENT ${env}`)

	if env == "prod" {
		host = getEnv("PROD_DB_HOST", "c3gtj1dt5vh48j.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com")
		port = getEnv("PROD_DB_PORT", "5432")
		dbUser = getEnv("PROD_DB_USER", "u8v8gfju2nbfqn")
		dbPassword = getEnv("PROD_DB_PASSWORD", "p009827e7e08f28b44b9ba56751f82543345f445a718fb44624dd5b653e0238bd")
		dbName = getEnv("PROD_DB_NAME", "ddcvr5rele9132")
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
