package web

type DisabledRequest struct {
	IsDisabled bool   `json:"is_disabled"`
	Balance    int    `json:"balance"`
	CustomerId string `json:"customer_id"`
	Id         int    `json:"id"`
}
