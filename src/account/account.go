package account

import (
	"fmt"
	"sync"
)

type Status int

const (
	Active Status = iota
	Blocked
)

func (s Status) String() string {
	switch s {
	case Active:
		return "active"
	case Blocked:
		return "blocked"
	}
	return "unknown"
}

type Account struct {
	IBAN     string // Unique account number in IBAN format
	balance  int    // Account Balance
	status   Status // Account status ("active" or "blocked")
	canBlock bool   // If account can be blocked
	mutex    sync.Mutex
}

type JsonAccount struct {
	IBAN    string `json:"IBAN"`
	Balance int    `json:"balance"`
	Status  string `json:"status"`
}

func NewAccount(IBAN string, canBlock bool) Account {
	return Account{IBAN: IBAN, balance: 0, status: Active, canBlock: canBlock}
}

func (account *Account) GetUid() string {
	return account.IBAN
}

func (account *Account) ToJson() *JsonAccount {
	return &JsonAccount{account.IBAN, account.balance, account.status.String()}
}

func (account *Account) GetBalance() int {
	return account.balance
}

// Withdraw subtracts funds from the account.
func (account *Account) withdraw(amount int) error {
	if account.status == Blocked {
		return fmt.Errorf("failed to withdraw account %s is blocked", account.IBAN)
	}
	if amount > account.balance {
		return fmt.Errorf("insufficient balance in account %s", account.IBAN)
	}
	account.balance -= amount
	return nil
}

// Deposit adds funds to the account.
func (account *Account) Deposit(amount int, lock bool) error {
	if lock {
		account.mutex.Lock()
		defer account.mutex.Unlock()
	}

	if account.status == Blocked {
		return fmt.Errorf("failed to deposit. account %s is blocked", account.IBAN)
	}
	account.balance += amount
	return nil
}

func (account *Account) TransferTo(toAccount *Account, amount int) error {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	toAccount.mutex.Lock()
	defer toAccount.mutex.Unlock()

	if err := account.withdraw(amount); err != nil {
		return err
	}

	if err := toAccount.Deposit(amount, false); err != nil {
		// Rollback withdrawal if deposit fails
		account.balance += amount
		return err
	}

	return nil
}

func (account *Account) Block() error {
	account.mutex.Lock()
	defer account.mutex.Unlock()
	if !account.canBlock {
		return fmt.Errorf("account %s can't be blocked", account.IBAN)
	}
	account.status = Blocked
	return nil
}

func (account *Account) Unblock() {
	account.mutex.Lock()
	defer account.mutex.Unlock()
	account.status = Active
}
