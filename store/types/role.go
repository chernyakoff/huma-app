package types

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleEditor Role = "editor"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleEditor:
		return true
	default:
		return false
	}
}
