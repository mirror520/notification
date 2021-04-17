package model

var Config = struct {
	InfluxDB struct {
		URL    string
		Org    string
		Bucket string
		Token  string
	}

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
