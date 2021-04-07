package model

type SMS struct {
	Phone   string `json:"phone" binding:"required"`
	Message string `json:"message" binding:"required"`
	Comment string `json:"comment"`
}

type SMSResult struct {
	ID     string
	Credit int
}
