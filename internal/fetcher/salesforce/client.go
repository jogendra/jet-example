package salesforce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"jet-example/internal/domain"
)

// Cache keys
const (
	cacheKeyAccessTokenKey     = "accessToken"
	cacheKeyRestInstanceURLKey = "restInstanceURL"
)

type client struct {
	config     Config
	httpClient *http.Client
	cache      *cache.Cache
}

func NewSalesforceClient(
	config Config,
	httpClient *http.Client,
	cache *cache.Cache,
) domain.Fetcher {
	return &client{
		config:     config,
		httpClient: httpClient,
		cache:      cache,
	}
}

// FetchContentBlocks fetches content blocks from Salesforce concurrently with error handling
func (c *client) FetchContentBlocks(
	ctx context.Context,
	request domain.ContentBlocksRequest,
) ([]domain.ContentBlock, error) {
	tokenResponse, err := c.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch access token: %w", err)
	}

	// fetch the first page to get the total count
	currentPage := 1
	request.Page.Page = currentPage
	firstPageResponse, err := fetchSingleAssetPage(
		ctx,
		tokenResponse.RestInstanceURL,
		tokenResponse.AccessToken,
		request,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch first page of assets: %w", err)
	}

	// calculate the total number of pages
	totalPages := int(math.Ceil(float64(firstPageResponse.Count) / float64(firstPageResponse.PageSize)))

	// create a channel to receive content blocks from goroutines
	contentBlockChan := make(chan []domain.ContentBlock, totalPages)

	var wg sync.WaitGroup
	wg.Add(totalPages)

	// launch a goroutine for each page
	for page := 1; page <= totalPages; page++ {
		go func(page int) {
			defer wg.Done()
			request.Page.Page = page
			response, err := fetchSingleAssetPage(
				ctx,
				tokenResponse.RestInstanceURL,
				tokenResponse.AccessToken,
				request,
			)
			if err != nil {
				// I left tech debt here
				// not sure if I should fail whole request if problem fetch one page
				// or continue with other pages
				contentBlockChan <- []domain.ContentBlock{} // Send an empty slice
				fmt.Printf("Error fetching page %d: %v\n", page, err)
				return
			}
			contentBlockChan <- response.Items
		}(page)
	}

	// wait for all goroutines to finish
	wg.Wait()
	close(contentBlockChan)

	// collect results from the channel
	var allContentBlocks []domain.ContentBlock
	for blocks := range contentBlockChan {
		allContentBlocks = append(allContentBlocks, blocks...)
	}

	return allContentBlocks, nil
}

// FetchAccessToken fetches an access token
func (c *client) fetchAccessToken(
	ctx context.Context,
) (response TokenResponse, err error) {
	// check cache for token and instance URL
	if cachedToken, found := c.getCachedToken(); found {
		return cachedToken, nil
	}

	authURL := c.config.AuthURL + "/v2/token"

	request := TokenRequest{
		GrantType:    "client_credentials", // since this is server-to-server integration (according to docs)
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		authURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		if httpResponse.StatusCode == http.StatusUnauthorized {
			return TokenResponse{}, fmt.Errorf("unauthorized")
		}

		return TokenResponse{},
			fmt.Errorf("failed to fetch access token, status code: %d", httpResponse.StatusCode)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&tokenResponse); err != nil {
		return TokenResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	// store the new token and instance URL in the cache
	expiration := time.Duration(tokenResponse.ExpiresIn) * time.Second
	// API documentations recommend that we refresh our token two minutes before its lifetime ends.
	safeExpiration := expiration - 2*time.Minute
	if safeExpiration > 0 { // Ensure expiration is not negative
		c.cache.Set(cacheKeyAccessTokenKey, tokenResponse.AccessToken, expiration)
		c.cache.Set(cacheKeyRestInstanceURLKey, tokenResponse.RestInstanceURL, expiration)
	}

	return tokenResponse, nil
}

// getCachedToken retrieves the token and instance URL from the cache.
func (c *client) getCachedToken() (TokenResponse, bool) {
	accessToken, foundToken := c.cache.Get(cacheKeyAccessTokenKey)
	instanceURL, foundInstance := c.cache.Get(cacheKeyRestInstanceURLKey)

	if foundToken && foundInstance {
		return TokenResponse{
			AccessToken:     accessToken.(string),
			RestInstanceURL: instanceURL.(string),
		}, true
	}

	return TokenResponse{}, false
}

// fetchSinglePage fetches a single page of assets based on the query
func fetchSingleAssetPage(
	ctx context.Context,
	instanceURL,
	accessToken string,
	request domain.ContentBlocksRequest,
) (ContentAssetsResponse, error) {
	url := instanceURL + "/asset/v1/content/assets/query"

	requestBody, err := json.Marshal(request)
	if err != nil {
		return ContentAssetsResponse{}, fmt.Errorf("failed to marshal query request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return ContentAssetsResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+accessToken)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return ContentAssetsResponse{}, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return ContentAssetsResponse{},
			fmt.Errorf("failed to fetch assets, status code: %d", httpResponse.StatusCode)
	}

	var response ContentAssetsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return ContentAssetsResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return response, nil
}
