package entity

type (
	Balance struct {
		Id        string  `json:"id"`
		UserId    string  `json:"userId"`
		Balance   float64 `json:"balance"`
		Currency  string  `json:"currency"`
		CreatedAt string  `json:"createdAt"`
	}

	BalanceResponse struct {
		Balance  float64 `json:"balance"`
		Currency string  `json:"currency"`
	}
)
