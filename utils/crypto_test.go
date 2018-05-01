package utils

import (
	"math"
	"testing"
)

const (
	privateKeyFile = "./../testing/private_key"
)

func TestGeneratePrivateKey(t *testing.T) {
	_, err := GeneratePrivateKey()
	if err != nil {
		t.Errorf("Error while generating private key: %v", err)
	}
}

func TestLoadPrivateKey(t *testing.T) {
	_, err := LoadPrivateKeyFromFile(privateKeyFile)
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

func TestGetOnionAddress(t *testing.T) {
	privateKey, _ := LoadPrivateKeyFromFile(privateKeyFile)
	address, err := GetOnionAddress(privateKey)
	if err != nil {
		t.Errorf("Error generating onion address from private key: %v", err)
	}
	if address != "kwke2hntvyfqm7dr" {
		t.Errorf("Error: onion address for private key not expected value")
	}
}
