package handlers_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"simple_bank/server"
	"strings"
	"testing"
)

type TestCaseStatusCode struct {
	input    string
	statusCode int
}

func TestCreateAccountHandler(t *testing.T) {

	cases := map[string]TestCaseStatusCode {
		"8 cents" : {input: `{"balance" : "0.08"}`, statusCode: http.StatusOK},
		"3.80": {input: `{"balance" : "3.80"}`, statusCode: http.StatusOK},
		"zero": {input: `{"balance" : "0"}`, statusCode: http.StatusOK  },
		"Normal 100": {input: `{"balance" : "100"}`,  statusCode: http.StatusOK},
		"Normal float":{input: `{"balance" : "99999.99"}`,  statusCode: http.StatusOK  },
		/// negative
		"bad format float" :{input: `{"balance" : "99999.999"}`,  statusCode: http.StatusUnprocessableEntity  },
		"malformed json":{input: `{"balance" : "100"`,  statusCode: http.StatusBadRequest  },
		"overflow":{input: `{"balance" : "9223372036854775807"}`,  statusCode: http.StatusUnprocessableEntity  },
		"wrong separator" :{input: `{"balance" : "290,75"`,  statusCode: http.StatusBadRequest  },
		"negative int":{input: `{"balance" : "-50"}`,  statusCode: http.StatusUnprocessableEntity  },
		"negative float":{input: `{"balance" : "-0.12"}`,  statusCode: http.StatusUnprocessableEntity  },
		"wrong json" :{input: `{"balance" : "-50"sadas}`,  statusCode: http.StatusBadRequest  },
	}

	// Switch to test mode and get the router
	gin.SetMode(gin.TestMode)
	r := server.NewRouter()

	for _, item := range cases {

		url := "/createAccount"

		req := httptest.NewRequest("PUT", url, strings.NewReader(item.input))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, item.statusCode, w.Code)

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
