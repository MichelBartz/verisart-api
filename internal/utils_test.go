package internal

import "testing"

func TestGenerateIDFromString(t *testing.T) {
	str := "A dummy string"

	hashed := GenerateIDFromString(str)

	if len(hashed) != 32 {
		t.Errorf("Expected a 32 length hash, got %s of length %d", hashed, len(hashed))
	}
}
