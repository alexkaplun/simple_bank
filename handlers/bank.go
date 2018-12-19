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

type TransferRequest struct {
	From string `json:"from"`
	To string `json:"to"`
	Amount string `json:"amount"`
}

type JSONResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

func CreateAccountHandler(c *gin.Context) {

	var r CreateAccountRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	balance, err := stringToBalanceInt64(r.Balance)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	_bank := bank.GetBank()
	uid, err := _bank.CreateAccount(balance)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, &JSONResponse{0, CreateAccountResponse{uid.String()}})
}

func GetBalanceByIdHandler(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	_bank := bank.GetBank()
	balance, err := _bank.GetAccountBalance(uid)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, &JSONResponse{ 0, GetBalanceResponse{balanceInt64ToString(balance)}})
}

func TransferHandler(c *gin.Context) {

	var r TransferRequest

	// Bind JSON
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	// validate balance
	amount, err := stringToBalanceInt64(r.Amount)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	// validating from UID
	from, err := uuid.Parse(r.From)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	// validating to UID
	to, err := uuid.Parse(r.To)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	// attempt to transfer
	_bank := bank.GetBank()
	err = _bank.Transfer(from, to, amount)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, &JSONResponse{-1, ErrorResponse{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, &JSONResponse{ 0, nil})
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
			if mulRes, ok := overflow.Mul64(parsed, 100); !ok {
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
