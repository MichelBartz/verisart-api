package internal

import (
	"crypto/md5"
	"fmt"
	"io"
)

// GenerateIDFromString takes an extremely naive approach to generate a loose UUID from a string
func GenerateIDFromString(str string) string {
	h := md5.New()
	io.WriteString(h, str)

	return fmt.Sprintf("%x", h.Sum(nil))
}
