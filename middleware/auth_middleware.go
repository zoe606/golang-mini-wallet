package middleware

import (
	"golang-mini-wallet/helpers"
	"golang-mini-wallet/model/web"
	"net/http"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if len(request.Header.Get("Authorization")) > 0 || request.URL.Path == "/api/v1/init" {
		middleware.Handler.ServeHTTP(writer, request)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		webResponse := web.Response{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized!",
		}
		helpers.WriteToResponseBody(writer, webResponse)
	}
}
