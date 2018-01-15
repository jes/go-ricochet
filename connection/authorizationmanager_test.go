package connection

import (
	"testing"
)

func TestAuthorizationManager(t *testing.T) {
	am := new(AuthorizationManager)
	am.Init()
	am.AddAuthorization("test")
	if am.Authorized("test") != nil {
		t.Errorf("Authorized(test) should return nil, instead returned error: %v", am.Authorized("test"))
	}

	if am.Authorized("not_authed") == nil {
		t.Errorf("Authorized(not_authed) should return error, instead returned nil: %v", am.Authorized("not_authed"))
	}
}
