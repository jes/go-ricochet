package utils

import (
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
