package model

var Config = struct {
	Every8D struct {
		BaseURL string
	}

	Mitake struct {
		BaseURL string
	}

	Providers []SMSProvider
}{}
