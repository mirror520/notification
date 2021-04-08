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

var Role = map[SMSRole]string{
	Backup: "backup",
	Master: "master",
}

type SMSAccount struct {
	Username string
	Password string
}

type SMSProvider struct {
	Name    string
	Type    SMSProviderType
	Role    SMSRole
	Account SMSAccount
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
