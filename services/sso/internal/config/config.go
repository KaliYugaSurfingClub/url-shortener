package config

import "time"

type JWT struct {
	Secret    string `yaml:"secret"`
	UserIdKey string `yaml:"user_id_key"`
}

type Cookie struct {
	Name     string        `yaml:"name"`
	Lifetime time.Duration `yaml:"lifetime"`
}

type HTTP struct {
	Port string `yaml:"port"`
	IP   string `yaml:"ip"`
}

type Config struct {
	JWT    JWT    `yaml:"jwt"`
	HTTP   HTTP   `yaml:"http"`
	Cookie Cookie `yaml:"cookie"`
}

func MustLoad() *Config {
	return &Config{
		JWT: JWT{
			Secret:    "secret",
			UserIdKey: "user_id",
		},
		Cookie: Cookie{
			Name:     "user_id",
			Lifetime: 24 * time.Hour,
		},
		HTTP: HTTP{
			Port: "8081",
			IP:   "0.0.0.0",
		},
	}
}
