package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

const tokenKey = "auth.token.context.info.key"

// SetTokenInfo stores token information into context
func SetTokenInfo(ctx *gin.Context, token *Token) {
	ctx.Set(tokenKey, token)
}

// GetTokenInfo returns token information from context
func GetTokenInfo(ctx *gin.Context) (*Token, error) {
	value, exists := ctx.Get(tokenKey)
	if !exists {
		return nil, errors.New("there is no token in context")
	}

	if token, ok := value.(*Token); !ok {
		return nil, fmt.Errorf("expected %T, got %T", &Token{}, value)
	} else if token == nil {
		return nil, errors.New("nil token in context")
	} else {
		return token, nil
	}
}

// MustGetTokenInfo returns token information from context, panic if err != nil
func MustGetTokenInfo(ctx *gin.Context) *Token {
	info, err := GetTokenInfo(ctx)
	if err != nil {
		panic(err)
	}
	return info
}
