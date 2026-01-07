package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Config struct {
	RunAddr             string `mapstructure:"run_address" env:"RUN_ADDRESS"`
	DataBaseURI         string `mapstructure:"database_uri" env:"DATABASE_URI"`
	AccrualSysAddr      string `mapstructure:"accrual_system_address" env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecretKey        string `mapstructure:"jwt_secret_key" env:"JWT_SECRET_KEY"`
	Workers             int    `mapstructure:"workers" env:"WORKERS"`
	WorkerTimeout       int    `mapstructure:"worker_timeout" env:"WORKER_TIMEOUT"`
	ProcessRetryTimeout int    `mapstructure:"process_retry_timeout" env:"PROCESS_RETRY_TIMEOUT"`
}

func New() *Config {
	return &Config{
		RunAddr:             "localhost:8080",
		AccrualSysAddr:      "http://localhost:8000",
		Workers:             3,
		WorkerTimeout:       5,
		ProcessRetryTimeout: 30,
	}
}

func (c *Config) ReadFromYaml() error {
	configFilename := flag.String("config", "", "config filename")
	flag.Parse()

	f, err := os.ReadFile(*configFilename)
	if err != nil {
		return err
	}

	var raw any

	if err := yaml.Unmarshal(f, &raw); err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &c})
	if err != nil {
		return err
	}

	if err := decoder.Decode(raw); err != nil {
		return err
	}

	return nil
}

func (c *Config) ReadFromFlags() {
	runAddr := flag.String("a", "", "listen addres")
	dataBaseURI := flag.String("d", "", "db uri")
	accrualSysAddr := flag.String("r", "", "accrual system address")
	flag.Parse()

	if runAddr != nil && *runAddr != "" {
		c.RunAddr = *runAddr
	}
	if dataBaseURI != nil && *dataBaseURI != "" {
		c.DataBaseURI = *dataBaseURI
	}
	if accrualSysAddr != nil && *accrualSysAddr != "" {
		c.AccrualSysAddr = *accrualSysAddr
	}
}

func (c *Config) ReadFromEnv() error {
	return cleanenv.ReadEnv(c)
}
