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

var config Config
var privateKey string
var publicKey string

func Initialize() {
	fmt.Println("Initializing config...")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("\nfatal error config file: %s", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("\nunable to decode into struct: %s", err))
	}

	b, err := os.ReadFile(config.PrivateKey)
	if err != nil {
		panic(err)
	}

	privateKey = string(b)

	fmt.Println("Initializing config completed!")
}

func GetConfig() Config {
	return config
}

func GetPrivateKey() string {
	return privateKey
}
