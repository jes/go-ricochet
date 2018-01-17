package connection

import (
	"github.com/s-rah/go-ricochet/utils"
)

// AuthorizationManager helps keep track of permissions for a connection
type AuthorizationManager struct {
	Authorizations map[string]bool
}

// Init sets up an AuthorizationManager to be used.
func (am *AuthorizationManager) Init() {
	am.Authorizations = make(map[string]bool)
}

// AddAuthorization adds the string authz to the map of allowed authorizations
func (am *AuthorizationManager) AddAuthorization(authz string) {
	am.Authorizations[authz] = true
}

// Authorized returns no error in the case an authz type is authorized, error otherwise.
func (am *AuthorizationManager) Authorized(authz string) error {
	_, authed := am.Authorizations[authz]
	if !authed {
		return utils.UnauthorizedActionError
	}
	return nil
}
