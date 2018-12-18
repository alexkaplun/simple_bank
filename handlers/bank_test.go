package handlers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type TestCase struct {
	input  interface{}
	output interface{}
	err    error
}

func Test_balanceInt64ToString(t *testing.T) {
	assert := assert.New(t)

	cases := map[string]TestCase{
		"zero":                {input: "0", output: "0", err: nil},
		"1 digit after zero":  {input: "2", output: "0.02", err: nil},
		"2 digits after zero": {input: "19", output: "0.19", err: nil},
		"1.19":                {input: "119", output: "1.19", err: nil},
		"Regular number":      {input: "90000", output: "900.00", err: nil},
		"Big number":          {input: "129009090909010", output: "1290090909090.10", err: nil},
	}

	for key, item := range cases {
		output := balanceInt64ToString(item.input.(string))
		assert.Equal(output, item.output.(string), key)
	}
}

func Test_stringToBalanceInt64(t *testing.T) {
	assert := assert.New(t)

	cases := map[string]TestCase{
		"zero":                       {input: "0", output: int64(0), err: nil},
		"two":                        {input: "2", output: int64(200), err: nil},
		"0.19":                       {input: "0.19", output: int64(19), err: nil},
		"reg 1":                      {input: "119", output: int64(11900), err: nil},
		"reg 2":                      {input: "90000", output: int64(9000000), err: nil},
		"float 2":                    {input: "290.75", output: int64(29075), err: nil},
		"float 1 symbol after comma": {input: "100.5", output: int64(-1), err: errors.New("can not parse balance")},
		"wrong separator":            {input: "290,75", output: int64(-1), err: &strconv.NumError{}},
		"negative":                   {input: "-10", output: int64(-1), err: errors.New("can not parse balance")},
		"overflow":                   {input: "9223372036854775807", output: int64(-1), err: errors.New("overflow, balance too high")},
		"overflow float":             {input: "922337203685477580.47", output: int64(-1), err: &strconv.NumError{}},
		"negative overflow":          {input: "-922337203685477580.7", output: int64(-1), err: errors.New("can not parse balance")},
		"not a number":               {input: "ffuu", output: int64(-1), err: errors.New("can not parse balance")},
	}

	for key, item := range cases {
		output, err := stringToBalanceInt64(item.input.(string))
		assert.IsType(int64(0), output, key)
		assert.IsType(item.err, err, key)
		assert.Equal(item.output, output, key)
	}

}
