package util

import (
	"sync"
	"testing"
)

func TestIBANGenerator(t *testing.T) {
	var BY = "BY"
	var RU = "RU"
	generatorBY := NewIBANGenerator(BY, 26)
	generatorRU := NewIBANGenerator(RU, 28)

	var wg sync.WaitGroup
	const routinesCount = 101
	const iterationsPerThread = 999

	wg.Add(routinesCount * 2)
	for i := 0; i < routinesCount; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterationsPerThread; j++ {
				generatorBY.Generate()
			}
		}()
		go func() {
			defer wg.Done()
			for k := 0; k < iterationsPerThread; k++ {
				generatorRU.Generate()
			}
		}()
	}
	wg.Wait()

	var iterations = routinesCount * iterationsPerThread
	var expectedBY = BY + PadLeft(iterations+1, 26)
	var expectedRU = RU + PadLeft(iterations+1, 28)

	var actualBY = generatorBY.Generate()
	var actualRU = generatorRU.Generate()

	if actualBY != expectedBY {
		t.Errorf("Expected BY IBAN: expected %s, got: %s", expectedBY, actualBY)
	}
	if actualRU != expectedRU {
		t.Errorf("Expected RU IBAN: expected %s, got: %s", expectedRU, actualRU)
	}
}
