package v1

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"golang-mini-wallet/helpers"
	"golang-mini-wallet/model/domain"
	"golang-mini-wallet/model/web"
	"golang-mini-wallet/repositories"
)

type WallerServiceImpl struct {
	WalletRepository repositories.WalletRepository
	DB               *sql.DB
	Validate         *validator.Validate
}

func NewWallerServiceImpl(walletRepository repositories.WalletRepository, DB *sql.DB, validate *validator.Validate) *WallerServiceImpl {
	return &WallerServiceImpl{WalletRepository: walletRepository, DB: DB, Validate: validate}
}

func (w *WallerServiceImpl) Init(ctx context.Context, request web.InitRequest) web.InitResponse {
	err := w.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx, err := w.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	tmpToken := helpers.RandomString(25)
	wallet := domain.Wallet{Token: tmpToken, CustomerId: request.CustomerXid}

	init := w.WalletRepository.Init(ctx, tx, wallet)
	return helpers.ToInitResponse(init)
}

func (w *WallerServiceImpl) Enable(ctx context.Context, request web.TokenRequest) web.EnableResponse {
	err := w.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx, err := w.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	wallet := domain.Wallet{Token: request.Token}
	enableWallet, err := w.WalletRepository.Enable(ctx, tx, wallet)

	return helpers.ToEnabledResponse(enableWallet)
}

func (w WallerServiceImpl) Get(ctx context.Context, request web.TokenRequest) web.EnableResponse {
	err := w.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx, err := w.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	wallet := domain.Wallet{Token: request.Token}
	getWallet, err := w.WalletRepository.Get(ctx, tx, wallet)

	return helpers.ToEnabledResponse(getWallet)
}

func (w *WallerServiceImpl) Deposit(ctx context.Context, request web.DepositRequest) web.DepositResponse {
	err := w.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx, err := w.DB.Begin()
	helpers.PanicIfError(err)
	deposit := w.WalletRepository.Deposit(ctx, tx, request)
	//defer helpers.CommitOrRollback(tx)

	//todo add to wallet with goroutine
	//error context
	//currBalance := request.Wallet.Balance
	//newBalance := currBalance + request.Amount
	//updateBalance := domain.Wallet{
	//	Balance: newBalance,
	//	Id:      request.Wallet.Id,
	//}
	//
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	//tx, err := w.DB.Begin()
	//	//helpers.PanicIfError(err)
	//	w.WalletRepository.UpdateBalance(ctx, tx, updateBalance)
	//	helpers.CommitOrRollback(tx)
	//}()

	return helpers.ToDepositResponse(deposit)
}

func (w *WallerServiceImpl) Withdrawal(ctx context.Context, request web.WithdrawalRequest) web.WithdrawalResponse {
	err := w.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx, err := w.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	deposit := w.WalletRepository.Withdrawal(ctx, tx, request)

	//todo reduce to wallet with goroutine

	return helpers.ToWithdrawalResponse(deposit)
}

func (w *WallerServiceImpl) Disabled(ctx context.Context, request web.DisabledRequest) web.DisabledResponse {
	err := w.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx, err := w.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	disabled := w.WalletRepository.Disabled(ctx, tx, request)

	return helpers.ToDisabledResponse(disabled)
}
