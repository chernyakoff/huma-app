package controller

import (
	"context"
	"fmt"
	"huma-app/store"
	"huma-app/store/types"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

type LoginInputBody struct {
	Email    string `format:"email" required:"true"`
	Password string
}

type LoginInput struct {
	Body LoginInputBody
}

type LoginOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Status    int
}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type AuthHeader struct {
	Session http.Cookie `cookie:"jwt"`
	UserId  int64       `hidden:"true"`
	Role    types.Role  `hidden:"true"`
}

type ProfileOutput struct {
	Body store.GetUserByIdRow
}

func (rs *UserResource) RegisterLogin(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Summary:     "Login",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Tags:        []string{"Auth"},
		Errors: []int{
			http.StatusUnauthorized,
		},
	}, func(ctx context.Context, input *LoginInput) (*LoginOutput, error) {

		user, err := rs.repo.GetUserByEmail(ctx, input.Body.Email)

		if err != nil {
			return nil, huma.Error401Unauthorized("Wrong password or email")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Body.Password))
		if err != nil {
			return nil, huma.Error401Unauthorized("Wrong password or email")
		}
		cookie, _ := rs.security.GenerateTokenToCookies(user.ID, user.Role)

		return &LoginOutput{
			SetCookie: *cookie,
			Status:    http.StatusNoContent,
		}, nil

	})
}

func (m *AuthHeader) Resolve(ctx huma.Context) []error {
	m.UserId, _ = ctx.Context().Value("user_id").(int64)
	m.Role, _ = ctx.Context().Value("user_role").(types.Role)
	return nil
}

var _ huma.Resolver = (*AuthHeader)(nil)

type Dummy struct{}

func (rs *UserResource) RegisterMe(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "me",
		Summary:     "Me",
		Method:      http.MethodGet,
		Path:        "/api/auth/me",
		Tags:        []string{"Auth"},
		Security: []map[string][]string{
			{"Bearer": {}},
		},
		Errors: []int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusBadRequest,
		},
	}, func(ctx context.Context, input *AuthHeader) (*ProfileOutput, error) {

		user, err := rs.repo.GetUserById(ctx, input.UserId)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return &ProfileOutput{Body: user}, nil
	})
}

func (rs *UserResource) RegisterLogout(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Summary:     "Logout",
		Method:      http.MethodGet,
		Path:        "/api/auth/logout",
		Tags:        []string{"Auth"},
		Security: []map[string][]string{
			{"Bearer": {}},
		},
	}, func(ctx context.Context, input *AuthHeader) (*LogoutOutput, error) {

		return &LogoutOutput{
			SetCookie: *rs.security.DeleteCookie(),
		}, nil

	})
}
