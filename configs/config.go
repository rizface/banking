package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbName     string
	DbPort     string
	DbHost     string
	DbUsername string
	DbPassword string
	DbParams   string
	ENV        string

	PrometheusAddress string

	JWTSecret  string
	BcryptSalt int

	S3ID         string
	S3SecretKey  string
	S3BucketName string
	S3Region     string
}

func LoadConfig() (Config, error) {
	config := Config{
		DbName:     os.Getenv("DB_NAME"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbParams:   os.Getenv("DB_PARAMS"),

		ENV: os.Getenv("ENV"),

		PrometheusAddress: os.Getenv("PROMETHEUS_ADDRESS"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		S3ID:         os.Getenv("S3_ID"),
		S3SecretKey:  os.Getenv("S3_SECRET_KEY"),
		S3BucketName: os.Getenv("S3_BUCKET_NAME"),
		S3Region:     os.Getenv("S3_REGION"),
	}

	salt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		return Config{}, fmt.Errorf("failed get bcrypt salt %v", err)
	}

	config.BcryptSalt = salt

	return config, nil
}
