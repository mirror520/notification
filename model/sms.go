package model

type SMSProviderType int

const (
	Every8D SMSProviderType = iota
	Mitake
)

var ProviderType = map[SMSProviderType]string{
	Every8D: "every8d",
	Mitake:  "mitake",
}

type SMSRole int

const (
	Backup SMSRole = iota
	Master
)

var ProviderRole = map[SMSRole]string{
	Backup: "backup",
	Master: "master",
}

type SMSAccount struct {
	Username string
	Password string
}

type SMSProviderProfile struct {
	ID      string
	Type    SMSProviderType
	Role    SMSRole
	Account SMSAccount
}

func (p *SMSProviderProfile) ProviderType() string {
	return ProviderType[p.Type]
}

func (p *SMSProviderProfile) ProviderRole() string {
	return ProviderRole[p.Role]
}

type SMS struct {
	Phone   string `json:"phone" binding:"required"`
	Message string `json:"message" binding:"required"`
	Comment string `json:"comment"`
}

type SMSResult struct {
	ID     string
	Credit int
}
