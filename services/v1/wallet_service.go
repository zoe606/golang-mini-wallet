package v1

import (
	"context"
	"golang-mini-wallet/model/web"
)

type WalletService interface {
	// Init create wallet and get token
	Init(ctx context.Context, request web.InitRequest) web.InitResponse

	//Enable set active wallet on db
	//If the wallet is already enabled, this endpoint would fail
	//This endpoint should be usable again only if the wallet is disabled
	//Before enabling the wallet, the customer cannot view, add, or use its virtual money.
	Enable(ctx context.Context, request web.TokenRequest) web.EnableResponse

	//Get Wallet after adding or using virtual money, it is not expected to have the balance immediately updated.
	// The maximum delay for updating the balance is 5 seconds.
	Get(ctx context.Context, request web.TokenRequest) web.EnableResponse

	//Deposit customer can add virtual money to the wallet balance as a deposit once the wallet is enabled.
	//Reference ID passed must be unique for every deposit.
	Deposit(ctx context.Context, request web.DepositRequest) web.DepositResponse

	//Withdrawal a customer can use the virtual money from the wallet balance as a withdrawal once the wallet is enabled
	//The amount being used must not be more than the current balance
	// Reference ID passed must be unique for every withdrawal
	Withdrawal(ctx context.Context, request web.WithdrawalRequest) web.WithdrawalResponse

	// Disabled customer's wallet can be disabled as determined by the service.
	//Once disabled, the customer cannot view, add, or use its virtual money.
	Disabled(ctx context.Context, request web.DisabledRequest) web.DisabledResponse
}
