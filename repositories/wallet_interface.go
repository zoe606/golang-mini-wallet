package repositories

import (
	"context"
	"database/sql"
	"golang-mini-wallet/model/domain"
	"golang-mini-wallet/model/web"
)

type WalletRepository interface {
	Init(ctx context.Context, tx *sql.Tx, data domain.Wallet) domain.Token
	Enable(ctx context.Context, tx *sql.Tx, wallet domain.Wallet) (domain.Wallet, error)
	Get(ctx context.Context, tx *sql.Tx, wallet domain.Wallet) (domain.Wallet, error)
	Deposit(ctx context.Context, tx *sql.Tx, data web.DepositRequest) domain.Wallet
	Withdrawal(ctx context.Context, tx *sql.Tx, data web.WithdrawalRequest) domain.Wallet
	Disabled(ctx context.Context, tx *sql.Tx, data web.DisabledRequest) domain.Wallet
}
