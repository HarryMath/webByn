package account

type Account struct {
	IBAN    string
	balance int
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
