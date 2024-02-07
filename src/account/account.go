package account

import (
	"fmt"
	"sync"
)

type Account struct {
	IBAN    string // Unique account number in IBAN format
	balance int    // Account Balance
	status  string // Account status ("active" or "blocked")
	mutex   *sync.Mutex
}

func NewAccount(IBAN string) Account {
	return Account{IBAN: IBAN, balance: 0}
}

func (account *Account) GetUid() string {
	return account.IBAN
}

func (account *Account) GetBalance() int {
	return account.balance
}

// Withdraw subtracts funds from the account.
func (account *Account) Withdraw(amount int) error {
	account.mutex.Lock()
	defer account.mutex.Unlock()
	if amount > account.balance {
		return fmt.Errorf("insufficient balance in account %s", account.IBAN)
	}
	account.balance -= amount
	return nil
}

// Deposit adds funds to the account.
func (account *Account) Deposit(amount int) {
	account.mutex.Lock()
	defer account.mutex.Unlock()
	account.balance += amount
}
