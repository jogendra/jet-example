package domain

type TokenPayload struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	AccountID    string `json:"account_id"`
}

type TokenResponse struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn       int    `json:"expires_in"`
	RestInstanceURL string `json:"rest_instance_url"`
}

type ContentRequest struct {
}

type ContentResponse struct{}
