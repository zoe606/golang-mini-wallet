package helpers

import (
	"golang-mini-wallet/model/domain"
	"golang-mini-wallet/model/web"
)

func ToInitResponse(token domain.Token) web.InitResponse {
	return web.InitResponse{
		Token: token.Token,
	}
}

func ToEnabledResponse(wallet domain.Wallet) web.EnableResponse {
	var status string

	if wallet.IsActive == true {
		status = "enabled"
	} else {
		status = "disabled"
	}

	return web.EnableResponse{
		Wallet: web.ObjEnabled{
			Id:         wallet.Id,
			CustomerId: wallet.CustomerId,
			Status:     status,
			Balance:    wallet.Balance,
			ActiveAt:   wallet.ActiveAt,
		},
	}
}

func ToDepositResponse(wallet domain.Wallet) web.DepositResponse {
	return web.DepositResponse{
		Deposit: web.ObjDeposit{
			Id:          wallet.Id,
			CreatedBy:   wallet.CreatedBy,
			CreatedAt:   wallet.CreatedAt,
			Status:      "success",
			Amount:      wallet.Amount,
			ReferenceId: wallet.ReferenceId,
		},
	}
}

func ToWithdrawalResponse(wallet domain.Wallet) web.WithdrawalResponse {
	return web.WithdrawalResponse{
		Withdrawal: web.ObjWithdrawal{
			Id:          wallet.Id,
			CreatedBy:   wallet.CreatedBy,
			CreatedAt:   wallet.CreatedAt,
			Status:      "success",
			Amount:      wallet.Amount,
			ReferenceId: wallet.ReferenceId,
		},
	}
}

func ToDisabledResponse(wallet domain.Wallet) web.DisabledResponse {
	return web.DisabledResponse{
		Wallet: web.ObjDisabled{
			Id:         wallet.Id,
			CustomerId: wallet.CustomerId,
			Status:     "disabled",
			Balance:    wallet.Balance,
			DisabledAt: wallet.DisabledAt,
		},
	}
}
