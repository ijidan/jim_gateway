package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
)

const EnvLocal = "local"
const EnvTest = "test"
const EnvStage = "stage"
const EnvProduction = "production"

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Ver  string `yaml:"ver"`
		Env  string `yaml:"env"`
	}
	Websocket struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
		Log  string `yaml:"log"`
	}
	Tcp struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
		Log  string `yaml:"log"`
	}
	Rpc struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
		Ttl  int64  `yaml:"ttl"`
		Log  string `yaml:"log"`
	}
	Jaeger struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
	}
	PubSub struct{
		Brokers []string `yaml:"brokers"`
	}
	Gateway struct{
		Id string `yaml:"id"`
	}
	Runtime struct{
		Mode string `yaml:"mode"`
	}
}

var (
	onceConfig     sync.Once
	instanceConfig *Config
)

func GetConfigInstance(root string) *Config {
	onceConfig.Do(func() {
		instanceConfig = &Config{}
		v := viper.New()
		v.AddConfigPath(root)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.WatchConfig()
		v.OnConfigChange(func(in fsnotify.Event) {
		})
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := v.Unmarshal(instanceConfig); err != nil {
			panic(err)
		}
	})
	return instanceConfig
}
