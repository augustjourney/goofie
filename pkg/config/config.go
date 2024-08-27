package config

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseDSN           string `mapstructure:"DATABASE_DSN"`
	JwtTokenExp           time.Duration
	JwtTokenSecretKey     string `mapstructure:"JWT_TOKEN_SECRET_KEY"`
	JwtSigningMethod      jwt.SigningMethod
	TempFolder            string
	S3Bucket              string
	S3Provider            string
	S3Region              string
	S3Endpoint            string
	S3AccessKeyId         string `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey     string `mapstructure:"S3_SECRET_ACCESS_KEY"`
	CdnURL                string
	DefaultImageExtension string
}

var (
	config *Config
	once   sync.Once
)

func LoadConfig(path string) (*Config, error) {
	var err error
	once.Do(func() {
		config = &Config{}

		viper.AddConfigPath(path)
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			return
		}
		err = viper.Unmarshal(&config)

		config.JwtSigningMethod = jwt.SigningMethodHS256
		config.JwtTokenExp = time.Hour * 720
		config.TempFolder = "temp"
		config.S3Provider = "selectel"
		config.S3Bucket = "photos"
		config.S3Region = "ru-1a"
		config.S3Endpoint = "https://s3.storage.selcloud.ru"
		config.CdnURL = "https://713726.selcdn.ru"
		config.DefaultImageExtension = ".webp"
	})
	return config, err
}

func Get() *Config {
	return config
}
