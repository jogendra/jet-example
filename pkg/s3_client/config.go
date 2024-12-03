package s3_client

type ClientConf struct {
	CustomEndpoint string `env:"S3_CUSTOM_ENDPOINT" envDefault:""`
	UsePathStyle   bool   `env:"S3_USE_PATH_STYLE" envDefault:"false"`
}
