package env

import (
	"github.com/spf13/viper"
	"savebite/pkg/log"
	"time"
)

type Env struct {
	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	AppURL  string `mapstructure:"APP_URL"`
	AppPort string `mapstructure:"APP_PORT"`

	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`

	JwtSecretKey string        `mapstructure:"JWT_SECRET_KEY"`
	JwtExpTime   time.Duration `mapstructure:"JWT_EXP_TIME"`

	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`

	GeminiAPIKey string `mapstructure:"GEMINI_API_KEY"`
	GeminiModel  string `mapstructure:"GEMINI_MODEL"`
}

var AppEnv = getEnv()

func getEnv() *Env {
	env := &Env{}

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[Env][getEnv] Failed to read config file")
	}

	err = viper.Unmarshal(env)
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[Env][getEnv] Failed to unmarshal viper")
	}

	switch env.AppEnv {
	case "development":
		log.Info(nil, "Application is running on development mode")
	case "production":
		log.Info(nil, "Application is running on production mode")
	}

	return env
}
