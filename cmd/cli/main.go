package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/patrickmn/go-cache"

	"jet-example/internal/config"
	"jet-example/internal/fetcher/salesforce"
	"jet-example/internal/scheduler"
	"jet-example/internal/uploader/s3"
	pkgS3 "jet-example/pkg/s3_client"
)

func main() {
	// load configuration - env variables
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// parent context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize dependencies
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	cacheClient := cache.New(
		cfg.CacheConfig.DefaultExpirationTime,
		cfg.CacheConfig.CleanupInterval,
	)
	sfClient := salesforce.NewSalesforceClient(
		cfg.Salesforce,
		httpClient,
		cacheClient,
	)
	awsCfg, _ := awsconfig.LoadDefaultConfig(ctx)
	s3Client := pkgS3.NewS3Client(awsCfg, cfg.S3ClientConfig)
	s3Uploader := s3.NewS3Uploader(
		cfg.S3.PathPrefix,
		cfg.S3.Bucket,
		s3Client,
	)

	// create and start the scheduler
	s := scheduler.NewScheduler(sfClient, s3Uploader)
	go func() {
		if err := s.Start(ctx); err != nil {
			log.Fatalf("failed to start scheduler: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	s.Stop()
}
