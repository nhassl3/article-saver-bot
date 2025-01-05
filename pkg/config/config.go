package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nhassl3/article-saver-bot/pkg/e"
	"gopkg.in/yaml.v2"
	"os"
)

var ErrToGetEnvVariable = errors.New("error to ger environment variable")

type Config struct {
	Server `yaml:"server"`
	Token  string
}

type Server struct {
	Protocol string `yaml:"protocol" env-default:"https"`
	Host     string `yaml:"host" env-default:"api.telegram.org"`
}

func (c *Config) MustLoad() (err error) {
	defer func() { err = e.WrapIfErr("error to load config", err) }()
	if err = godotenv.Load(); err != nil {
		return err
	}

	var token string
	if err = GetVariable("TELEGRAM_BOT_TOKEN", &token); err != nil {
		return err
	}
	c.Token = token

	var configPath string
	if err = GetVariable("CONFIG_PATH", &configPath); err != nil {
		return err
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	return e.WrapIfErr("err to unmarshal config file", err)
}

func GetVariable(key string, settableVariable *string) error {
	token, exists := os.LookupEnv(key)
	if !exists {
		// create error
		return e.Wrap(fmt.Sprintf("%s environment variable not set", key), ErrToGetEnvVariable)
	}
	*settableVariable = token
	return nil
}
