package bank

import (
	"webByn/src/account"
	"webByn/src/repository"
)

type WebByn struct {
	accounts repository.Repository[*account.Account]
}

func (byn *WebByn) CreateNewAccount() (*account.Account, error) {
	var iban = "IBAN"
	var accountInstance = account.NewAccount(iban)
	err := byn.accounts.Add(&accountInstance)
	if err != nil {
		return nil, err
	}
	return &accountInstance, nil
}
