package auth

import "github.com/golang-jwt/jwt/v5"

type TokenPayload struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
	Remember bool   `json:"remember"`
}

// TokenClaims is payload info in jwt
type TokenClaims struct {
	TokenPayload
	jwt.RegisteredClaims
}

// Token represents a jwt token
type Token struct {
	Token       *jwt.Token
	Claims      TokenClaims
	TokenString string
}

// TokenPair represents a jwt token pair composed of access token and refresh token
type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
}

type TokenResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

const (
	UsageUnknown  Usage = 0
	UsageRegister Usage = 1
	UsageReset    Usage = 2
)

type Usage int

func (u Usage) Name() string {
	switch u {
	case 1:
		return "register"
	case 2:
		return "reset"
	default:
		return "unknown"
	}
}

func (u Usage) String() string {
	switch u {
	case 1:
		return "register account"
	case 2:
		return "reset password"
	default:
		return "unknown usage"
	}
}

func CheckValidUsage(u Usage) error {
	if u.String() == UsageUnknown.String() {
		return ErrVerifyCodeUsageUnsupported
	}
	return nil
}
