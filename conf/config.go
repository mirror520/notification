package conf

type Config struct {
	Providers map[string]Provider
}

type ProviderType string

const (
	Every8D ProviderType = "every8d"
	Mitake  ProviderType = "mitake"
)

type Provider struct {
	Type        ProviderType `yaml:"type"`
	Username    string       `yaml:"username"`
	Password    string       `yaml:"password"`
	BaseURL     string       `yaml:"baseUrl"`
	CallbackURL string       `yaml:"callbackUrl"`
}
