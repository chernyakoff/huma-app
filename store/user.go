package store

import (
	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (m *CreateUserParams) Resolve(ctx huma.Context) []error {
	// Get request info you don't normally have access to.
	m.Password, _ = GetHashPassword(m.Password)

	return nil
}

var _ huma.Resolver = (*CreateUserParams)(nil)
