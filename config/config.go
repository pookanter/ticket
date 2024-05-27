package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Services struct {
		Gateway struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
			URL  string `mapstructure:"url"`
		} `mapstructure:"gateway"`
		Authen struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
			URL  string `mapstructure:"url"`
		} `mapstructure:"authen"`
		Ticket struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
			URL  string `mapstructure:"url"`
		} `mapstructure:"ticket"`
		Database struct {
			Host     string `mapstructure:"host"`
			Dbname   string `mapstructure:"dbname"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
		}
	} `mapstructure:"services"`
	PrivateKey         string `mapstructure:"private_key"`
	PublicKey          string `mapstructure:"public_key"`
	AccessTokenExpire  int    `mapstructure:"access_token_expire"`
	RefreshTokenExpire int    `mapstructure:"refresh_token_expire"`
}

func ReadConfig() (Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("\nunable to decode into struct: %s", err)
	}

	return config, nil
}

func ReadPrivateKey(config Config) ([]byte, error) {
	b, err := os.ReadFile(config.PrivateKey)
	if err != nil {
		return []byte{}, fmt.Errorf("\nunable to decode into struct: %s", err)
	}

	return b, nil
}

func ReadPublicKey(config Config) ([]byte, error) {
	b, err := os.ReadFile(config.PublicKey)
	if err != nil {
		return []byte{}, fmt.Errorf("\nunable to decode into struct: %s", err)
	}

	return b, nil
}
