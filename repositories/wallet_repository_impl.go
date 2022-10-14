package repositories

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"golang-mini-wallet/helpers"
	"golang-mini-wallet/model/domain"
	"golang-mini-wallet/model/web"
	"time"
)

type WalletRepositoryImpl struct{}

func NewWalletRepository() WalletRepository {
	return &WalletRepositoryImpl{}
}

func (w *WalletRepositoryImpl) Init(ctx context.Context, tx *sql.Tx, data domain.Wallet) domain.Token {
	query := "SELECT id, token FROM wallet_m where customer_id = ?"
	rows, err := tx.QueryContext(ctx, query, data.CustomerId)
	helpers.PanicIfError(err)
	defer rows.Close()

	token := domain.Token{}
	wallet := domain.Wallet{}

	if rows.Next() {
		err2 := rows.Scan(&wallet.Id, &wallet.Token)
		helpers.PanicIfError(err2)
		token.Token = wallet.Token
		return token
	} else {
		qry := "INSERT INTO wallet_m(customer_id,token) values (?,?)"
		_, err := tx.ExecContext(ctx, qry, data.CustomerId, data.Token)
		helpers.PanicIfError(err)
		token.Token = data.Token
		return token
	}
}

func (w *WalletRepositoryImpl) Enable(ctx context.Context, tx *sql.Tx, wallet domain.Wallet) (domain.Wallet, error) {
	query := "SELECT id, customer_id, is_active,balance FROM wallet_m where token = ? and is_active = false"
	rows, err := tx.QueryContext(ctx, query, wallet.Token)
	helpers.PanicIfError(err)

	accountWallet := domain.Wallet{}

	if rows.Next() {
		err2 := rows.Scan(&accountWallet.Id, &accountWallet.CustomerId, &accountWallet.IsActive, &accountWallet.Balance)
		helpers.PanicIfError(err2)
	} else {
		return accountWallet, errors.New("Wallet not found")
	}

	rows.Close()

	if accountWallet.Id != 0 {
		query = "UPDATE wallet_m set is_active = true , active_at = NOW() where id = ?"
		res, err := tx.QueryContext(ctx, query, accountWallet.Id)
		helpers.PanicIfError(err)
		res.Close()

		loc, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(loc)
		accountWallet.CreatedAt = now
		accountWallet.IsActive = true
	}

	return accountWallet, nil
}

func (w *WalletRepositoryImpl) Get(ctx context.Context, tx *sql.Tx, wallet domain.Wallet) (domain.Wallet, error) {
	query := "SELECT id, customer_id, is_active, balance , active_at FROM wallet_m where token = ? and is_active = true"
	rows, err := tx.QueryContext(ctx, query, wallet.Token)
	helpers.PanicIfError(err)
	defer rows.Close()
	accountWallet := domain.Wallet{}

	if rows.Next() {
		err2 := rows.Scan(&accountWallet.Id, &accountWallet.CustomerId, &accountWallet.IsActive, &accountWallet.Balance, &accountWallet.ActiveAt)
		helpers.PanicIfError(err2)
		return accountWallet, nil
	} else {
		return accountWallet, errors.New("Wallet not found")
	}

}

func (w *WalletRepositoryImpl) Deposit(ctx context.Context, tx *sql.Tx, data web.DepositRequest) (domain.Wallet, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	query := "SELECT COUNT(id) FROM wallet_t where reference_id = ? "
	rows, err := tx.QueryContext(ctx, query, data.ReferenceId)
	helpers.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		accountWallet := domain.Wallet{}
		return accountWallet, errors.New("Reference id must unique!")
	} else {
		query := "INSERT INTO wallet_t(wallet_id, amount, reference_id, type, created_at, created_by) values (?,?,?,?,?,?)"
		rows, err := tx.QueryContext(ctx, query, data.Wallet.Id, data.Amount, data.ReferenceId, "deposit", now.Format("2006-01-02 03:04:05"), data.Wallet.CustomerId)
		helpers.PanicIfError(err)
		rows.Close()
	}

	accountWallet := domain.Wallet{
		Id:          data.Wallet.Id,
		Amount:      data.Amount,
		CreatedBy:   data.Wallet.CustomerId,
		CreatedAt:   now,
		ReferenceId: data.ReferenceId,
	}

	//todo add to wallet with goroutine
	currBalance := data.Wallet.Balance
	newBalance := currBalance + data.Amount

	go func() {
		time.Sleep(5 * time.Second)
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/golang_mini_wallet")
		defer db.Close()
		res, err := db.Query("UPDATE wallet_m set balance = ? where id = ?", newBalance, data.Wallet.Id)
		defer res.Close()
		//tx, err := w.DB.Begin()
		helpers.PanicIfError(err)
		//UpdateBalance(ctx, tx, updateBalance)
		//query := "UPDATE wallet_m set balance = ? where id = ?"
		//_, err := tx.ExecContext(ctx, query, newBalance, data.Wallet.Id)
		//helpers.PanicIfError(err)

		//helpers.CommitOrRollback(tx)
	}()

	return accountWallet, nil
}

func (w *WalletRepositoryImpl) Withdrawal(ctx context.Context, tx *sql.Tx, data web.WithdrawalRequest) (domain.Wallet, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	query := "SELECT COUNT(id) FROM wallet_t where reference_id = ? "
	rows, err := tx.QueryContext(ctx, query, data.ReferenceId)
	helpers.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		accountWallet := domain.Wallet{}
		return accountWallet, errors.New("Reference id must unique!")
	} else {
		query := "INSERT INTO wallet_t(wallet_id, amount, reference_id, type, created_at, created_by) values (?,?,?,?,?,?)"
		rows, err := tx.QueryContext(ctx, query, data.Wallet.Id, data.Amount, data.ReferenceId, "withdrawal", now.Format("2006-01-02 03:04:05"), data.Wallet.CustomerId)
		helpers.PanicIfError(err)
		defer rows.Close()
	}

	accountWallet := domain.Wallet{
		Id:          data.Wallet.Id,
		Amount:      data.Amount,
		CreatedBy:   data.Wallet.CustomerId,
		CreatedAt:   now,
		ReferenceId: data.ReferenceId,
	}

	//todo add to wallet with goroutine
	currBalance := data.Wallet.Balance
	newBalance := currBalance - data.Amount

	go func() {
		time.Sleep(5 * time.Second)
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/golang_mini_wallet")
		defer db.Close()
		res, err := db.Query("UPDATE wallet_m set balance = ? where id = ?", newBalance, data.Wallet.Id)
		defer res.Close()
		helpers.PanicIfError(err)
	}()

	return accountWallet, nil
}

func (w *WalletRepositoryImpl) Disabled(ctx context.Context, tx *sql.Tx, data web.DisabledRequest) domain.Wallet {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	query := "UPDATE wallet_m set is_active = false , disabled_at = NOW() where id = ?"
	_, err := tx.ExecContext(ctx, query, data.Id)
	helpers.PanicIfError(err)

	accountWallet := domain.Wallet{
		Id:         data.Id,
		Balance:    data.Balance,
		CreatedBy:  data.CustomerId,
		CustomerId: data.CustomerId,
		DisabledAt: now,
	}
	return accountWallet
}
