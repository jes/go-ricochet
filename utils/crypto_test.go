package utils

import (
	"math"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	_, err := GeneratePrivateKey()
	if err != nil {
		t.Errorf("Error while generating private key: %v", err)
	}
}

func TestLoadPrivateKey(t *testing.T) {
	_, err := LoadPrivateKeyFromFile("../testing/private_key")
	if err != nil {
		t.Errorf("Error while loading private key from file: %v", err)
	}
}

func TestGetRandNumber(t *testing.T) {
	num := GetRandNumber()
	if !num.IsUint64() || num.Uint64() > uint64(math.MaxUint32) {
		t.Errorf("Error random number outside of expected bounds %v", num)
	}
}
