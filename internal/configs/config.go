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
		host = getEnv("PROD_DB_HOST", "c9uss87s9bdb8n.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com")
		port = getEnv("PROD_DB_PORT", "5432")
		dbUser = getEnv("PROD_DB_USER", "urgpqrbhtjorj")
		dbPassword = getEnv("PROD_DB_PASSWORD", "p53c7b9a425f1437b8b0fea043bf700ec4dfe5c218633a2ec533700d4cea9cce7")
		dbName = getEnv("PROD_DB_NAME", "d5m311pvqgff2e")
	} else {
		host = getEnv("DEV_DB_HOST", "localhost")
		port = getEnv("DEV_DB_PORT", "5432")
		dbUser = getEnv("DEV_DB_USER", "user_soldim_development")
		dbPassword = getEnv("DEV_DB_PASSWORD", "pwd_soldim_development")
		dbName = getEnv("DEV_DB_NAME", "soldim_development")
	}

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
