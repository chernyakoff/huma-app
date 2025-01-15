package store

import (
	"github.com/danielgtaylor/huma/v2"
)

func (m *CreateUserParams) Resolve(ctx huma.Context) []error {
	// Get request info you don't normally have access to.
	m.Password, _ = GetHashPassword(m.Password)

	return nil
}

var _ huma.Resolver = (*CreateUserParams)(nil)
