package handlers_test

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"simple_bank/server"
	"strings"
	"testing"
)

type TestCase struct {
	id       int
	input    string
	expected int
	positive bool
}

func TestCreateAccountHandler(t *testing.T) {
	cases := []TestCase{
		// positive
		{id: 1, input: `{"balance" : "0.08"}`, expected: http.StatusOK, positive: true},
		{id: 2, input: `{"balance" : "3.80"}`, expected: http.StatusOK, positive: true},
		{id: 3, input: `{"balance" : "0"}`, expected: http.StatusOK, positive: true},
		{id: 4, input: `{"balance" : "100"}`, expected: http.StatusOK, positive: true},
		{id: 5, input: `{"balance" : "99999.99"}`, expected: http.StatusOK, positive: true},
		/// negative
		{id: 6, input: `{"balance" : "99999.999"}`, expected: http.StatusUnprocessableEntity, positive: false},
		{id: 7, input: `{"balance" : "100"`, expected: http.StatusBadRequest, positive: false},
		{id: 8, input: `{"balance" : "9223372036854775807"}`, expected: http.StatusUnprocessableEntity, positive: false},
		{id: 9, input: `{"balance" : "290,75"`, expected: http.StatusBadRequest, positive: false},
		{id: 10, input: `{"balance" : "-50"}`, expected: http.StatusUnprocessableEntity, positive: false},
		{id: 11, input: `{"balance" : "-0.12"}`, expected: http.StatusUnprocessableEntity, positive: false},
		{id: 12, input: `{"balance" : "-50"sadas}`, expected: http.StatusBadRequest, positive: false},
	}

	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	for _, item := range cases {

		url := "/createAccount"

		req := httptest.NewRequest("PUT", url, strings.NewReader(item.input))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != item.expected {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				item.id, w.Code, item.expected)
		}
		/*
				resp := w.Result()
				body, _ := ioutil.ReadAll(resp.Body)

				bodyStr := string(body)
				if bodyStr != item.Response {
					t.Errorf("[%d] wrong Response: got %+v, expected %+v",
						caseNum, bodyStr, item.Response)
				}*/
	}

}

func TestGetBalanceByIdHandler(t *testing.T) {

}
