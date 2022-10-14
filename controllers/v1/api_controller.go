package v1

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ApiController interface {
	Init(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Enable(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Get(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Deposit(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Withdrawal(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Disabled(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}
