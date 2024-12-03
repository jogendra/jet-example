package salesforce

import "jet-example/internal/domain"

type ContentAssetsResponse struct {
	Count    int                   `json:"count"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
	Items    []domain.ContentBlock `json:"items"`
}
