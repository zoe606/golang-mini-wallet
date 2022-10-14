package v1

import (
	"github.com/julienschmidt/httprouter"
	"golang-mini-wallet/helpers"
	"golang-mini-wallet/model/domain"
	"golang-mini-wallet/model/web"
	v1 "golang-mini-wallet/services/v1"
	"net/http"
	"strconv"
	"strings"
)

type ApiControllerImpl struct {
	WalletService v1.WalletService
}

func NewApiControllerImpl(walletService v1.WalletService) *ApiControllerImpl {
	return &ApiControllerImpl{WalletService: walletService}
}

type ErrWarper struct {
	Error any `json:"error"`
}

type ErrCustomerNotInclude struct {
	CustomerXid []string `json:"customer_xid"`
}

type ErrReferencIdNotInclude struct {
	ReferenceId string `json:"reference_id"`
}

func (a *ApiControllerImpl) Init(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//helpers.ReadFormRequestBody(r, &payload)
	//
	//bodyBytes, _ := io.ReadAll(r.Body)
	//json.Unmarshal(bodyBytes, &payload)

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil || len(r.PostFormValue("customer_xid")) == 0 {
		//helpers.PanicIfError(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		//m := []string{"Missing data for required field."}
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: ErrCustomerNotInclude{
					CustomerXid: []string{"Missing data for required field."},
				},
			},
		}

		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	payload := web.InitRequest{
		CustomerXid: r.PostFormValue("customer_xid"),
	}

	initRes := a.WalletService.Init(r.Context(), payload)

	webResponse := web.Response{
		Code:   200,
		Status: "success",
		Data:   initRes,
	}

	helpers.WriteToResponseBody(w, webResponse)
}

func (a *ApiControllerImpl) Enable(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Token ")
	reqToken = splitToken[1]

	payload := web.TokenRequest{
		Token: reqToken,
	}

	enabledRes := a.WalletService.Enable(r.Context(), payload)

	if enabledRes.Wallet.Id == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "Already enabled",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	} else {
		webResponse := web.Response{
			Code:   200,
			Status: "success",
			Data:   enabledRes,
		}

		helpers.WriteToResponseBody(w, webResponse)
	}

}

func (a *ApiControllerImpl) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Token ")
	reqToken = splitToken[1]

	payload := web.TokenRequest{
		Token: reqToken,
	}

	getWallet := a.WalletService.Get(r.Context(), payload)

	if getWallet.Wallet.Id == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "Disabled",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	} else {
		webResponse := web.Response{
			Code:   200,
			Status: "success",
			Data:   getWallet,
		}

		helpers.WriteToResponseBody(w, webResponse)
	}
}

func (a *ApiControllerImpl) Deposit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Token ")
	reqToken = splitToken[1]

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil || len(r.PostFormValue("reference_id")) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		//m := []string{"Missing data for required field."}
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: ErrReferencIdNotInclude{
					ReferenceId: "Missing data for required field.",
				},
			},
		}

		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	token := web.TokenRequest{
		Token: reqToken,
	}

	getWallet := a.WalletService.Get(r.Context(), token)

	if getWallet.Wallet.Id == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "Disabled",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	amount, _ := strconv.Atoi(r.PostFormValue("amount"))
	payload := web.DepositRequest{
		Amount:      amount,
		ReferenceId: r.PostFormValue("reference_id"),
		Token:       domain.Token{Token: reqToken},
		Wallet: domain.Wallet{
			Id:         getWallet.Wallet.Id,
			CustomerId: getWallet.Wallet.CustomerId,
			Balance:    getWallet.Wallet.Balance,
		},
	}

	depositRes := a.WalletService.Deposit(r.Context(), payload)

	webResponse := web.Response{
		Code:   200,
		Status: "success",
		Data:   depositRes,
	}

	helpers.WriteToResponseBody(w, webResponse)
}

func (a *ApiControllerImpl) Withdrawal(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Token ")
	reqToken = splitToken[1]

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil || len(r.PostFormValue("reference_id")) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: ErrReferencIdNotInclude{
					ReferenceId: "Missing data for required field.",
				},
			},
		}

		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	token := web.TokenRequest{
		Token: reqToken,
	}

	getWallet := a.WalletService.Get(r.Context(), token)

	if getWallet.Wallet.Id == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "Disabled",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	amount, _ := strconv.Atoi(r.PostFormValue("amount"))

	if amount > getWallet.Wallet.Balance {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "Balance not available",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	payload := web.WithdrawalRequest{
		Amount:      amount,
		ReferenceId: r.PostFormValue("reference_id"),
		Token:       domain.Token{Token: reqToken},
		Wallet: domain.Wallet{
			Id:         getWallet.Wallet.Id,
			CustomerId: getWallet.Wallet.CustomerId,
			Balance:    getWallet.Wallet.Balance,
		},
	}

	withdrawalRes := a.WalletService.Withdrawal(r.Context(), payload)

	webResponse := web.Response{
		Code:   200,
		Status: "success",
		Data:   withdrawalRes,
	}

	helpers.WriteToResponseBody(w, webResponse)
}

func (a *ApiControllerImpl) Disabled(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Token ")
	reqToken = splitToken[1]

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil || len(r.PostFormValue("is_disabled")) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "Missing data for required field.",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	token := web.TokenRequest{
		Token: reqToken,
	}

	getWallet := a.WalletService.Get(r.Context(), token)

	if getWallet.Wallet.Status == "disabled" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "Wallet already disabled!",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	var status bool
	if r.PostFormValue("is_disabled") == "true" {
		status = true
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: "this endpoint purpose only for disabled wallet!",
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return
	}

	payload := web.DisabledRequest{
		Id:         getWallet.Wallet.Id,
		IsDisabled: status,
		Balance:    getWallet.Wallet.Balance,
		CustomerId: getWallet.Wallet.CustomerId,
	}
	disabledRes := a.WalletService.Disabled(r.Context(), payload)

	webResponse := web.Response{
		Code:   200,
		Status: "success",
		Data:   disabledRes,
	}

	helpers.WriteToResponseBody(w, webResponse)
}
