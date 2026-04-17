package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIId    string
	Login    string
	Password string
	SoapURL  string
	Port     string
}

func Load() (*Config, error) {
	// ignore error — .env file is optional
	_ = godotenv.Load()

	cfg := &Config{
		APIId:    os.Getenv("DD_API_ID"),
		Login:    os.Getenv("DD_LOGIN"),
		Password: os.Getenv("DD_PASSWORD"),
		SoapURL:  os.Getenv("DD_SOAP_URL"),
		Port:     os.Getenv("PORT"),
	}

	if cfg.APIId == "" {
		return nil, fmt.Errorf("DD_API_ID is required")
	}
	if cfg.Login == "" {
		return nil, fmt.Errorf("DD_LOGIN is required")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("DD_PASSWORD is required")
	}
	if cfg.SoapURL == "" {
		cfg.SoapURL = "http://www.drebedengi.ru/soap/"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
