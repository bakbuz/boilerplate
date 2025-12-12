package bootstrap

import "github.com/spf13/viper"

type Config struct {
	DataSources struct {
		Default string
	}
	Supabase struct {
		ProjectId      string
		AnonPublicKey  string
		ServiceRoleKey string
	}
	Jwt struct {
		SecretKey string
	}
}

func LoadConfig(in string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.SetConfigFile(in)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
