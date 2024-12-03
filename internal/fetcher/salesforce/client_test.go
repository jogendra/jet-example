package salesforce

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"

	"jet-example/internal/domain"
)

func TestSalesforceClient_FetchAccessToken(t *testing.T) {
	tests := []struct {
		name          string
		setupCache    func(*cache.Cache)
		setupHTTP     func() *httptest.Server
		expectedToken TokenResponse
		wantErr       require.ErrorAssertionFunc
	}{
		{
			name: "Token retrieved from cache",
			setupCache: func(c *cache.Cache) {
				c.Set(cacheKeyAccessTokenKey, "cachedToken", time.Minute)
				c.Set(cacheKeyRestInstanceURLKey, "cachedInstanceURL", time.Minute)
			},
			setupHTTP: func() *httptest.Server { return nil }, // No need for mock server when using cache
			expectedToken: TokenResponse{
				AccessToken:     "cachedToken",
				RestInstanceURL: "cachedInstanceURL",
			},
			wantErr: require.NoError,
		},
		{
			name:       "Successful token retrieval from API",
			setupCache: func(c *cache.Cache) {},
			setupHTTP: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"access_token": "newToken",
						"rest_instance_url": "newInstanceURL",
						"expires_in": 3600
					}`))
				}))
			},
			expectedToken: TokenResponse{
				AccessToken:     "newToken",
				RestInstanceURL: "newInstanceURL",
				ExpiresIn:       3600,
			},
			wantErr: require.NoError,
		},
		{
			name:       "Unauthorized error",
			setupCache: func(c *cache.Cache) {},
			setupHTTP: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
				}))
			},
			expectedToken: TokenResponse{},
			wantErr: func(t require.TestingT, err error, _ ...interface{}) {
				require.ErrorContains(t, err, "unauthorized")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheInstance := cache.New(5*time.Minute, 10*time.Minute)
			tt.setupCache(cacheInstance)
			authURL := "test.com"
			var server *httptest.Server
			if tt.setupHTTP() != nil {
				server = tt.setupHTTP()
				authURL = server.URL
				defer server.Close()
			}

			client := NewSalesforceClient(
				Config{
					AuthURL:      authURL,
					ClientID:     "some-client-id",
					ClientSecret: "some-client-secret",
				},
				http.DefaultClient,
				cacheInstance,
			)

			token, err := client.FetchAccessToken(context.Background())

			tt.wantErr(t, err)
			require.Equal(t, tt.expectedToken, token)

			cacheInstance.Flush()
		})
	}
}

func TestSalesforceClient_FetchContentBlocks(t *testing.T) {
	tests := []struct {
		name              string
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		request           domain.ContentBlocksRequest
		want              []domain.ContentBlock
		wantErr           require.ErrorAssertionFunc
	}{
		{
			name: "Success - Single Page",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				response := domain.ContentAssetsResponse{
					Count:    2,
					Page:     1,
					PageSize: 10,
					Items: []domain.ContentBlock{
						{Content: "Block 1"},
						{Content: "Block 2"},
					},
				}
				json.NewEncoder(w).Encode(response)
			},
			request: domain.ContentBlocksRequest{
				Page: struct {
					Page     int `json:"page"`
					PageSize int `json:"pageSize"`
				}{Page: 1, PageSize: 10},
			},
			want: []domain.ContentBlock{
				{Content: "Block 1"},
				{Content: "Block 2"},
			},
			wantErr: require.NoError,
		},
		{
			name: "Success - Multiple Pages",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				// Mock responses for 2 pages
				var request domain.ContentBlocksRequest
				json.NewDecoder(r.Body).Decode(&request)
				page := request.Page.Page
				if page == 1 {
					response := domain.ContentAssetsResponse{
						Count:    3,
						Page:     1,
						PageSize: 2,
						Items: []domain.ContentBlock{
							{Content: "Block 1"},
							{Content: "Block 2"},
						},
					}
					json.NewEncoder(w).Encode(response)
				} else if page == 2 {
					response := domain.ContentAssetsResponse{
						Count:    3,
						Page:     2,
						PageSize: 2,
						Items: []domain.ContentBlock{
							{Content: "Block 3"},
						},
					}
					json.NewEncoder(w).Encode(response)
				}
			},
			request: domain.ContentBlocksRequest{
				Page: struct {
					Page     int `json:"page"`
					PageSize int `json:"pageSize"`
				}{Page: 1, PageSize: 2},
			},
			want: []domain.ContentBlock{
				{Content: "Block 1"},
				{Content: "Block 2"},
				{Content: "Block 3"},
			},
			wantErr: require.NoError,
		},
		{
			name: "Error - Fetching First Page",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				// Mock an error response for the first page
				w.WriteHeader(http.StatusInternalServerError)
			},
			request: domain.ContentBlocksRequest{},
			want:    nil,
			wantErr: require.Error,
		},
		{
			name: "Fetching Subsequent Page - one of page fail to fetch",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				// Mock a successful response for the first page, but an error for the second
				var request domain.ContentBlocksRequest
				json.NewDecoder(r.Body).Decode(&request)
				page := request.Page.Page
				if page == 1 {
					response := domain.ContentAssetsResponse{
						Count:    3,
						Page:     1,
						PageSize: 2,
						Items: []domain.ContentBlock{
							{Content: "Block 1"},
							{Content: "Block 2"},
						},
					}
					json.NewEncoder(w).Encode(response)
				} else if page == 2 {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			request: domain.ContentBlocksRequest{
				Page: struct {
					Page     int `json:"page"`
					PageSize int `json:"pageSize"`
				}{Page: 1, PageSize: 2},
			},
			want:    []domain.ContentBlock{{Content: "Block 1"}, {Content: "Block 2"}}, // Partial result
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(tt.mockServerHandler))
			defer mockServer.Close()

			c := &client{
				config: Config{
					AuthURL:      mockServer.URL, // Use the mock server URL
					ClientID:     "testClientID",
					ClientSecret: "testClientSecret",
				},
				httpClient: http.DefaultClient,
				cache:      cache.New(5*time.Minute, 10*time.Minute),
			}

			// Set up the access token in the cache (to avoid fetching it)
			c.cache.Set(cacheKeyAccessTokenKey, "testAccessToken", cache.DefaultExpiration)
			c.cache.Set(cacheKeyRestInstanceURLKey, mockServer.URL, cache.DefaultExpiration)

			got, err := c.FetchContentBlocks(context.Background(), tt.request)

			// check error match
			tt.wantErr(t, err)

			// check result - not in order hence ElementsMatch
			require.ElementsMatch(t, got, tt.want)

			c.cache.Flush()
		})
	}
}
