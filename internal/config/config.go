package config

import (
	"errors"
	"os"
)

type Config struct {
	Port            string
	DatabaseURL     string
	AWSRegion       string
	AWSS3BucketName string
}

func Load() (*Config, error) {
	databaseURL := getEnv("DATABASE_URL", "postgres://alex:wantam@localhost:5432/file_storage?sslmode=disable")

	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	if bucketName == "" {
		return nil, errors.New("AWS_S3_BUCKET_NAME environment variable is required")
	}

	return &Config{
		Port:            getEnv("PORT", ":8080"),
		DatabaseURL:     databaseURL,
		AWSRegion:       getEnv("AWS_REGION", "af-south-1"),
		AWSS3BucketName: bucketName,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
