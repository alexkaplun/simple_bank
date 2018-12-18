package handlers

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

import (
	"testing"
)

type TestCase struct {
	id       int
	input    interface{}
	expected interface{}
	positive bool
}

func Test_balanceInt64ToString(t *testing.T) {

	cases := []TestCase{
		{id: 1, input: "0", expected: "0", positive: true},
		{id: 2, input: "2", expected: "0.02", positive: true},
		{id: 3, input: "19", expected: "0.19", positive: true},
		{id: 4, input: "119", expected: "1.19", positive: true},
		{id: 5, input: "90000", expected: "900.00", positive: true},
		{id: 6, input: "129009090909010", expected: "1290090909090.10", positive: true},
	}

	for _, item := range cases {
		got := balanceInt64ToString(item.input.(string))
		if item.positive && got != item.expected.(string) {
			t.Errorf("Error in [%d]: expected: %s, got %s", item.id, item.expected.(string), got)
		}
	}
}

func Test_stringToBalanceInt64(t *testing.T) {
	cases := []TestCase{
		{id: 1, input: "0", expected: int64(0), positive: true},
		{id: 2, input: "2", expected: int64(200), positive: true},
		{id: 3, input: "0.19", expected: int64(19), positive: true},
		{id: 4, input: "119", expected: int64(11900), positive: true},
		{id: 5, input: "90000", expected: int64(9000000), positive: true},
		{id: 6, input: "290.75", expected: int64(29075), positive: true},
		{id: 7, input: "100.5", positive: false},
		// negative cases
		{id: 8, input: "290,75", positive: false},
		{id: 9, input: "-10", positive: false},
		//more than int64 max value
		{id: 10, input: "9223372036854775807", positive: false},
		{id: 11, input: "922337203685477580.47", positive: false},
		//less than int64 min value
		{id: 12, input: "-922337203685477580.7", positive: false},
		{id: 13, input: "ffuu", positive: false},
	}

	for _, item := range cases {
		got, err := stringToBalanceInt64(item.input.(string))
		if !item.positive && err == nil {
			t.Errorf("Error in [%d]: expected error, got %d", item.id, got)
		} else if item.positive && err != nil {
			t.Errorf("Error in [%d]: expected %d, got error: %s", item.id, got, err)
		} else if item.positive && got != item.expected {
			t.Errorf("Error in [%d]: expected %#v, got %d", item.id, item.expected, got)
		}
	}

}

