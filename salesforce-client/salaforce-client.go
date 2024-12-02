package salesforce_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"jet-example/domain"
)

type SalesforceClient interface {
	GetAccessToken(url string, payload domain.TokenPayload) (response domain.TokenResponse, err error)
	FetchContent(path, request domain.ContentRequest) (domain.ContentResponse, error)
}

type salesforceClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewSalesforceClient(baseURL string, httpClient *http.Client) SalesforceClient {
	return &salesforceClientImpl{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (sf *salesforceClientImpl) GetAccessToken(
	url string,
	payload domain.TokenPayload,
) (response domain.TokenResponse, err error) {
	authURL := sf.baseURL + "/v2/token"

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("error marshalling payload %v", err)
		// log error
		return domain.TokenResponse{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("error creating new request %v", err)
		return domain.TokenResponse{}, nil
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := sf.httpClient.Do(req)

	if err != nil {
		fmt.Printf("error sending request %v", err)
		return domain.TokenResponse{}, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response %v", err)
		return domain.TokenResponse{}, nil
	}
	defer resp.Body.Close()

	var tokenResponse domain.TokenResponse
	err = json.Unmarshal(responseBody, &tokenResponse)

	if err != nil {
		fmt.Printf("error unmarshalling response %v", err)
		return domain.TokenResponse{}, nil
	}

	return tokenResponse, nil
}
