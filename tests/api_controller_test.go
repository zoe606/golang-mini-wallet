package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"golang-mini-wallet/app"
	controllers "golang-mini-wallet/controllers/v1"
	"golang-mini-wallet/helpers"
	"golang-mini-wallet/middleware"
	"golang-mini-wallet/repositories"
	services "golang-mini-wallet/services/v1"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const BASE_URL = "http://localhost:8010"
const TOKEN = "Token 4IbHYDlYGXLyKEkvemW1A6f6W"

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/golang_mini_wallet")
	helpers.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	walletRepository := repositories.NewWalletRepository()
	walletService := services.NewWallerServiceImpl(walletRepository, db, validate)
	apiController := controllers.NewApiControllerImpl(walletService)

	router := app.NewRouter(apiController)

	return middleware.NewAuthMiddleware(router)
}

func TestInitSucces(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	err := writer.WriteField("customer_xid", helpers.RandomString(15))
	if err != nil {
		return
	}
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/init", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 200, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "success", bodyResponse["status"])
}

func TestInitFailed(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	err := writer.WriteField("customer_xid", "")
	if err != nil {
		return
	}
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/init", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 400, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "fail", bodyResponse["status"])
	//fmt.Println(bodyResponse["data"].(map[string]interface{})["error"])
	//assert.Equal(t, "error", bodyResponse["data"].(map[string]interface{})["error"])
}

func TestEnableSucces(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 200, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "success", bodyResponse["status"])
}

func TestEnableFailed(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 400, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "fail", bodyResponse["status"])
}

func TestGetWalletSucess(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.Close()

	request := httptest.NewRequest(http.MethodGet, BASE_URL+"/api/v1/wallet", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 200, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "success", bodyResponse["status"])
}

func TestGetWalletDisabled(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	writer.Close()

	request := httptest.NewRequest(http.MethodGet, BASE_URL+"/api/v1/wallet", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 400, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "fail", bodyResponse["status"])
}

func TestDepositSucess(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	_ = writer.WriteField("amount", "123456")
	_ = writer.WriteField("reference_id", helpers.RandomString(15))
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet/deposits", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 200, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "success", bodyResponse["status"])
}

func TestDepositFailed(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	_ = writer.WriteField("amount", "123456")
	_ = writer.WriteField("reference_id", "")
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet/deposits", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 400, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "fail", bodyResponse["status"])
}

func TestWithDrawalSucess(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	_ = writer.WriteField("amount", "123456")
	_ = writer.WriteField("reference_id", helpers.RandomString(15))
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet/withdrawals", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 200, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "success", bodyResponse["status"])
}

func TestWithDrawalFailed(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	_ = writer.WriteField("amount", "123456")
	_ = writer.WriteField("reference_id", "")
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet/withdrawals", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 400, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "fail", bodyResponse["status"])
}

func TestDisabledWalletSucess(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	_ = writer.WriteField("is_disabled", "true")
	writer.Close()

	request := httptest.NewRequest(http.MethodPatch, BASE_URL+"/api/v1/wallet", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 200, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "success", bodyResponse["status"])
}

func TestDisabledFailed(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	reqBody := new(bytes.Buffer)
	writer := multipart.NewWriter(reqBody)
	_ = writer.WriteField("is_disabled", "false")
	writer.Close()

	request := httptest.NewRequest(http.MethodPatch, BASE_URL+"/api/v1/wallet", reqBody)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", TOKEN)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	var bodyResponse map[string]interface{}
	json.Unmarshal(body, &bodyResponse)

	assert.Equal(t, 400, int(bodyResponse["code"].(float64)))
	assert.Equal(t, "fail", bodyResponse["status"])
}
