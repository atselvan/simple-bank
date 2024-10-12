package util

import "github.com/spf13/viper"

type (
	// Config stores all the configuration of the application.
	Config struct {
		DBDriver      string `mapstructure:"DB_DRIVER"`
		DBHost        string `mapstructure:"DB_HOST"`
		DBPort        string `mapstructure:"DB_PORT"`
		DBUserName    string `mapstructure:"DB_USERNAME"`
		DBPassword    string `mapstructure:"DB_PASSWORD"`
		DBName        string `mapstructure:"DB_NAME"`
		ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	}
)

// LoadConfig reads configuration from a file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
