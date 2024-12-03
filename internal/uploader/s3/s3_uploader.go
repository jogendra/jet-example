package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"jet-example/internal/domain"
)

type s3Uploader struct {
	s3Client     *s3.Client
	s3Bucket     string
	s3PathPrefix string
}

func NewS3Uploader(
	bucket string,
	pathPrefix string,
	client *s3.Client,
) domain.Uploader {
	return &s3Uploader{
		s3Client:     client,
		s3Bucket:     strings.Trim(bucket, "/"),
		s3PathPrefix: strings.Trim(pathPrefix, "/"),
	}
}

func (u *s3Uploader) UploadContentBlocks(
	ctx context.Context,
	contentBlocks []domain.ContentBlock,
) error {
	jsonData, err := json.Marshal(contentBlocks)
	if err != nil {
		return err
	}

	objectKey := fmt.Sprintf(
		"%s/%s.json",
		time.Now().Format("2006-01-02"),
		"content-block",
	)
	err = u.upload(ctx, objectKey, jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *s3Uploader) upload(ctx context.Context, key string, body []byte) error {
	_, err := u.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.s3Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	return err
}
