package config

import "github.com/spf13/viper"

type Config struct {
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`
	DBUrl   string `mapstructure:"DB_URL"`
}

func (cfg *Config) IsDevelopment() bool {
	return cfg.AppEnv == "development"
}

func Load(config *Config) error {
	viper.BindEnv("APP_ENV")
	viper.BindEnv("APP_PORT")
	viper.BindEnv("DB_URL")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)

	return viper.Unmarshal(config)
}
