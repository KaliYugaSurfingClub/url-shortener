package config

import (
	"time"
)

type Config struct {
	PostgresURL string
	HTTPServer  HTTPServer
	Service     Service
	Auth        Auth
}

type HTTPServer struct {
	IP                string
	Port              string
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration
}

type Service struct {
	GeneratedAliasLength int
	TriesToGenerate      int
	Alp                  string
}

type Auth struct {
	JWTSecret  string
	JWTKey     string
	CookieName string
}

//todo

func MustLoad() *Config {
	cfg := &Config{
		PostgresURL: "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable",
		HTTPServer: HTTPServer{
			Port:              "8080",
			IP:                "0.0.0.0",
			WriteTimeout:      time.Second * 15,
			ReadHeaderTimeout: time.Second * 15,
			ReadTimeout:       time.Second * 15,
			IdleTimeout:       time.Second * 60,
		},
		Service: Service{
			GeneratedAliasLength: 4,
			TriesToGenerate:      3,
			Alp:                  "abc123",
		},
		Auth: Auth{
			JWTSecret:  "secret",
			JWTKey:     "user_id",
			CookieName: "user_id",
		},
	}

	return cfg
}
