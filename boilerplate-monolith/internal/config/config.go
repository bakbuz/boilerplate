package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

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

func LoadConfig(path string) (*Config, error) {
	var config Config

	// 1. Önce .env dosyasını okumayı dener (local geliştirme için)
	// 2. Ardından ortam değişkenlerini (ENV vars) okur ve struct'ı doldurur.
	// Canlı ortamda .env dosyası olmasa bile sistem ENV'lerini okur.
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		// .env yoksa sadece ENV var'lardan okumayı dene (Production senaryosu)
		err = cleanenv.ReadEnv(&config)
		if err != nil {
			return nil, errors.WithMessagef(err, "Config yüklenemedi: %s", path)
		}
	}

	return &config, nil
}
