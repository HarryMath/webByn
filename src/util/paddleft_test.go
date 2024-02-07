package util

import "testing"

func TestPadLeft(t *testing.T) {
	testCases := []struct {
		number         int
		length         int
		expectedPadded string
	}{
		{123, 6, "000123"},
		{123, 3, "123"},
		{0, 5, "00000"},
		{987654321, 9, "987654321"},
		{987654321, 3, "987654321"},
	}

	for _, tc := range testCases {
		padded := PadLeft(tc.number, tc.length)
		if padded != tc.expectedPadded {
			t.Errorf("PadLeft(%d, %d) = %s; expected %s", tc.number, tc.length, padded, tc.expectedPadded)
		}
	}
}
