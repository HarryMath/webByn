package util

import "sync"

type IncrementalGenerator struct {
	mu      sync.Mutex
	current int
}

func NewIncrementalGenerator() *IncrementalGenerator {
	return &IncrementalGenerator{}
}

func (ng *IncrementalGenerator) Next() int {
	ng.mu.Lock()
	defer ng.mu.Unlock()
	ng.current++
	return ng.current
}
