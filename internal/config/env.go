package config

import (
	"fmt"
	"os"
	"strings"
)

var requiredEnvVars = []string{
	"APP_ENV",
	"APP_PORT",
	"DB_CONNECTION",
	"DB_DRIVER",
}

type EnviConfig struct {
	AppPort      string
	AppEnv       string
	DbDriver     string
	DbConnection string
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitAndCheckEnv() (*EnviConfig, []error) {
	var errs []error
	if !strings.EqualFold(strings.ToLower(os.Getenv("APP_ENV")), "local") {
		for _, v := range requiredEnvVars {
			if GetEnv(v, "") == "" {
				errs = append(errs, fmt.Errorf("%s is required", v))
			}
		}
	}

	return &EnviConfig{
		AppEnv:       GetEnv("APP_ENV", "local"),
		AppPort:      GetEnv("APP_PORT", "8010"),
		DbConnection: GetEnv("DB_CONNECTION", ""),
		DbDriver:     GetEnv("DB_DRIVER", "sqlite"),
	}, nil
}
