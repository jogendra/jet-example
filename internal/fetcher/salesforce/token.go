package salesforce

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// TokenType       string `json:"token_type"` -> Always "Bearer"
	RestInstanceURL string `json:"rest_instance_url"`
}
