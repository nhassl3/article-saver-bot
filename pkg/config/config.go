package config

import (
	"github.com/joho/godotenv"
	"github.com/nhassl3/article-saver-bot/pkg/e"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server `yaml:"server"`
	Token  string
}

type Server struct {
	Protocol string `yaml:"protocol" env-default:"https"`
	Host     string `yaml:"host" env-default:"api.telegram.org"`
}

func (c *Config) MustLoad() error {
	if err := godotenv.Load(); err != nil {
		return e.Wrap("error to load env", nil)
	}

	var token string
	if err := GetVariable("TELEGRAM_BOT_TOKEN", &token); err != nil {
		return err
	}
	c.Token = token

	var configPath string
	if err := GetVariable("CONFIG_PATH", &configPath); err != nil {
		return err
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return e.Wrap("err to load config file", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	return e.WrapIfErr("err to unmarshal config file", err)
}

func GetVariable(key string, settableVariable *string) error {
	token, exists := os.LookupEnv(key)
	if !exists {
		return e.Wrap("TELEGRAM_BOT_TOKEN environment variable not set", nil)
	}
	*settableVariable = token
	return nil
}
