package bank

import (
	"errors"
	"github.com/JohnCGriffin/overflow"
	"github.com/google/uuid"
	"strconv"
	"sync"
	"time"
)

type Account struct {
	createdAt time.Time
	updatedAt time.Time
	balance   int64
}

type Bank struct {
	accounts map[uuid.UUID]*Account
	mu       *sync.Mutex
}

var bank = &Bank{
	accounts: make(map[uuid.UUID]*Account),
	mu:       &sync.Mutex{},
}

func GetBank() *Bank {
	return bank
}

func (b *Bank) CreateAccount(balance int64) (uuid.UUID, error) {
	if balance < 0 {
		return uuid.Nil, errors.New("сan not be negative balance")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	newId := uuid.New()
	if _, ok := b.accounts[newId]; ok {
		return uuid.Nil, errors.New("сan not generate Account ID")
	}

	b.accounts[newId] = &Account{
		createdAt: time.Now(),
		updatedAt: time.Now(),
		balance:   balance,
	}

	return newId, nil
}

func (b *Bank) GetAccountBalance(id uuid.UUID) (string, error) {
	// checking if account exists
	b.mu.Lock()
	defer b.mu.Unlock()

	ac, ok := b.accounts[id]
	if !ok {
		return "", errors.New("no account found")
	}

	return strconv.FormatInt(ac.balance, 10), nil
}

func (b *Bank) Transfer(from uuid.UUID, to uuid.UUID, amount int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// checking 0 or negative amount
	if amount <= 0 {
		return errors.New("can not be negative balance")
	}

	if _, ok := b.accounts[to]; !ok {
		return errors.New("terminating account not found")
	}

	if _, ok := b.accounts[from]; !ok {
		return errors.New("originating account not found")
	}

	fromBalance := b.accounts[from].balance
	toBalance := b.accounts[to].balance

	// Check if from has enough balance
	if fromBalance < amount {
		return errors.New("originating balance not enough")
	}

	// check for overflow after operation
	addRes, ok := overflow.Add64(toBalance, amount)
	if !ok {
		return errors.New("overflow of to balance")
	}

	//all validated, let's transfer
	b.accounts[from].balance = fromBalance - amount
	b.accounts[from].updatedAt = time.Now()

	b.accounts[to].balance = addRes
	b.accounts[to].updatedAt = time.Now()

	return nil

}
