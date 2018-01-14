package identity

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestIdentity(t *testing.T) {
	id := Identity{}
	if id.Initialized() != false {
		t.Errorf("Identity should not be initialized")
	}

	id = Init("../testing/private_key")
	if id.Initialized() == false {
		t.Errorf("Identity should be initialized")
	}

	if id.Hostname() != "kwke2hntvyfqm7dr" {
		t.Errorf("Expected %v as Hostname() got: %v", "kwke2hntvyfqm7dr", id.Hostname())
	}

	mac := hmac.New(sha256.New, []byte("Hello"))
	mac.Write([]byte("World"))
	hmac := mac.Sum(nil)
	bytes, err := id.Sign(hmac)
	expected := "b0a0a0562735b559e0efb5b3431f1aa31ddc90d2cff114d0dc05980351a4ddc6086d92efdded8a7c447a2bab4afc5f031755738d1b21edba72680dea0e33b62e914faa1f596d5f76ca0ee91cb06e4ebab748a222cc860437b7c7afd12ebee6d6998b52183bd9eb9d5b96ea95900245480539464fa889719925e569cac0cecbc1"
	actual := fmt.Sprintf("%x", bytes)
	if expected != actual || err != nil {
		t.Errorf("Identity sign failed, %v %v", actual, err)
	}
}
