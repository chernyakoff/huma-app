package security

import (
	"errors"
	"fmt"
	"huma-app/lib/config"
	"huma-app/store/types"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

var JWTCookieName = "jwt"

type AppToken struct {
	jwt.RegisteredClaims
	UserID    uuid.UUID
	UserRole  types.Role
	TokenType TokenType
}

var _ jwt.Claims = &AppToken{}

type TokenType string

const (
	AccessToken   TokenType = "access"
	EmailToken    TokenType = "email"    // for verify email
	PasswordToken TokenType = "password" // for password forgot
)

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrTokenNotFound    = errors.New("token not found")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrInvalidRolesType = errors.New("invalid role type. Must be []string")
	ErrExpired          = errors.New("token is expired")
)

type Security struct {
	key             []byte
	Now             func() time.Time
	ExpiresInterval time.Duration
}

func NewSecurity() *Security {
	return &Security{
		key:             []byte(config.Get().Secret.Jwt),
		Now:             time.Now,
		ExpiresInterval: 24 * time.Hour,
	}
}

func (security Security) GenerateToken(tokenType TokenType, expiresIn time.Duration, userID uuid.UUID, userRole types.Role) (string, error) {
	claims := AppToken{
		UserID:    userID,
		UserRole:  userRole,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "my-app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(security.key)
}

func (security Security) VerifyToken(tokenString string, tokenType TokenType) (*AppToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AppToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return security.key, nil
	})
	if err != nil {
		return nil, ErrInvalidTokenType
	}

	if claims, ok := token.Claims.(*AppToken); ok && token.Valid {
		if claims.TokenType != tokenType {
			return nil, ErrInvalidTokenType
		}
		return claims, nil

	} else {
		return nil, ErrInvalidTokenType

	}
}

func (security Security) GenerateTokenToCookies(tokenType TokenType, expiresIn time.Duration, userID uuid.UUID, userRole types.Role) (*http.Cookie, error) {
	token, err := security.GenerateToken(tokenType, expiresIn, userID, userRole)
	if err != nil {
		return &http.Cookie{}, err
	}
	return &http.Cookie{
		Name:     JWTCookieName,
		Value:    token,
		Expires:  security.Now().Add(security.ExpiresInterval),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		MaxAge:   int(security.ExpiresInterval.Seconds()),
	}, nil
}

func (security Security) DeleteCookie() *http.Cookie {
	return &http.Cookie{
		Name:    JWTCookieName,
		Value:   "",
		Path:    "/",
		Expires: security.Now(),
	}
}
