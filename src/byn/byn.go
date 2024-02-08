package byn

import (
	"encoding/json"
	"fmt"
	"sync"
	"webByn/src/account"
	"webByn/src/repository"
	"webByn/src/util"
)

type TransferRequest struct {
	From   string `json:"fromIBAN"`
	To     string `json:"toIBAN"`
	Amount int    `json:"amount"`
}

type PaymentService struct {
	ibanGenerator      util.IBANGenerator
	emissionAccount    *account.Account                                               // Account for emission money that can't be blocked
	destructionAccount *account.Account                                               // Account for destruction money that can't be blocked
	accounts           *repository.Repository[*account.Account, *account.JsonAccount] // wrapper for accounts array
}

var once sync.Once
var bynSystemInstance *PaymentService = nil

// GetBynSystem returns instance of PaymentService which must be singleton
func GetBynSystem() *PaymentService {
	once.Do(func() {
		var ibanGenerator = util.NewIBANGenerator("BY", 28-2)
		var accountsRepository = repository.NewRepository[*account.Account]()
		var emissionAccount = account.NewAccount(ibanGenerator.Generate(), false, true)
		var destructionAccount = account.NewAccount(ibanGenerator.Generate(), false, false)
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
			&emissionAccount,
			&destructionAccount,
			accountsRepository,
		}
	})
	return bynSystemInstance
}

// GetEmissionAccountNumber returns IBAN if "emission" account
func (paymentService *PaymentService) GetEmissionAccountNumber() string {
	return paymentService.emissionAccount.IBAN
}

// GetDestructionAccountNumber returns IBAN if "destruction" account
func (paymentService *PaymentService) GetDestructionAccountNumber() string {
	return paymentService.destructionAccount.IBAN
}

// OpenAccount creates new account and returns IBAN of new instance
func (paymentService *PaymentService) OpenAccount() string {
	var iban = paymentService.ibanGenerator.Generate()
	var accountInstance = account.NewAccount(iban, true, true)
	err := paymentService.accounts.Add(&accountInstance)
	if err != nil {
		panic(err)
	}
	return accountInstance.IBAN
}

// IssueMoney issues the specified amount to the "emission" account.
func (paymentService *PaymentService) IssueMoney(amount int) {
	err := paymentService.emissionAccount.Deposit(amount, true)
	if err != nil {
		panic("Failed to issue to 'emission' account")
	}
}

// TransferToDestruction sends a specified amount of money from a specified account to the "destruction" account.
func (paymentService *PaymentService) TransferToDestruction(fromIBAN string, amount int) error {
	fromAccount, err := paymentService.accounts.GetById(fromIBAN, false)
	if err != nil {
		return fmt.Errorf("account with IBAN %s not found", fromIBAN)
	}
	err = (*fromAccount).TransferTo(paymentService.destructionAccount, amount)
	if err != nil {
		return err
	}
	return nil
}

// Transfer sends a specified amount of money from account with id fromIBAN to account with id toIBAN.
func (paymentService *PaymentService) Transfer(fromIBAN string, toIBAN string, amount int) error {
	fromAccount, err := paymentService.accounts.GetById(fromIBAN, false)
	if err != nil {
		return fmt.Errorf("account with IBAN %s not found", fromIBAN)
	}
	toAccount, err := paymentService.accounts.GetById(toIBAN, false)
	if err != nil {
		return fmt.Errorf("account with IBAN %s not found", toIBAN)
	}
	err = (*fromAccount).TransferTo(*toAccount, amount)
	if err != nil {
		return err
	}
	return nil
}

func (paymentService *PaymentService) DumpAccountsAsJSON() string {
	accounts := paymentService.accounts.GelAll()
	jsonData, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

func (paymentService *PaymentService) TransferByJson(request TransferRequest) error {
	return paymentService.Transfer(request.From, request.To, request.Amount)
}

func (paymentService *PaymentService) TransferByJsonString(requestString string) error {
	var request TransferRequest
	err := json.Unmarshal([]byte(requestString), &request)
	if err != nil {
		return err
	}
	return paymentService.TransferByJson(request)
}

func (paymentService *PaymentService) GetBalance(IBAN string) (int, error) {
	accountInstance, err := paymentService.accounts.GetById(IBAN, false)
	if err != nil {
		return 0, err
	}
	return (*accountInstance).GetBalance(), nil
}

// GetAccountInfo returns account as json
func (paymentService *PaymentService) GetAccountInfo(IBAN string) (string, error) {
	accountInstance, err := paymentService.accounts.GetById(IBAN, false)
	if err != nil {
		return "", err
	}
	jsonAccount := (*accountInstance).ToJson()
	jsonData, err := json.MarshalIndent(jsonAccount, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (paymentService *PaymentService) BlockAccount(IBAN string) error {
	accountInstance, err := paymentService.accounts.GetById(IBAN, false)
	if err != nil {
		return err
	}
	return (*accountInstance).Block()
}

func (paymentService *PaymentService) UnblockAccount(IBAN string) error {
	accountInstance, err := paymentService.accounts.GetById(IBAN, false)
	if err != nil {
		return err
	}
	(*accountInstance).Unblock()
	return nil
}
