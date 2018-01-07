package identity

import (
	"github.com/s-rah/go-ricochet/identity"
	"testing"
)

func TestIdentity(t *testing.T) {
	id := identity.Init("../testing/private_key")
	if id.Hostname() != "kwke2hntvyfqm7dr" {
		t.Errorf("Expected %v as Hostname() got: %v", "kwke2hntvyfqm7dr", id.Hostname())
	}
}
