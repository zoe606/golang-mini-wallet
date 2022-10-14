package web

type DepositResponse struct {
	Deposit ObjDeposit `json:"deposit"`
}

type ObjDeposit struct {
	Id          int    `json:"id"` //id wallet
	CreatedBy   string `json:"deposited_by"`
	CreatedAt   any    `json:"deposited_at"`
	Status      string `json:"status"`
	Amount      int    `json:"amount"`
	ReferenceId string `json:"reference_id"`
}
