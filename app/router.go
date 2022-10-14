package app

import (
	"github.com/julienschmidt/httprouter"
	v1 "golang-mini-wallet/controllers/v1"
	"golang-mini-wallet/exceptions"
)

func NewRouter(apiController v1.ApiController) *httprouter.Router {
	router := httprouter.New()
	router.POST("/api/v1/init", apiController.Init)
	router.POST("/api/v1/wallet", apiController.Enable)
	router.GET("/api/v1/wallet", apiController.Get)
	router.POST("/api/v1/wallet/deposits", apiController.Deposit)
	router.POST("/api/v1/wallet/withdrawals", apiController.Withdrawal)
	router.PATCH("/api/v1/wallet", apiController.Disabled)
	router.PanicHandler = exceptions.ErrorHandler
	return router
}
