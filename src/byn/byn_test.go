package byn

import (
	"sync"
	"testing"
)

// TestSingleton gets BYN PaymentSystem instance multiple times in parallel threads and issue money
// expecting that all money were issued to one account doe to Payment System is singleton
func TestSingleton(t *testing.T) {
	var system = GetBynSystem()

	var wg sync.WaitGroup
	const routinesCount = 100
	const iterationsPerThread = 1000
	const issueAmount = 765

	wg.Add(routinesCount)
	for i := 0; i < routinesCount; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterationsPerThread; j++ {
				GetBynSystem().IssueMoney(issueAmount)
			}
		}()
	}
	wg.Wait()

	expectedBalance := routinesCount * iterationsPerThread * issueAmount
	actualBalance := system.emissionAccount.GetBalance()
	if expectedBalance != actualBalance {
		t.Errorf("Expected balance: %d, got: %d", expectedBalance, actualBalance)
	}
}
