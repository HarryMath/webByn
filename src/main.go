package main

import (
	"strconv"
	"webByn/src/byn"
)

func main() {
	Scenario1()
	Scenario2()
	Scenario3()
	Scenario4()
}

func printDestructionBalance() {
	bynSystem := byn.GetBynSystem()
	destructionIBAN := bynSystem.GetDestructionAccountNumber()
	balance, err := bynSystem.GetBalance(destructionIBAN)
	if err != nil {
		panic(err)
	}
	println("Destruction account balance: " + strconv.Itoa(balance))
}

func printEmissionBalance() {
	bynSystem := byn.GetBynSystem()
	emissionIBAN := bynSystem.GetEmissionAccountNumber()
	balance, err := bynSystem.GetBalance(emissionIBAN)
	if err != nil {
		panic(err)
	}
	println("Emission account balance: " + strconv.Itoa(balance))
}

// Scenario1 demonstrates how to issue money, create account and transfer money
func Scenario1() {
	println("\n------- SCENARIO 1 STARTED --------")
	var bynSystem = byn.GetBynSystem()

	println("\nissue 1000 money into circulation (выпуск 1000 денег в оборот)")
	bynSystem.IssueMoney(1000)

	println("\ncreate new account (создание нового аккаунта)")
	accountIBAN := bynSystem.OpenAccount()
	println("created account: " + accountIBAN)

	println("\nget emission IBAN (получние счета эмиссии)")
	emissionIBAN := bynSystem.GetEmissionAccountNumber()
	println("emission IBAN: " + emissionIBAN)

	println("\ntransfer 400 money to new account (перевод 1000 денег на новый аккаунт)")
	_ = bynSystem.Transfer(emissionIBAN, accountIBAN, 400)

	println("\nget account balance (получение баланса после перевода)")
	balance, _ := bynSystem.GetBalance(accountIBAN)
	println("balance is " + strconv.Itoa(balance))

	println("\ntransfer 600 money to account using json (перевод 600 денег на новый аккаунт)")
	jsonRequest := `{ "fromIBAN": "BY00000000000000000000000001", "toIBAN": "BY00000000000000000000000003", "amount": 600 }`
	println(jsonRequest)
	err := bynSystem.TransferByJsonString(jsonRequest)
	if err != nil {
		panic(err)
	}
	balance, _ = bynSystem.GetBalance(accountIBAN)
	println("balance is " + strconv.Itoa(balance))

	println("\n------- SCENARIO 1 FINISHED --------")
}

// Scenario2 demonstrates blocking and unblocking accounts
func Scenario2() {
	println("\n------- SCENARIO 2 STARTED --------")
	var bynSystem = byn.GetBynSystem()
	emissionIBAN := bynSystem.GetEmissionAccountNumber()

	println("\nissue 500 money into circulation (выпуск 500 денег в оборот)")
	bynSystem.IssueMoney(500)

	println("\ncreate new account (создание нового аккаунта)")
	accountIBAN := bynSystem.OpenAccount()
	println("created account: " + accountIBAN)

	println("\nblock created account (блокировка аккаунта)")
	_ = bynSystem.BlockAccount(accountIBAN)
	accountInfo, _ := bynSystem.GetAccountInfo(accountIBAN)
	println("check that status is 'blocked': \n" + accountInfo)

	println("\ntry to transfer money to blocked account (попытка перевода денег на заблокированный аккаунт)")
	err := bynSystem.Transfer(emissionIBAN, accountIBAN, 500)
	if err != nil {
		println("Error: ", err.Error())
		println("\ncheck that balance on emission account is 500 (проверка что баланс на аккаунте эмиссии не изменился)")
		println("expect 500")
		printEmissionBalance()
	}

	println("\nunblock account and try transfer again (разблокировка аккаунт и перевод 500 денег)")
	_ = bynSystem.UnblockAccount(accountIBAN)
	_ = bynSystem.Transfer(emissionIBAN, accountIBAN, 500)
	balance, _ := bynSystem.GetBalance(accountIBAN)
	println("balance after transfer is " + strconv.Itoa(balance))

	println("\ntry to block emission account")
	err = bynSystem.BlockAccount(emissionIBAN)
	if err != nil {
		println("Error: ", err.Error())
	}

	println("\n------- SCENARIO 2 FINISHED --------")
}

// Scenario3 demonstrates what happens if not enough money to transfer
func Scenario3() {
	println("\n------- SCENARIO 3 STARTED --------")
	var bynSystem = byn.GetBynSystem()
	emissionIBAN := bynSystem.GetEmissionAccountNumber()

	println("\nissue 1000 money into circulation (выпуск 1000 денег в оборот)")
	bynSystem.IssueMoney(1000)

	println("\ncreate two new accounts (создание двух новых аккаунтов)")
	account1 := bynSystem.OpenAccount()
	account2 := bynSystem.OpenAccount()
	println("created account 1: " + account1)
	println("created account 2: " + account2)

	println("\ngive money to both accounts (перевод средств на оба аккаунта)")
	_ = bynSystem.Transfer(emissionIBAN, account1, 300)
	_ = bynSystem.Transfer(emissionIBAN, account2, 700)
	balance, _ := bynSystem.GetBalance(account1)
	println("balance for account 1 is " + strconv.Itoa(balance))
	balance, _ = bynSystem.GetBalance(account2)
	println("balance for account 2 is " + strconv.Itoa(balance))

	println("\ntransfer 1500 money from account 1 to account 2 (перевод 1500 денег с аккаунта 1 на аккаунт 2)")
	err := bynSystem.Transfer(account1, account2, 1500)
	if err != nil {
		println("Error: ", err.Error())
	}

	println("\nCheck that balance not changed (проверка что балнс не изменился)")
	balance, _ = bynSystem.GetBalance(account1)
	println("balance for account 1 is " + strconv.Itoa(balance))
	balance, _ = bynSystem.GetBalance(account2)
	println("balance for account 2 is " + strconv.Itoa(balance))

	println("\n------- SCENARIO 3 FINISHED --------")
}

// Scenario4 sends money to destruction account prints all accounts as json
func Scenario4() {
	println("\n------- SCENARIO 4 STARTED --------")
	var bynSystem = byn.GetBynSystem()
	emissionIBAN := bynSystem.GetEmissionAccountNumber()
	destructionIBAN := bynSystem.GetDestructionAccountNumber()

	println("\nissue 300 money into circulation (выпуск 300 денег в оборот)")
	bynSystem.IssueMoney(300)

	println("\ncreate new account (создание новго аккаунта)")
	accountIBAN := bynSystem.OpenAccount()
	println("created account: " + accountIBAN)

	println("\ngive money to account (перевод средств на аккаунт)")
	_ = bynSystem.Transfer(emissionIBAN, accountIBAN, 300)
	balance, _ := bynSystem.GetBalance(accountIBAN)
	println("balance for account is " + strconv.Itoa(balance))

	println("\ntransfer 100 money to destruction account (вывод 100 едениц денег из оборота)")
	_ = bynSystem.TransferToDestruction(accountIBAN, 100)
	balance, _ = bynSystem.GetBalance(accountIBAN)
	println("balance for account is " + strconv.Itoa(balance))
	printDestructionBalance()

	println("\ntry to transfer money back from destruction account (попытка вернуть деньги обратно со счета 'уничтожения')")
	err := bynSystem.Transfer(destructionIBAN, accountIBAN, 100)
	if err != nil {
		println("Error: ", err.Error())
	}

	println("\nblocking created account (блокировка созданного аккаунта)")
	_ = bynSystem.BlockAccount(accountIBAN)

	println("\nprint all accounts as json (вывод всех аккаунтов в формате json)")
	result := bynSystem.DumpAccountsAsJSON()
	println(result)

	println("\n------- SCENARIO 4 FINISHED --------")
}
