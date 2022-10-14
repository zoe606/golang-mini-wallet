package web

type EnableResponse struct {
	Wallet ObjEnabled
}

type ObjEnabled struct {
	Id         int         `json:"id"`       //id wallet
	CustomerId string      `json:"owned_by"` //customer id
	Status     interface{} `json:"status"`
	Balance    int         `json:"balance"`
	ActiveAt   interface{} `json:"enabled_at"`
	//Token      string      `json:"token"`
}
