package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                     string     `yaml:"env" env-default:"prod"`
	ClearStorageRefreshRate int        `yaml:"clear_storage_refresh_rate" env-default:"1"`
	MemSize                 uint64     `yaml:"mem_size" env-default:"500"`
	FilePath                string     `yaml:"file_path" env-default:"storage.json"`
	HTTPServer              HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	RWTimeout   time.Duration `yaml:"rw_timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	StopTimeout time.Duration `yaml:"stop_timeout" env-default:"10s"`
}

func NewConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	cfg := new(Config)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("error reading config: %s", err)

	}

	return cfg, nil
}
