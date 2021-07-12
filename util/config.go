package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	GcpProject string `mapstructure:"GCPPROJECT"`
	HttpPort   string `mapstructure:"HTTPPORT"`
	UseDB      string `mapstructure:"USEDB"`
	DBInstance string `mapstructure:"DBINSTANCE"`
	DBRegion   string `mapstructure:"DBREGION"`
	DBHost     string `mapstructure:"DBHOST"`
	DBUser     string `mapstructure:"DBUSER"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBDriver   string `mapstructure:"DBDRIVER"`
	DBName     string `mapstructure:"DBNAME"`
	DBPort     string `mapstructure:"DBPORT"`
	DBSslmode  string `mapstructure:"DBSSLMODE"`
	DBTimezone string `mapstructure:"DBTIMEZONE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("videoapp")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
