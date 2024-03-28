package entity

type (
	Source struct {
		BankAccountNumber string `json:"bankAccountNumber"`
		BankName          string `json:"bankName"`
	}
	History struct {
		TransactionId    string  `json:"transactionId"`
		UserId           string  `json:"userId"`
		Balance          float64 `json:"balance"`
		Currency         string  `json:"currency"`
		TransferProofImg string  `json:"transferProofImg"`
		CreatedAt        int64   `json:"createdAt"`
		Source           Source  `json:"source"`
	}
)
