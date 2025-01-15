package controller

import (
	"context"
	"huma-app/lib/security"
	"huma-app/store"
	"net/http"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
)

type UserResource struct {
	basePath string
	repo     UserRepo
	security *security.Security
}

func (rs *UserResource) path(segments ...string) string {
	segment := ""
	if len(segments) > 0 {
		segment = segments[0]
	}
	return filepath.Clean(filepath.Join(rs.basePath, segment))

}

func NewUserResource(userRepo UserRepo, security *security.Security) *UserResource {
	return &UserResource{
		basePath: "/api/users",
		repo:     userRepo,
		security: security,
	}
}

type CreateUserInput struct {
	Body store.CreateUserParams
}
type CreateUserOutput struct {
	Body store.CreateUserRow
}

type DeleteUserInput struct {
	ID int64 `path:"id"`
}

type GetAllUsersOutput struct {
	Body []store.GetUsersRow
}

func (rs *UserResource) RegisterGetAll(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-all-users",
		Summary:     "Get all users",
		Method:      http.MethodGet,
		Path:        rs.path(),
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"Bearer": {"admin"}},
		},
	}, func(ctx context.Context, input *struct{}) (*GetAllUsersOutput, error) {
		rows, _ := rs.repo.GetUsers(ctx)
		return &GetAllUsersOutput{Body: rows}, nil
	})
}

func (rs *UserResource) RegisterCreate(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Summary:     "Create user",
		Method:      http.MethodPost,
		Path:        rs.path(),
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
		row, err := rs.repo.CreateUser(ctx, input.Body)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return &CreateUserOutput{Body: row}, nil
	})
}

func (rs *UserResource) RegisterDelete(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Summary:     "Delete user",
		Method:      http.MethodDelete,
		Path:        rs.path("{id}"),
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"Bearer": {"admin"}},
		},
	}, func(ctx context.Context, input *DeleteUserInput) (*struct{}, error) {
		err := rs.repo.DeleteUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, nil
	})
}

type UserRepo interface {
	CreateUser(ctx context.Context, arg store.CreateUserParams) (store.CreateUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (store.GetUserByEmailRow, error)
	GetUserById(ctx context.Context, id int64) (store.GetUserByIdRow, error)
	GetUsers(ctx context.Context) ([]store.GetUsersRow, error)
	DeleteUser(ctx context.Context, id int64) error
}

var _ UserRepo = (*store.Queries)(nil)
