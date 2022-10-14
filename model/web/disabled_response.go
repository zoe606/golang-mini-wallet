package web

type DisabledResponse struct {
	Wallet ObjDisabled `json:"wallet"`
}

type ObjDisabled struct {
	Id         int         `json:"id"`       //id wallet
	CustomerId string      `json:"owned_by"` //customer id
	Status     string      `json:"status"`
	Balance    int         `json:"balance"`
	DisabledAt interface{} `json:"disabled_at"`
}
