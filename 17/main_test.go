package main

import (
	"fmt"
	"math/big"
	"testing"
)

// Helper function to calculate sum of digits
func sumDigits(s string) int {
	sum := 0
	for _, c := range s {
		sum += int(c - '0')
	}
	return sum
}

// Helper to add big integers as strings
func addStrings(a, b string) string {
	n1 := new(big.Int)
	n1.SetString(a, 10)
	n2 := new(big.Int)
	n2.SetString(b, 10)
	n1.Add(n1, n2)
	return n1.String()
}

// Helper to multiply big integer string by int
func mulString(a string, k int) string {
	n1 := new(big.Int)
	n1.SetString(a, 10)
	n2 := big.NewInt(int64(k))
	n1.Mul(n1, n2)
	return n1.String()
}

func TestSolve(t *testing.T) {
	testCases := []int{10, 100, 1, 9, 999}

	for _, n := range testCases {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			a, d := solve(n)

			// Verify basic constraints
			if len(a) > 200000 {
				t.Errorf("a length > 2*10^5: %d", len(a))
			}
			if len(d) > 200000 {
				t.Errorf("d length > 2*10^5: %d", len(d))
			}
			if a == "" || a[0] == '0' {
				t.Errorf("a is invalid: %s", a)
			}
			if d == "" || d[0] == '0' {
				t.Errorf("d is invalid: %s", d)
			}

			// Verify arithmetic progression property of sum of digits
			// S(a), S(a+d), S(a+2d), ...

			// Get S(a)
			s0 := sumDigits(a)

			// Get S(a+d)
			term1 := addStrings(a, d)
			s1 := sumDigits(term1)

			diff := s1 - s0

			// Verify for all k from 2 to n-1
			for k := 2; k < n; k++ {
				termK := addStrings(a, mulString(d, k))
				sK := sumDigits(termK)

				expectedSK := s0 + k*diff
				if sK != expectedSK {
					t.Errorf("AP property violated at k=%d. Got S=%d, expected %d. a=%s, d=%s", k, sK, expectedSK, a, d)
					break
				}
			}
		})
	}
}

// Test with n=10000 (max input)
func TestSolveMaxN(t *testing.T) {
	n := 10000
	a, d := solve(n)

	s0 := sumDigits(a)
	term1 := addStrings(a, d)
	s1 := sumDigits(term1)
	diff := s1 - s0

	// Check random points to save time
	points := []int{2, 10, 100, 1000, 5000, 9999}
	for _, k := range points {
		termK := addStrings(a, mulString(d, k))
		sK := sumDigits(termK)

		expectedSK := s0 + k*diff
		if sK != expectedSK {
			t.Errorf("AP property violated at k=%d. Got S=%d, expected %d", k, sK, expectedSK)
		}
	}
}
