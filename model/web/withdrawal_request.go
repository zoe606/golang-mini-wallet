package web

import "golang-mini-wallet/model/domain"

type WithdrawalRequest struct {
	Amount      int    `json:"amount"`
	ReferenceId string `validate:"required" json:"reference_id"`
	Token       domain.Token
	Wallet      domain.Wallet
}
