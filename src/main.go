package main

import (
	"strconv"
	"webByn/src/account"
	"webByn/src/bank"
)

func main() {
	var bynSystem = bank.WebByn{}
	_, err1 := bynSystem.CreateNewAccount()
	if err1 != nil {
		panic(err1)
	}
	var acc = account.NewAccount("IBAN")
	var stringBalance = strconv.Itoa(acc.GetBalance())
	println("Hello world! First account with balance " + stringBalance)
}
