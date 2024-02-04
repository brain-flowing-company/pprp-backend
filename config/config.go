package config

import "github.com/spf13/viper"

type Config struct {
	AppEnv         string   `mapstructure:"APP_ENV"`
	AppPort        string   `mapstructure:"APP_PORT"`
	DBUrl          string   `mapstructure:"DB_URL"`
	GoogleClientId string   `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleSecret   string   `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirect string   `mapstructure:"GOOGLE_REDIRECT"`
	GoogleScopes   []string `mapstructure:"GOOGLE_SCOPE"`
	JWTSecret      string   `mapstructure:"JWT_SECRET"`
	SessionExpire  int      `mapstructure:"SESSION_EXPIRES"`
}

func (cfg *Config) IsDevelopment() bool {
	return cfg.AppEnv == "development"
}

func Load(config *Config) error {
	_ = viper.BindEnv("APP_ENV")
	_ = viper.BindEnv("APP_PORT")
	_ = viper.BindEnv("DB_URL")
	_ = viper.BindEnv("GOOGLE_CLIENT_ID")
	_ = viper.BindEnv("GOOGLE_CLIENT_SECRET")
	_ = viper.BindEnv("GOOGLE_REDIRECT")
	_ = viper.BindEnv("GOOGLE_SCOPE")
	_ = viper.BindEnv("JWT_SECRET")
	_ = viper.BindEnv("SESSION_EXPIRES")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)

	return viper.Unmarshal(config)
}
