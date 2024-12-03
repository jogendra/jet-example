package scheduler

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"

	"jet-example/internal/domain"
)

type Scheduler struct {
	fetcher  domain.Fetcher
	uploader domain.Uploader
	cron     *cron.Cron
}

func NewScheduler(
	fetcher domain.Fetcher,
	uploader domain.Uploader,
) *Scheduler {
	return &Scheduler{
		fetcher:  fetcher,
		uploader: uploader,
		cron:     cron.New(),
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	_, err := s.cron.AddFunc("@every 24h", func() {
		if err := s.fetchAndSyncContentBlocks(ctx); err != nil {
			log.Printf("syncing content blocks failed: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	log.Println("scheduler started")
	return nil
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("scheduler stopped")
}

func (s *Scheduler) fetchAndSyncContentBlocks(ctx context.Context) error {
	// fetch content blocks
	contentBlocks, err := s.fetcher.FetchContentBlocks(ctx, domain.ContentBlocksRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch content blocks: %w", err)
	}

	// upload to provided uploader
	if err := s.uploader.UploadContentBlocks(ctx, contentBlocks); err != nil {
		return fmt.Errorf("failed to upload content blocks: %w", err)
	}

	return nil
}
