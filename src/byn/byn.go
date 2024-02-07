package byn

import (
	"fmt"
	"sync"
	"webByn/src/account"
	"webByn/src/repository"
	"webByn/src/util"
)

type PaymentService struct {
	ibanGenerator      util.IBANGenerator
	emissionAccount    account.Account // Счет для эмиссии денег
	destructionAccount account.Account // Счет для уничтожения денег
	accounts           repository.Repository[*account.Account]
}

var once sync.Once
var bynSystemInstance *PaymentService = nil

// GetBynSystem returns instance of PaymentService which must be singleton
func GetBynSystem() *PaymentService {
	once.Do(func() {
		var ibanGenerator = util.NewIBANGenerator("BY", 28-2)
		var accountsRepository = repository.NewRepository[*account.Account]()
		var emissionAccount = account.NewAccount(ibanGenerator.Generate())
		var destructionAccount = account.NewAccount(ibanGenerator.Generate())
		err := accountsRepository.Add(&emissionAccount)
		if err != nil {
			return
		}
		err = accountsRepository.Add(&destructionAccount)
		if err != nil {
			return
		}
		bynSystemInstance = &PaymentService{
			*ibanGenerator,
			emissionAccount,
			destructionAccount,
			*accountsRepository,
		}
	})
	return bynSystemInstance
}

// OpenAccount creates new account and returns new account instance
func (paymentService *PaymentService) OpenAccount() (*account.Account, error) {
	var iban = paymentService.ibanGenerator.Generate()
	var accountInstance = account.NewAccount(iban)
	err := paymentService.accounts.Add(&accountInstance)
	if err != nil {
		return nil, err
	}
	return &accountInstance, nil
}

// IssueMoney issues the specified amount to the "emission" account.
func (paymentService *PaymentService) IssueMoney(amount int) {
	paymentService.emissionAccount.Deposit(amount)
}

// Transfer sends a specified amount of money from a specified account to the "destruction" account.
func (paymentService *PaymentService) Transfer(fromIBAN string, amount int) error {
	fromAccount, err := paymentService.accounts.GetById(fromIBAN, false)
	if err != nil {
		return fmt.Errorf("account with IBAN %s not found", fromIBAN)
	}
	err = (*fromAccount).Withdraw(amount)
	if err != nil {
		return err
	}
	paymentService.destructionAccount.Deposit(amount)
	return nil
}
