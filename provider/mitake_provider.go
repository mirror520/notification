package provider

type MitakeProvider struct {
}

func (p *MitakeProvider) SendSMS(phone, message string) {
}

func (p *MitakeProvider) Credit() (int, error) {
	return 0, nil
}
