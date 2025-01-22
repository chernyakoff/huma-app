package security

import (
	"errors"
	"fmt"
	"huma-app/lib/config"
	"huma-app/store/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppToken struct {
	jwt.RegisteredClaims // Required, this struct contains the standard claims
	UserID               int64
	UserRole             types.Role
}

var _ jwt.Claims = &AppToken{}

var JWTCookieName = "jwt"

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

func (security Security) GenerateToken(userID int64, role types.Role) (tokenString string, err error) {
	claims := AppToken{
		UserID:   userID,
		UserRole: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен истекает через 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "my-app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(security.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateTokenToCookies generates a JWT token with the given claims and writes it to the cookies.
func (security Security) GenerateTokenToCookies(userID int64, role types.Role) (*http.Cookie, error) {
	token, err := security.GenerateToken(userID, role)
	if err != nil {
		return &http.Cookie{}, err
	}
	return &http.Cookie{
		Name:     JWTCookieName,
		Value:    token,
		Expires:  security.Now().Add(security.ExpiresInterval),
		HttpOnly: true,
		Path:     "/",
		// SameSite: http.SameSiteStrictMode,
		// Secure:   true,
		MaxAge: int(security.ExpiresInterval.Seconds()),
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

func (security Security) ParseToken(tokenString string) (*AppToken, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AppToken{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return security.key, nil
	})
	if err != nil {
		return nil, ErrInvalidTokenType
	}

	if claims, ok := token.Claims.(*AppToken); ok && token.Valid {
		return claims, nil

	} else {
		return nil, ErrInvalidTokenType

	}

}
