package config

import (
	"github.com/caarlos0/env/v10"

	"jet-example/internal/fetcher/salesforce"
	"jet-example/internal/uploader/s3"
	"jet-example/pkg/s3_client"
)

type AppConfig struct {
	Salesforce     salesforce.Config
	S3             s3.Config
	CacheConfig    CacheConfig
	S3ClientConfig s3_client.ClientConf
}

func LoadAppConfig() (AppConfig, error) {
	appConfig := AppConfig{}

	if err := env.Parse(&appConfig); err != nil {
		return AppConfig{}, err
	}

	return appConfig, nil
}
