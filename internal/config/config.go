package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type AppConfig struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
}

var Config *AppConfig

func Load() error {
	viper.AutomaticEnv()
	env := strings.ToLower(viper.GetString("APP_ENV"))
	if env == "" {
		env = "dev"
	}

	file := viper.New()
	file.SetConfigName("config")
	file.SetConfigType("yaml")
	file.AddConfigPath("configs")
	file.AddConfigPath("../..")
	if err := file.ReadInConfig(); err != nil {
		return fmt.Errorf("read configs: %w", err)
	}
	sub := file.Sub(env)
	if sub == nil {
		return fmt.Errorf("no %q section in configs.yaml", env)
	}

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for _, k := range []string{
		"db.host", "db.port", "db.name",
		"db.user", "db.password", "db.sslmode",
	} {
		_ = v.BindEnv(k)
	}

	v.MergeConfigMap(sub.AllSettings())

	var c AppConfig
	if err := v.Unmarshal(&c); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	Config = &c
	return nil
}
