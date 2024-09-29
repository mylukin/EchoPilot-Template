package config

import (
	"github.com/golang-jwt/jwt/v5"
)

// UserLoginClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type UserLoginClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
