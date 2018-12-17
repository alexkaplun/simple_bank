package bank

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

import (
	"github.com/google/uuid"
	"strconv"
	"sync"
	"testing"
	"time"
)

var ids = []string{
	"714252f0-718b-44d9-a349-7bb30276bd95",
	"fc87c82f-46a1-4a38-b487-289bad552b65",
	"ea90026f-3272-4f2b-bef3-c84293569828",
	"c6a48141-1fcd-404c-9a9a-9a5a95e0d2af",
	"b636ab15-1ded-4e55-9fa4-e773c11dc0d6",
	"bf46513e-a4bb-4fe8-b80d-e26c71fc55e4",
	"05fed8ad-3581-48f3-89a8-fcb29b672fe0",
	"cd49335f-71bb-4856-b3be-6babf437064e",
	"b1e6ded9-d6b1-4966-a537-fbb10a891b67",
	////////////////////////////////////////
	"bc7b1e1b-e4d9-4176-b606-cf58847a0d22",
}

var balances = []string{
	"0",
	"1",
	"1000",
	"50000",
	"5000000",
	"5000000",
	"1000000000000000",
	"12121212",
	"34234234234",
	///////////////////////////////////////
	"50",
}

var testBank = &Bank{
	accounts: make(map[uuid.UUID]*Account),
	mu:       &sync.Mutex{},
}

func init() {
	for i := 0; i < 9; i++ {
		b, _ := strconv.ParseInt(balances[i], 10, 64)

		acc := Account{
			createdAt: time.Now(),
			updatedAt: time.Now(),
			balance:   b,
		}
		str, _ := uuid.Parse(ids[i])
		testBank.accounts[str] = &acc
	}
}

func TestBank_CreateAccount(t *testing.T) {
	//fmt.Println(len(testBank.accounts))

	var b int64 = 1000
	uid, err := testBank.CreateAccount(b)

	if err != nil {
		t.Errorf("Can't create account with uid %s", uid.String())
	} else {

		createdBalance := testBank.accounts[uid].balance
		if createdBalance != b {
			t.Errorf("New account balance wrong. Got %d, expected %d", createdBalance, b)
		}
	}
}

func TestBank_CreateAccountBalances(t *testing.T) {
	// test zero balance
	var b int64
	b = 0
	uid, err := testBank.CreateAccount(b)

	if err != nil {
		t.Errorf("Can't create account with uid %s", uid.String())
	} else {

		createdBalance := testBank.accounts[uid].balance
		if createdBalance != b {
			t.Errorf("New account balance wrong. Got %d, expected %d", createdBalance, b)
		}

	}

	// test very big balance
	b = 999999999999999999
	uid, err = testBank.CreateAccount(b)

	if err != nil {
		t.Errorf("Can't create account with uid %s", uid.String())
	} else {

		createdBalance := testBank.accounts[uid].balance
		if createdBalance != b {
			t.Errorf("New account balance wrong. Got %d, expected %d", createdBalance, b)
		}
	}

	// test  negative balance
	b = -10
	uid, err = testBank.CreateAccount(b)

	if err == nil {
		t.Errorf("Should not allow create negative account")
	}
}

func TestBank_GetAccountBalance(t *testing.T) {
	uid, _ := uuid.Parse("05fed8ad-3581-48f3-89a8-fcb29b672fe0")
	gotBalance, _ := testBank.GetAccountBalance(uid)

	if gotBalance != "1000000000000000" {
		t.Errorf("Get Balance error. Got %s expected 1000000000000000", gotBalance)
	}
}

func TestBank_GetAccountBalanceNoAccount(t *testing.T) {
	uid, _ := uuid.Parse("11111111-1111-1111-1111-1111111111")
	_, err := testBank.GetAccountBalance(uid)

	if err == nil {
		t.Errorf("Get balance error, got balance on non existing account")
	}
}

func TestBank_TransferFromNotExists(t *testing.T) {
	uidFrom, _ := uuid.Parse("11111111-1111-1111-1111-1111111111")
	uidTo, _ := uuid.Parse(ids[0])

	var amount int64 = 100
	if err := testBank.Transfer(uidFrom, uidTo, amount); err == nil {
		t.Errorf("Transfer error, transfer FROM non-existing account")
	}

}

func TestBank_TransferToNotExists(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[0])
	uidTo, _ := uuid.Parse("11111111-1111-1111-1111-1111111111")

	var amount int64 = 100
	if err := testBank.Transfer(uidFrom, uidTo, amount); err == nil {
		t.Errorf("Transfer error, transfer TO non-existing account")
	}
}

func TestBank_TransferZeroOrNegative(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[0])
	uidTo, _ := uuid.Parse(ids[1])

	var amountZero int64 = 0
	if err := testBank.Transfer(uidFrom, uidTo, amountZero); err == nil {
		t.Errorf("Transfer error, transfer zero amount")
	}

	var amountNegative int64 = -100
	if err := testBank.Transfer(uidFrom, uidTo, amountNegative); err == nil {
		t.Errorf("Transfer error, transfer negative")
	}
}

func TestBank_Transfer(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[3]) // balance = 50000
	uidTo, _ := uuid.Parse(ids[2])   // balace = 1000

	var amount int64 = 5000

	testBank.Transfer(uidFrom, uidTo, amount)

	if fromVal, _ := testBank.GetAccountBalance(uidFrom); fromVal != "45000" {
		t.Errorf("Wrong FROM balance after transfer. Got %s, expected 45000", fromVal)
	}

	if toBal, _ := testBank.GetAccountBalance(uidTo); toBal != "6000" {
		t.Errorf("Wrong TO balance after transfer. Got %s, expected 6000", toBal)
	}
}

func TestBank_TransferNotEnoughBalance(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[3]) // balance = 50000
	uidTo, _ := uuid.Parse(ids[2])   // balace = 1000

	var amount int64 = 50001

	err := testBank.Transfer(uidFrom, uidTo, amount)

	if err == nil {
		t.Errorf("Error transfer inadequste balance")
	}
}
