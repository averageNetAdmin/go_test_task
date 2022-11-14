package claim

import "github.com/golang-jwt/jwt"

type UserClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
}
