package main

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"golang-mini-wallet/app"
	controller "golang-mini-wallet/controllers/v1"
	"golang-mini-wallet/helpers"
	"golang-mini-wallet/middleware"
	"golang-mini-wallet/repositories"
	service "golang-mini-wallet/services/v1"
	"net/http"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	return &http.Server{
		Addr:    "localhost:8010",
		Handler: authMiddleware,
	}
}

func main() {
	server := InitServer()
	err := server.ListenAndServe()
	helpers.PanicIfError(err)
}

func InitServer() *http.Server {
	WalletRepository := repositories.NewWalletRepository()
	db := app.NewDB()
	validate := validator.New()
	walletService := service.NewWallerServiceImpl(WalletRepository, db, validate)
	apiController := controller.NewApiControllerImpl(walletService)
	router := app.NewRouter(apiController)
	authMiddleware := middleware.NewAuthMiddleware(router)
	server := NewServer(authMiddleware)
	return server
}
