package domain

import "context"

// Uploader uploads content blocks to implemented uploader e.g. local, s3_client
type Uploader interface {
	UploadContentBlocks(ctx context.Context, contentBlocks []ContentBlock) error
}

// Fetcher fetches content blocks to implemented e.g. salesforce (as of now)
type Fetcher interface {
	FetchContentBlocks(ctx context.Context, request ContentBlocksRequest) ([]ContentBlock, error)
}
