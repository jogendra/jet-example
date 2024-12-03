package salesforce

type Config struct {
	AuthURL      string `env:"SALESFORCE_AUTH_URL,notEmpty"`
	ClientID     string `env:"SALESFORCE_CLIENT_ID,notEmpty"`
	ClientSecret string `env:"SALESFORCE_CLIENT_SECRET,notEmpty"`
}
