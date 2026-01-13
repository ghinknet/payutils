package utils

import "fmt"

// CentsToYuan transfer cents to yuan
func CentsToYuan(cents int64) string {
	yuan := cents / 100
	remainder := cents % 100

	if cents < 0 {
		yuan = -yuan
		remainder = -remainder
		return fmt.Sprintf("-%d.%02d", yuan, remainder)
	}

	return fmt.Sprintf("%d.%02d", yuan, remainder)
}
