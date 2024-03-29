package entity

type Transaction struct {
	SenderId                   string
	RecipientBankAccountNumber string `json:"recipientBankAccountNumber"`
	RecipientBankName          string `json:"recipientBankName"`
	FromCurrency               string `json:"fromCurrency"`
	Balances                   int    `json:"balances"`
}
