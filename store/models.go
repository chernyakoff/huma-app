// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package store

import (
	"time"

	"github.com/google/uuid"
	"huma-app/store/types"
)

type User struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	Role      types.Role `json:"role"`
	Verified  int64      `json:"verified"`
}
