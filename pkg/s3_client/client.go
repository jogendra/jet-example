package s3_client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(awsCfg aws.Config, cfg ClientConf) *s3.Client {
	client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		setS3ClientOptions(cfg, options)
	})
	return client
}

func setS3ClientOptions(cfg ClientConf, options *s3.Options) *s3.Options {
	if cfg.CustomEndpoint != "" {
		options.BaseEndpoint = aws.String(cfg.CustomEndpoint)
	}
	if cfg.UsePathStyle {
		options.UsePathStyle = true
	}
	return options
}
