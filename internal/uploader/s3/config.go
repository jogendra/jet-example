package s3

type Config struct {
	Bucket     string `env:"S3_BUCKET,notEmpty"`
	PathPrefix string `env:"S3_PATH_PREFIX,notEmpty"`
}
