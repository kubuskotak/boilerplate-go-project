package config

import (
	"sync"
	"time"

	"github.com/kubuskotak/valkyrie"
	"github.com/rs/zerolog/log"
)

// Config struct generate
type Config struct {
	App struct {
		Name         string         `yaml:"name"`
		Latency      int            `yaml:"latency"`
		ReadTimeout  int            `yaml:"read_timeout"`
		WriteTimeout int            `yaml:"write_timeout"`
		Timezone     string         `yaml:"timezone"`
		Debug        bool           `yaml:"debug"`
		Env          string         `yaml:"env"`
		SecretKey    string         `yaml:"secret_key"`
		ExpireIn     *time.Duration `yaml:"expire_in"`
	} `yaml:"App"`

	Port struct {
		Http int `yaml:"http"`
		Grpc int `yaml:"grpc"`
	} `yaml:"Ports"`

	DB struct {
		DsnMain string `yaml:"dsn_main" env:"DSN_MAIN"`
	}
}

var (
	Instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		Instance = &Config{}
		if err := valkyrie.Config(valkyrie.ConfigOpts{
			Config:    Instance,
			Paths:     []string{"./config"},
			Filenames: []string{"app.config.yaml", ".env"},
		}); err != nil {
			log.Error().Err(err).Msg("get config error")
		}
	})
	return Instance
}
