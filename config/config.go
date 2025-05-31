package config

import (
	"log"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Server Server
	Mongo  Mongo
	Bcrypt Bcrypt
	JWT    JWT
}

type Server struct {
	Port int `envconfig:"PORT" default:"8080"`
}

type Mongo struct {
	URI      string `envconfig:"MONGO_URI" default:"mongodb://localhost:27017"`
	Database string `envconfig:"DB_NAME" default:"user-api"`
}

type Bcrypt struct {
	SaltRounds int `mapstructure:"saltRounds"`
}

type JWT struct {
	SecretKey string `mapstructure:"secretKey"`
	Algorithm string `mapstructure:"algorithm"`
	ExpiresIn string `mapstructure:"expiresIn"`
}

var cfg Config

func Init() {
	runtime.GOMAXPROCS(1)
	_ = godotenv.Load()
	initViper()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unmarshal viper config error: %s", err)
	}

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("read env error: %s", err.Error())
	}
}

func Get() Config {
	return cfg
}

func initViper() {
	viper.AddConfigPath("..")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read viper config: %s", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
