package util

import (
	"fmt"
	"strconv"
	"strings"
)

func PadLeft(number int, length int) string {
	numberStr := strconv.Itoa(number)
	if len(numberStr) >= length {
		return numberStr
	}
	padding := length - len(numberStr)
	return fmt.Sprintf("%s%s", strings.Repeat("0", padding), numberStr)
}
