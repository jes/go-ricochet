package identity

import (
	"crypto"
	"crypto/rsa"
	"encoding/asn1"
	"github.com/s-rah/go-ricochet/utils"
)

// Identity is an encapsulation of Name, PrivateKey and other features
// that make up a Ricochet client.
// The purpose of Identity is to prevent other classes directly accessing private key
// and to ensure the integrity of security-critical functions.
type Identity struct {
	Name string
	pk   *rsa.PrivateKey
}

// Init loads an identity from a file. Currently file should be a private_key
// but this may change in the future. //XXX
func Init(filename string) Identity {
	pk, err := utils.LoadPrivateKeyFromFile(filename)
	if err == nil {
		return Identity{"", pk}
	}
	return Identity{}
}

// Initialize is a courtesy function for initializing an Identity in-code.
func Initialize(name string, pk *rsa.PrivateKey) Identity {
	return Identity{name, pk}
}

// Initialized ensures that an Identity has been assigned a private_key and
// is ready to perform operations.
func (i *Identity) Initialized() bool {
	if i.pk == nil {
		return false
	}
	return true
}

// PublicKeyBytes returns the public key associated with this Identity in serializable-friendly
// format. //TODO Not sure I like this.
func (i *Identity) PublicKeyBytes() []byte {
	publicKeyBytes, _ := asn1.Marshal(rsa.PublicKey{
		N: i.pk.PublicKey.N,
		E: i.pk.PublicKey.E,
	})

	return publicKeyBytes
}

// Hostname provides the onion address associated with this Identity.
func (i *Identity) Hostname() string {
	return utils.GetTorHostname(i.PublicKeyBytes())
}

// Sign produces a cryptographic signature using this Identities private key.
func (i *Identity) Sign(challenge []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(nil, i.pk, crypto.SHA256, challenge)
}
