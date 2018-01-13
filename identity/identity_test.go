package identity

import (
	"testing"
)

func TestIdentity(t *testing.T) {
	id := Init("../testing/private_key")
	if id.Initialized() == false {
		t.Errorf("Identity should be initialized")
	}

	if id.Hostname() != "kwke2hntvyfqm7dr" {
		t.Errorf("Expected %v as Hostname() got: %v", "kwke2hntvyfqm7dr", id.Hostname())
	}
}
