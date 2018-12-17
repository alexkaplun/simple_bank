package bank

import (
	"errors"
	"github.com/google/uuid"
	"math/big"
	"sync"
	"time"
)

type Account struct {
	createdAt time.Time
	updatedAt time.Time
	balance   *big.Int
}

type Bank struct {
	accounts map[uuid.UUID]*Account
	mu       *sync.Mutex
}

var bank = &Bank{
	accounts: make(map[uuid.UUID]*Account),
	mu:       &sync.Mutex{}}

func GetBank() *Bank {
	return bank
}

func (b *Bank) CreateAccount(balance big.Int) (uuid.UUID, error) {
	if balance.Cmp(new(big.Int)) == -1 {
		return uuid.Nil, errors.New("Can't be negative balance")
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	newId := uuid.New()
	if _, ok := b.accounts[newId]; ok {
		return uuid.Nil, errors.New("Can't generate Account ID")
	}

	newBalance := balance

	b.accounts[newId] = &Account{
		createdAt: time.Now(),
		updatedAt: time.Now(),
		balance:   &newBalance,
	}

	return newId, nil
}

func (b *Bank) GetAccountBalance(id uuid.UUID) (string, error) {
	// checking if account exists
	b.mu.Lock()
	defer b.mu.Unlock()

	ac, ok := b.accounts[id]
	if !ok {
		return "", errors.New("No account found")
	}

	return ac.balance.String(), nil
}

func (b *Bank) Transfer(from uuid.UUID, to uuid.UUID, amount big.Int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// checking 0 or negative amount
	if amount.Cmp(new(big.Int)) <= 0 {
		return errors.New("Can't be negative balance")
	}

	if _, ok := b.accounts[to]; !ok {
		return errors.New("Terminating account not found")
	}

	if _, ok := b.accounts[from]; !ok {
		return errors.New("Originating account not found")
	} else {

		// Check if from has enough balance
		fromBalance := b.accounts[from].balance
		if amount.Cmp(fromBalance) == 1 {
			return errors.New("Originating balance not enough")
		}

		//all validated, let's transfer
		fromBalance.Sub(fromBalance, &amount)
		b.accounts[from].updatedAt = time.Now()

		toBalance := b.accounts[to].balance
		toBalance.Add(toBalance, &amount)
		b.accounts[to].updatedAt = time.Now()

		return nil
	}

}
