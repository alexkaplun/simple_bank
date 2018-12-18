package handlers

import (
	"errors"
	"github.com/JohnCGriffin/overflow"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"simple_bank/models/bank"
	"strconv"
	"strings"
)

var validBalanceRegexp = regexp.MustCompile(`^[0-9]+(.[0-9][0-9])?$`)

type CreateAccountRequest struct {
	Balance string `json:"balance"`
}

type CreateAccountResponse struct {
	Uid string `json:"account_id"`
}

type GetBalanceResponse struct {
	Balance string `json:"balance"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

func CreateAccountHandler(c *gin.Context) {

	var r CreateAccountRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &ErrorResponse{err.Error()})
		return
	}

	balance, err := stringToBalanceInt64(r.Balance)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &ErrorResponse{err.Error()})
		return
	}

	bank := bank.GetBank()
	uid, err := bank.CreateAccount(balance)

	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, &CreateAccountResponse{uid.String()})
}

func GetBalanceByIdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, &ErrorResponse{"No account id provided"})
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{err.Error()})
		return
	}

	bank := bank.GetBank()
	balance, err := bank.GetAccountBalance(uid)

	if err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, &GetBalanceResponse{balanceInt64ToString(balance)})
}

func stringToBalanceInt64(s string) (int64, error) {

	if !validBalanceRegexp.MatchString(s) {
		return -1, errors.New("can not parse balance")
	}

	part := strings.Split(s, ".")
	if len(part) == 1 {
		parsed, err := strconv.ParseInt(part[0], 10, 64)
		if err != nil {
			return -1, err
		} else {
			if mulRes, ok  := overflow.Mul64(parsed, 100); !ok {
				return -1, errors.New("overflow, balance too high")
			} else {
				return mulRes, nil
			}
		}
	} else {
		// assuming len(part) can not be more than 2
		parsed, err := strconv.ParseInt(part[0]+part[1], 10, 64)
		if err != nil {
			return -1, err
		} else {
			return parsed, nil
		}
	}
}

func balanceInt64ToString(b string) string {
	var left, right string
	switch {
	case b == "0":
		return "0"
	case len(b) == 1:
		left, right = "0", "0"+b
	case len(b) == 2:
		left, right = "0", b
	default:
		left, right = b[:len(b)-2], b[len(b)-2:]
	}
	return left + "." + right
}
