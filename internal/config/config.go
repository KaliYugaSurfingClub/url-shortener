package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	PostgresURL string     `mapstructure:"postgres_url"`
	HTTPServer  HTTPServer `mapstructure:"http_server"`
	Service     Service    `mapstructure:"service"`
	Auth        Auth       `mapstructure:"auth"`
}

type HTTPServer struct {
	Address           string        `mapstructure:"address"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
}

type Service struct {
	GeneratedAliasLength int `mapstructure:"generated_alias_length"`
}

type Auth struct {
	JWTSecret       string        `mapstructure:"jwt_secret"`
	UserIdJWTKey    string        `mapstructure:"user_id_jwt_key"`
	UserIdCookieKey string        `mapstructure:"user_id_cookie_name"`
	CookieLifetime  time.Duration `mapstructure:"cookie_lifetime"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load .env file")
	}
}

func MustLoad() *Config {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		log.Fatal(".env file with CONFIG_PATH variable not found")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Yaml config file does not exist: %s", configPath)
	}

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return &cfg
}
