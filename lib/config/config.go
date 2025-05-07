package config

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Storage struct {
	Path string `yaml:"path" env-required:"false" env-default:"./storage.db"`
}

type Smtp struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type FrontendUrls struct {
	Verify string `yaml:"verify"`
}

type Frontend struct {
	Path string       `yaml:"path" env-required:"false" env-default:"./frontend/build/index.html"`
	Urls FrontendUrls `yaml:"urls"`
}

type Server struct {
	Port int    `yaml:"port" env-required:"false" env-default:"8888"`
	Host string `yaml:"host" env-required:"false" env-default:"localhost"`
}

type Api struct {
	Name    string `yaml:"name" env-required:"false" env-default:"My API"`
	Version string `yaml:"version" env-required:"false" env-default:"1.0.0"`
}
type Secret struct {
	Jwt string `yaml:"jwt" env-required:"false" env-default:"x4FdvXiOwcP65Jc1VJyMbpun"`
}

type Config struct {
	Frontend `yaml:"frontend"`
	Secret   `yaml:"secret"`
	Server   `yaml:"server"`
	Storage  `yaml:"storage"`
	Api      `yaml:"api"`
	Smtp     `yaml:"smtp" env-required:"false"`
}

var (
	once         sync.Once   //nolint:gochecknoglobals // desired behavior
	globalConfig = &Config{} //nolint:gochecknoglobals // desired behavior
)

func Get() Config {
	once.Do(func() {
		globalConfig = MustLoad()
	})
	return *globalConfig
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	if res == "" {
		res = "config.yaml"
	}

	return res
}
