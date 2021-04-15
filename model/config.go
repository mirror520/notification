package model

var Config = struct {
	Every8D struct {
		BaseURL     string
		CallbackURL string
	}

	Mitake struct {
		BaseURL     string
		CallbackURL string
	}

	Providers []SMSProviderProfile
}{}
