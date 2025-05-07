package handlers

import (
	"context"
	"huma-app/store"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type DeleteUserInput struct {
	ID uuid.UUID `path:"id"`
}

type GetAllUsersOutput struct {
	Body []store.GetUsersRow
}

func (rs *ApiHandlers) RegisterGetAll(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-all-users",
		Summary:     "Get all users",
		Method:      http.MethodGet,
		Path:        "/api/users",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"Bearer": {"admin"}},
		},
	}, func(ctx context.Context, input *struct{}) (*GetAllUsersOutput, error) {
		rows, _ := rs.repo.GetUsers(ctx)
		return &GetAllUsersOutput{Body: rows}, nil
	})
}

func (rs *ApiHandlers) RegisterDelete(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Summary:     "Delete user",
		Method:      http.MethodDelete,
		Path:        "/api/users/{id}",
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
