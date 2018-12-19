package handlers_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"simple_bank/handlers"
	bankModel "simple_bank/models/bank"
	"simple_bank/server"
	"strings"
	"testing"
)

type TestCaseStatusCode struct {
	input      string
	statusCode int
}

func TestCreateAccountHandler(t *testing.T) {

	cases := map[string]TestCaseStatusCode{
		"8 cents":      {input: `{"balance" : "0.08"}`, statusCode: http.StatusOK},
		"3.80":         {input: `{"balance" : "3.80"}`, statusCode: http.StatusOK},
		"zero":         {input: `{"balance" : "0"}`, statusCode: http.StatusOK},
		"Normal 100":   {input: `{"balance" : "100"}`, statusCode: http.StatusOK},
		"Normal float": {input: `{"balance" : "99999.99"}`, statusCode: http.StatusOK},
		/// negative
		"bad format float": {input: `{"balance" : "99999.999"}`, statusCode: http.StatusUnprocessableEntity},
		"malformed json":   {input: `{"baldfddfance" : "100"`, statusCode: http.StatusBadRequest},
		"overflow":         {input: `{"balance" : "9223372036854775807"}`, statusCode: http.StatusUnprocessableEntity},
		"wrong separator":  {input: `{"balance" : "290,75"}`, statusCode: http.StatusUnprocessableEntity},
		"negative int":     {input: `{"balance" : "-50"}`, statusCode: http.StatusUnprocessableEntity},
		"negative float":   {input: `{"balance" : "-0.12"}`, statusCode: http.StatusUnprocessableEntity},
		"wrong json":       {input: `{"balance" : "-50"sadas}`, statusCode: http.StatusBadRequest},
	}

	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	for key, item := range cases {

		url := "/createAccount"

		req := httptest.NewRequest("PUT", url, strings.NewReader(item.input))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, item.statusCode, w.Code, key)

		// read the response
		res := w.Result()
		body, _ := ioutil.ReadAll(res.Body)

		resp := new(handlers.JSONResponse)
		err := json.Unmarshal(body, resp)

		// assert json format
		assert.Nil(t, err, key)
		assert.IsType(t, &handlers.JSONResponse{}, resp, key)
	}

}

func TestGetBalanceByIdHandler_StatusCodesNegative(t *testing.T) {
	cases := map[string]TestCaseStatusCode{
		"invalid id":       {input: "123", statusCode: http.StatusUnprocessableEntity},
		"non-existing UID": {input: "e4517c6a-b2e2-4257-997b-2e5cc7356483", statusCode: http.StatusUnprocessableEntity},
		"empty uid":        {input: "", statusCode: http.StatusNotFound},
	}

	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	for key, item := range cases {

		url := "/balance/" + item.input

		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, item.statusCode, w.Code, key)
	}
}

func TestGetBalanceByIdHandler(t *testing.T) {
	//Create the account
	bank := bankModel.GetBank()
	input := int64(123 * 100)
	expected := "123.00"
	uid, err := bank.CreateAccount(input)
	if err != nil {
		assert.Fail(t, "Can't create account")
	}

	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	url := "/balance/" + uid.String()
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)

	resp := &handlers.JSONResponse{
		Body: &handlers.GetBalanceResponse{},
	}
	if err := json.Unmarshal(body, resp); err != nil {
		assert.Fail(t, "Can't unmarshal get balance response")
	}

	assert.IsType(t, &handlers.JSONResponse{Body: &handlers.GetBalanceResponse{}}, resp)
	assert.Equal(t, expected, resp.Body.(*handlers.GetBalanceResponse).Balance)
}

func TestTransferHandler_Validation(t *testing.T) {

	cases := map[string]TestCaseStatusCode{
		"Bad JSON": {input: `{Iambadjson`, statusCode: http.StatusBadRequest},
		"Invalid Balance": {
			input:
			`{"from":"800227ee-362c-4382-809d-ea39f0807418","to":"e4517c6a-b2e2-4257-997b-2e5cc7356483","amount":"0.342401"}`,
			statusCode: http.StatusUnprocessableEntity},
		"Bad From UID": {
			input:
			`{"from":"800227ee-362c-4382-809d-ea39f007418","to":"e4517c6a-b2e2-4257-997b-2e5cc7356483","amount":"100"}`,
			statusCode: http.StatusUnprocessableEntity},
		"Bad To UID": {
			input:
			`{"from":"800227ee-362c-4382-809d-ea39f0807418","to":"e4517c6a-2e2-4257-997b-2e5cc7356483","amount":"1000"}`,
			statusCode: http.StatusUnprocessableEntity},
	}

	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	url := "/transfer"

	for key, item := range cases {
		req := httptest.NewRequest("POST", url, strings.NewReader(item.input))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, item.statusCode, w.Code, key)
	}

}

func TestTransferHandler_TransferError(t *testing.T) {
	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	url := "/transfer"

	bank := bankModel.GetBank()
	fromBalance, toBalance := int64(12000*100), int64(9000*100)

	uidTo, err := bank.CreateAccount(toBalance)
	if err != nil {
		assert.Fail(t, "Can't create account")
	}
	uidFrom, err := bank.CreateAccount(fromBalance)
	if err != nil {
		assert.Fail(t, "Can't create account")
	}

	// from account not found
	goodJSON := `{"from":"800227ee-362c-4382-809d-ea39f0807418","to":"` + uidTo.String() + `","amount":"100"}`

	req := httptest.NewRequest("POST", url, strings.NewReader(goodJSON))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	// to account not found
	goodJSON = `{"from":"` + uidFrom.String() + `","to":"800227ee-362c-4382-809d-ea39f0807418","amount":"100"}`

	req = httptest.NewRequest("POST", url, strings.NewReader(goodJSON))
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	// not enough balance
	goodJSON = `{"from":"` + uidFrom.String() + `","to":"` + uidTo.String() + `","amount":"100000000"}`

	req = httptest.NewRequest("POST", url, strings.NewReader(goodJSON))
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestTransferHandler(t *testing.T) {
	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	url := "/transfer"

	bank := bankModel.GetBank()
	// 12345.60 and 9000.30
	// expected 12300.00 and 9045.90
	fromBalance, toBalance := int64(1234560), int64(900030)
	amount := "45.60"

	uidTo, err := bank.CreateAccount(toBalance)
	if err != nil {
		assert.Fail(t, "Can't create To account")
	}
	uidFrom, err := bank.CreateAccount(fromBalance)
	if err != nil {
		assert.Fail(t, "Can't create From account")
	}

	// from account not found
	goodJSON := `{"from":"` + uidFrom.String() + `","to":"` + uidTo.String() + `","amount":"` + amount + `"}`

	req := httptest.NewRequest("POST", url, strings.NewReader(goodJSON))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// HTTP code mus be ok since expect transfer to be processed
	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)

	resp := &handlers.JSONResponse{Body: nil}

	if err := json.Unmarshal(body, resp); err != nil {
		assert.Fail(t, "Can't unmarshal get balance response")
	}

	// check json response format
	assert.IsType(t, &handlers.JSONResponse{Body: nil}, resp)

	// Check updated balances
	updatedFrom, _ := bank.GetAccountBalance(uidFrom)
	updatedTo, _ := bank.GetAccountBalance(uidTo)
	assert.Equal(t, "1230000", updatedFrom)
	assert.Equal(t, "904590", updatedTo)
}
