package handlers

import (
	"context"
	"fmt"
	"huma-app/lib/config"
	"huma-app/lib/mail"
	"huma-app/lib/middleware"
	"huma-app/lib/security"
	"huma-app/store"
	"huma-app/store/types"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StatusOutput struct {
	Status int
}

type RegisterInputBody struct {
	Email    string `format:"email" json:"email" required:"true"`
	Password string `json:"password"`
}
type RegisterInput struct {
	Body RegisterInputBody
}

func GetHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (m *RegisterInputBody) Resolve(ctx huma.Context) []error {
	// Get request info you don't normally have access to.
	m.Password, _ = GetHashPassword(m.Password)

	return nil
}

var _ huma.Resolver = (*RegisterInputBody)(nil)

func (rs *ApiHandlers) RegisterRegister(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Summary:     "Register",
		Method:      http.MethodPost,
		Path:        "/api/auth/register",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *RegisterInput) (*StatusOutput, error) {
		_, err := rs.repo.GetUserByEmail(ctx, input.Body.Email)
		if err == nil {
			return nil, huma.Error400BadRequest("This email address is already in use")
		}
		user, err := rs.repo.CreateUser(ctx, store.CreateUserParams{
			ID:       uuid.New(),
			Email:    input.Body.Email,
			Password: input.Body.Password,
		})
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}
		token, _ := rs.security.GenerateToken(security.EmailToken, time.Hour*1, user.ID, user.Role)
		link := fmt.Sprintf(config.Get().Frontend.Urls.Verify, token)
		err = mail.SendVerifyMail(user.Email, mail.VerifyEmailParams{
			AppName: config.Get().Api.Name,
			Link:    link,
		})
		if err != nil {
			//
		}

		return &StatusOutput{Status: http.StatusOK}, nil
	})
}

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

func (rs *ApiHandlers) RegisterLogin(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Summary:     "Login",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{middleware.RateLimitMiddleware(api, rs.limiter.Login)},
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
		cookie, _ := rs.security.GenerateTokenToCookies(security.AccessToken, time.Hour*24, user.ID, user.Role)

		return &LoginOutput{
			SetCookie: *cookie,
			Status:    http.StatusNoContent,
		}, nil

	})
}

type AuthHeader struct {
	Session http.Cookie `cookie:"jwt"`
	UserId  uuid.UUID   `hidden:"true"`
	Role    types.Role  `hidden:"true"`
}

func (m *AuthHeader) Resolve(ctx huma.Context) []error {
	m.UserId, _ = ctx.Context().Value("user_id").(uuid.UUID)
	return nil
}

var _ huma.Resolver = (*AuthHeader)(nil)

type ProfileOutputBody struct {
	ID    uuid.UUID  `hidden:"true" json:"id"`
	Email string     `format:"email" json:"email" required:"true"`
	Role  types.Role `json:"role"`
}

type ProfileOutput struct {
	Body store.GetUserByIdRow
}

func (rs *ApiHandlers) RegisterMe(api huma.API) {
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
		},
	}, func(ctx context.Context, input *AuthHeader) (*ProfileOutput, error) {

		user, err := rs.repo.GetUserById(ctx, input.UserId)
		if err != nil {
			fmt.Println(err.Error())
			return nil, huma.Error404NotFound("User not found")
		}
		return &ProfileOutput{Body: user}, nil
	})
}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Status    int
}

func (rs *ApiHandlers) RegisterLogout(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Summary:     "Logout",
		Description: "This endpoint is used to log out a currently authenticated user. It invalidates the user's session or authentication token, ensuring that the token can no longer be used for accessing protected resources. This is a critical security feature to prevent unauthorized access after a user logs out.",
		Method:      http.MethodGet,
		Path:        "/api/auth/logout",
		Tags:        []string{"Auth"},
		Security: []map[string][]string{
			{"Bearer": {}},
		},
	}, func(ctx context.Context, input *AuthHeader) (*LogoutOutput, error) {

		return &LogoutOutput{
			SetCookie: *rs.security.DeleteCookie(),
			Status:    http.StatusOK,
		}, nil

	})
}

type VerifyEmailInput struct {
	Token string `query:"token"`
}

func (rs *ApiHandlers) RegisterVerify(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "verify-email",
		Summary:     "Verify email",
		Description: "This endpoint is used to verify a user's email address. When a user registers or updates their email, a verification token is sent to their email address. The user must submit this token through this endpoint to confirm the validity of their email.",
		Method:      http.MethodGet,
		Path:        "/api/auth/verify-email",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{middleware.RateLimitMiddleware(api, rs.limiter.Verify)},
		Errors: []int{
			http.StatusUnauthorized,
		},
	}, func(ctx context.Context, input *VerifyEmailInput) (*StatusOutput, error) {

		token, err := rs.security.VerifyToken(input.Token, security.EmailToken)
		if err != nil {
			return &StatusOutput{
				Status: http.StatusUnauthorized,
			}, nil
		}

		user, err := rs.repo.GetUserById(ctx, token.UserID)
		if err != nil {
			return &StatusOutput{
				Status: http.StatusUnauthorized,
			}, nil
		}

		rs.repo.VerifyUser(ctx, user.ID)
		return &StatusOutput{
			Status: http.StatusOK,
		}, nil

	})
}
