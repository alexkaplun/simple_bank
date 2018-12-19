package bank

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	"9223372036854775807",
}

var testBank = &Bank{
	accounts: make(map[uuid.UUID]*Account),
	mu:       &sync.Mutex{},
}

func init() {
	for i := 0; i <= 9; i++ {
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
	var b int64 = 1000
	uid, err := testBank.CreateAccount(b)

	assert.Nil(t, err)

	createdBalance := testBank.accounts[uid].balance
	assert.Equal(t, b, createdBalance)
}

func TestBank_CreateAccountBalances(t *testing.T) {
	// tests zero balance
	var b int64
	b = 0
	uid, err := testBank.CreateAccount(b)
	assert.Nil(t, err)

	createdBalance := testBank.accounts[uid].balance
	assert.Equal(t, b, createdBalance)

	// tests very big balance
	b = 999999999999999999
	uid, err = testBank.CreateAccount(b)
	assert.Nil(t, err)

	createdBalance = testBank.accounts[uid].balance
	assert.Equal(t, b, createdBalance)

	// tests  negative balance
	b = -10
	uid, err = testBank.CreateAccount(b)
	assert.NotNil(t, err)
}

func TestBank_GetAccountBalance(t *testing.T) {
	uid, _ := uuid.Parse("05fed8ad-3581-48f3-89a8-fcb29b672fe0")
	gotBalance, err := testBank.GetAccountBalance(uid)
	assert.Nil(t, err)
	assert.Equal(t, "1000000000000000", gotBalance)
}

func TestBank_GetAccountBalanceNoAccount(t *testing.T) {
	uid, _ := uuid.Parse("11111111-1111-1111-1111-1111111111")
	b, err := testBank.GetAccountBalance(uid)
	assert.NotNil(t, err)
	assert.Equal(t, "", b)
}

func TestBank_TransferFromNotExists(t *testing.T) {
	uidFrom, _ := uuid.Parse("11111111-1111-1111-1111-1111111111")
	uidTo, _ := uuid.Parse(ids[0])

	err := testBank.Transfer(uidFrom, uidTo, int64(100))
	assert.NotNil(t, err)

}

func TestBank_TransferToNotExists(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[0])
	uidTo, _ := uuid.Parse("11111111-1111-1111-1111-1111111111")

	err := testBank.Transfer(uidFrom, uidTo, int64(100))
	assert.NotNil(t, err)
}

func TestBank_TransferZeroOrNegative(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[0])
	uidTo, _ := uuid.Parse(ids[1])

	err := testBank.Transfer(uidFrom, uidTo, int64(0))
	assert.NotNil(t, err)

	err = testBank.Transfer(uidFrom, uidTo, int64(-100))
	assert.NotNil(t, err)
}

func TestBank_Transfer(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[3]) // balance = 50000
	uidTo, _ := uuid.Parse(ids[2])   // balace = 1000

	var amount int64 = 5000

	err := testBank.Transfer(uidFrom, uidTo, amount)
	assert.Nil(t, err)

	fromVal, err := testBank.GetAccountBalance(uidFrom)
	assert.Nil(t, err)
	assert.Equal(t, "45000", fromVal)

	toVal, err := testBank.GetAccountBalance(uidTo)
	assert.Nil(t, err)
	assert.Equal(t, "6000", toVal)
}

func TestBank_TransferNotEnoughBalance(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[3]) // balance = 50000
	uidTo, _ := uuid.Parse(ids[2])   // balace = 1000

	err := testBank.Transfer(uidFrom, uidTo, int64(50001))
	assert.NotNil(t, err)
}

func TestBank_TransferOverflow(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[2]) // balance = 1000
	uidTo, _ := uuid.Parse(ids[9])   // balace = 9223372036854775807

	err := testBank.Transfer(uidFrom, uidTo, int64(1000))
	assert.NotNil(t, err)
}

func TestBank_TransferSameAccount(t *testing.T) {
	uidFrom, _ := uuid.Parse(ids[2]) // balance = 1000
	uidTo, _ := uuid.Parse(ids[2])   // balace = 1000

	err := testBank.Transfer(uidFrom, uidTo, int64(500))
	assert.NotNil(t, err)
}
