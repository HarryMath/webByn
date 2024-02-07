package util

import (
	"sync"
	"testing"
)

func TestIncrementalNumberGenerator(t *testing.T) {
	generator := NewIncrementalGenerator()

	var wg sync.WaitGroup
	const routinesCount = 100
	const iterationsPerThread = 1000

	wg.Add(routinesCount)
	for i := 0; i < routinesCount; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterationsPerThread; j++ {
				generator.Next()
			}
		}()
	}
	wg.Wait()

	expectedValue := routinesCount * iterationsPerThread
	if generator.current != expectedValue {
		t.Errorf("Expected value: %d, got: %d", expectedValue, generator.current)
	}
}
