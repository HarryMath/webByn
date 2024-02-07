package util

import (
	"fmt"
)

type IBANGenerator struct {
	countryCode string
	codeLength  int
	codeGen     *IncrementalGenerator
}

func NewIBANGenerator(countryCode string, codeLength int) *IBANGenerator {
	return &IBANGenerator{
		countryCode,
		codeLength,
		NewIncrementalGenerator(),
	}
}

func (ibanGenerator *IBANGenerator) Generate() string {
	code := ibanGenerator.codeGen.Next()
	paddedCode := PadLeft(code, ibanGenerator.codeLength)
	return fmt.Sprintf("%s%s", ibanGenerator.countryCode, paddedCode)
}
