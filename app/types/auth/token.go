package auth

import "github.com/golang-jwt/jwt/v4"

type TokenPayload struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

type TokenClaims struct {
	TokenPayload
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  Token `json:"accessToken"`
	RefreshToken Token `json:"refreshToken"`
}

type Token struct {
	Token       *jwt.Token
	Claims      TokenClaims
	TokenString string
}
