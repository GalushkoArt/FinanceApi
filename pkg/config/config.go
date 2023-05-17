package config

import (
	"FinanceApi/pkg/log"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Database struct {
		Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
		Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
		Name     string `yaml:"name" env:"DB_NAME" env-default:"finance_app"`
		User     string `yaml:"user" env:"DB_USER" env-default:"finance_user"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	API struct {
		TwelveData struct {
			Host      string        `yaml:"host" env:"TWELVE_DATA_HOST" env-default:"https://api.twelvedata.com"`
			Timeout   time.Duration `yaml:"timeout" env:"TWELVE_DATA_TIMEOUT" env-default:"2m"`
			RateLimit int           `yaml:"rateLimit" env:"TWELVE_DATA_RATE_LIMIT" env-default:"8"`
			ApiKey    string        `yaml:"apiKey" env:"TWELVE_DATA_API_KEY"`
		} `yaml:"twelveData"`
	} `yaml:"api"`
	Logs struct {
		Level string `yaml:"level" env:"LOGS_LEVEL" env-default:"INFO"`
		Path  string `yaml:"path" env:"LOGS_PATH" env-default:"logs.txt"`
	} `yaml:"logs"`
	Server struct {
		Port            string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
		ReadTimeoutSec  int    `yaml:"readTimeoutSec" env:"SERVER_READ_TIMEOUT" env-default:"10"`
		WriteTimeoutSec int    `yaml:"writeTimeoutSec" env:"SERVER_WRITE_TIMEOUT" env-default:"10"`
	} `yaml:"server"`
	Cache struct {
		SymbolTTL time.Duration `yaml:"symbolTtl" env:"CACHE_SYMBOL_TTL" env-default:"1h"`
	} `yaml:"cache"`
}

var Conf Config

func Init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Panic("Couldn't open current working directory!", err)
	}
	err = cleanenv.ReadConfig(filepath.Join(wd, "config.yaml"), &Conf)
	if err != nil {
		log.Panic("Error on reading config!", err)
	}
}
