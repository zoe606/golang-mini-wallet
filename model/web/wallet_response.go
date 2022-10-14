package web

type WalletResponse struct {
	Id         int    `json:"id"`       //id wallet
	CustomerId string `json:"owned_by"` //customer id
	Status     string `json:"status"`
	Balance    int    `json:"balance"`
	ActiveAt   any    `json:"enabled_at"`
}
