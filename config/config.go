package config

import "github.com/spf13/viper"

const envFileName = ".env"

type Config struct {
	AppHost                           string `mapstructure:"APP_HOST"`
	AppPort                           int    `mapstructure:"APP_PORT"`
	AppProtocol                       string `mapstructure:"APP_PROTOCOL"`
	DbHost                            string `mapstructure:"DB_HOST"`
	DbPort                            int    `mapstructure:"DB_PORT"`
	DbUsername                        string `mapstructure:"DB_USERNAME"`
	DbPassword                        string `mapstructure:"DB_PASSWORD"`
	DbName                            string `mapstructure:"DB_NAME"`
	LoggerLevel                       string `mapstructure:"LOGGER_LEVEL"`
	LoggerOutput                      string `mapstructure:"stdout"`
	TokenAccessExpirationMinute       int    `mapstructure:"TOKEN_ACCESS_EXPIRATION_MINUTE"`
	TokenRefreshExpirationMinute      int    `mapstructure:"TOKEN_REFRESH_EXPIRATION_MINUTE"`
	TokenVerificationExpirationMinute int    `mapstructure:"TOKEN_VERIFICATION_EXPIRATION_MINUTE"`
}

func LoadConfig(path string) (Config, error) {
	config := Config{}
	viper.AddConfigPath(path)
	viper.SetConfigFile(envFileName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
