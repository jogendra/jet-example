package local

import (
	"context"
	"fmt"

	"jet-example/internal/domain"
)

type localUploader struct {
	directory string
}

// NewLocalUploader implement ContentUploader method and save data to local machine
// It returns an error if the directory cannot be created.
func NewLocalUploader() (domain.Uploader, error) {
	return &localUploader{}, nil
}

// UploadContentBlocks unimplemented - ideally store data to local machine
func (u *localUploader) UploadContentBlocks(
	ctx context.Context,
	contentBlocks []domain.ContentBlock,
) error {
	return fmt.Errorf("unimplemented")
}
