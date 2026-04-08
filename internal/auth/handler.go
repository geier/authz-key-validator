package auth

import (
	"context"

	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

// AuthorizationServer implements the Envoy ext_authz API
type AuthorizationServer struct {
	envoyauth.UnimplementedAuthorizationServer
}

// Check performs the external authorization check.
func (s *AuthorizationServer) Check(ctx context.Context, req *envoyauth.CheckRequest) (*envoyauth.CheckResponse, error) {
	return nil, errors.New("not implemented")
}
