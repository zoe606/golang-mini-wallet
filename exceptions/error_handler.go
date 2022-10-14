package exceptions

import (
	"github.com/go-playground/validator/v10"
	"golang-mini-wallet/helpers"
	"golang-mini-wallet/model/web"
	"net/http"
)

type ErrWarper struct {
	Error interface{} `json:"error"`
}

type ApiError struct {
	Field string
	Msg   string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if notFoundError(w, r, err) {
		return
	}

	if validationErrors(w, r, err) {
		return
	}
	internalServerError(w, r, err)
}

func validationErrors(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		webResponse := web.Response{
			Code:   http.StatusBadRequest,
			Status: "fail",
			Data: ErrWarper{
				Error: exception.Error(),
			},
		}
		helpers.WriteToResponseBody(w, webResponse)
		return true
	} else {
		return false
	}

	//var ve validator.ValidationErrors
	//m := map[string]interface{}{}
	//if errors.As(exception, &ve) {
	//	//out := make([]interface{}, len(ve))
	//	for _, fe := range ve {
	//		//out[i] = msgForTag(fe.Tag())
	//		//out[i] = []string{fieldForTag(fe.Field())}
	//		m[fieldForTag(fe.Field())] = []string{msgForTag(fe.Tag())}
	//	}
	//	webResponse := web.Response{
	//		Code:   http.StatusBadRequest,
	//		Status: "fail",
	//		Data: ErrWarper{
	//			Error: m,
	//		},
	//	}
	//	helpers.WriteToResponseBody(w, webResponse)
	//}
	//return true
}

func notFoundError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		webResponse := web.Response{
			Code:   http.StatusNotFound,
			Status: "not found",
			Data:   exception.Error,
		}
		helpers.WriteToResponseBody(w, webResponse)

		return true
	} else {
		return false
	}
}

func internalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	webResponse := web.Response{
		Code:   http.StatusInternalServerError,
		Status: "internal server error",
		Data:   err,
	}
	helpers.WriteToResponseBody(w, webResponse)
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "Missing data for required field."
	case "min":
		return "Min data for required field."
	}
	return ""
}

func fieldForTag(field string) string {
	switch field {
	case "CustomerXid":
		return "customer_xid"
	}
	return ""
}
