package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const (
	// InvalidPrivateKeyFileError is a library error, thrown when the given key file fials to load
	InvalidPrivateKeyFileError = Error("InvalidPrivateKeyFileError")

	// RicochetKeySize - tor onion services currently use rsa key sizes of 1024 bits
	RicochetKeySize = 1024
)

// GeneratePrivateKey generates a new private key for use
func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RicochetKeySize)
	if err != nil {
		return nil, errors.New("Could not generate key: " + err.Error())
	}
	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	return x509.ParsePKCS1PrivateKey(privateKeyDer)
}

// LoadPrivateKeyFromFile loads a private key from a file...
func LoadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParsePrivateKey(pemData)
}

// ParsePrivateKey Convert a private key string to a usable private key
func ParsePrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, InvalidPrivateKeyFileError
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// PrivateKeyToString turns a private key into storable string
func PrivateKeyToString(privateKey *rsa.PrivateKey) string {
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return string(pem.EncodeToMemory(&privateKeyBlock))
}
