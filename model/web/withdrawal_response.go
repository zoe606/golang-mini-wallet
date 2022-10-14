package web

type WithdrawalResponse struct {
	Withdrawal ObjWithdrawal `json:"withdrawal"`
}

type ObjWithdrawal struct {
	Id          int    `json:"id"` //id wallet
	CreatedBy   any    `json:"withdrawn_by"`
	CreatedAt   any    `json:"withdrawn_at"`
	Status      string `json:"status"`
	Amount      int    `json:"amount"`
	ReferenceId string `json:"reference_id"`
}
