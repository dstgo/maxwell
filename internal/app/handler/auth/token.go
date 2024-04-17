package auth

import (
	"context"
	"github.com/dstgo/maxwell/internal/app/conf"
	"github.com/dstgo/maxwell/internal/app/data/cache"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/jwtx"
	"github.com/ginx-contribs/str2bytes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenPayload struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
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

func NewTokenHandler(jwtConf conf.JwtConf, client *redis.Client) *TokenHandler {
	return &TokenHandler{
		method:       jwt.SigningMethodHS256,
		accessCache:  cache.NewRedisTokenCache("access", client),
		refreshCache: cache.NewRedisTokenCache("refresh", client),
		jwtConf:      jwtConf,
	}
}

// TokenHandler is responsible for maintaining authentication tokens
type TokenHandler struct {
	method       jwt.SigningMethod
	accessCache  cache.TokenCache
	refreshCache cache.TokenCache
	jwtConf      conf.JwtConf
}

func (t *TokenHandler) Issue(ctx context.Context, payload TokenPayload, refresh bool) (TokenPair, error) {
	now := time.Now()
	var tokenPair TokenPair

	// issue access token
	accessToken, err := t.newToken(now, t.jwtConf.Access.Key, payload)
	if err != nil {
		return tokenPair, err
	}

	// store into the cache
	if err := t.accessCache.Set(ctx, accessToken.Claims.ID, accessToken.Claims.ID, t.jwtConf.Access.Expire); err != nil {
		return TokenPair{}, err
	}

	tokenPair.AccessToken = accessToken
	// no need to refresh the token
	if !refresh {
		return tokenPair, nil
	}

	// issue refresh token
	refreshToken, err := t.newToken(now, t.jwtConf.Refresh.Key, payload)
	if err != nil {
		return tokenPair, err
	}

	// associated with access token
	if err := t.accessCache.Set(ctx, refreshToken.Claims.ID, accessToken.Claims.ID, t.jwtConf.Refresh.Expire); err != nil {
		return tokenPair, nil
	}
	tokenPair.RefreshToken = refreshToken

	return tokenPair, nil
}

// Refresh refreshes the access token lifetime with the given refresh token
func (t *TokenHandler) Refresh(ctx context.Context, accessToken string, refreshToken string) (TokenPair, error) {
	now := time.Now()
	var pair TokenPair
	// return directly if refresh token is expired
	refresh, err := t.VerifyRefresh(ctx, refreshToken)
	if err != nil {
		return pair, err
	}
	pair.RefreshToken = refresh

	// parse access token
	access, err := t.VerifyAccess(ctx, accessToken)
	if errors.Is(err, jwt.ErrTokenExpired) {
		// return if over the delay time
		if access.Claims.ExpiresAt.Add(t.jwtConf.Access.Delay).Sub(now) < 0 {
			return pair, jwt.ErrTokenExpired
		}
	} else if err != nil {
		return pair, err
	}

	// check access token if is associated with refresh token
	id, err := t.refreshCache.Get(ctx, refresh.Claims.ID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return pair, err
	}
	if access.Claims.ID != id {
		return pair, jwt.ErrTokenUnverifiable
	}

	// use a new token to replace the old one
	newAccess, err := t.newToken(now, t.jwtConf.Access.Key, access.Claims.TokenPayload)
	if err != nil {
		return pair, err
	}
	pair.AccessToken = newAccess

	// get rest ttl
	ttl, err := t.accessCache.TTL(ctx, access.Claims.ID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return pair, statuserr.InternalError(err)
	}
	// extend lifetime of access token
	ttl += t.jwtConf.Access.Expire
	if err := t.accessCache.Set(ctx, newAccess.Claims.ID, newAccess.Claims.ID, ttl); err != nil {
		return pair, statuserr.InternalError(err)
	}

	// update association
	if err := t.refreshCache.Set(ctx, refresh.Claims.ID, newAccess.Claims.ID, -1); err != nil {
		return pair, statuserr.InternalError(err)
	}

	return pair, nil
}

func (t *TokenHandler) VerifyAccess(ctx context.Context, token string) (Token, error) {
	return t.verify(ctx, t.jwtConf.Access.Key, token)
}

func (t *TokenHandler) VerifyRefresh(ctx context.Context, token string) (Token, error) {
	return t.verify(ctx, t.jwtConf.Refresh.Key, token)
}

func (t *TokenHandler) newToken(now time.Time, key string, payload TokenPayload) (Token, error) {
	// create the token claims
	claims := TokenClaims{
		TokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.jwtConf.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.jwtConf.Access.Expire)),
			ID:        uuid.NewString(),
		},
	}

	// issue the token
	token, err := jwtx.IssueWithClaims(str2bytes.Str2Bytes(key), t.method, claims)
	if err != nil {
		return Token{}, err
	}

	return Token{
		Token:       token.Token,
		Claims:      claims,
		TokenString: token.SignedString,
	}, err
}

func (t *TokenHandler) parse(token, secret string) (Token, error) {
	parseJwt, err := jwtx.VerifyWithClaims(token, str2bytes.Str2Bytes(secret), t.method, &TokenClaims{})
	if err == nil || errors.Is(err, jwt.ErrTokenExpired) {
		return Token{
			Token:       parseJwt.Token,
			Claims:      *parseJwt.Claims.(*TokenClaims),
			TokenString: parseJwt.SignedString,
		}, nil
	} else {
		return Token{}, err
	}
}

func (t *TokenHandler) verify(ctx context.Context, key, token string) (Token, error) {
	parsedToken, err := t.parse(token, key)
	if err != nil {
		return parsedToken, err
	}
	// check if exists in cache
	if _, err := t.accessCache.Get(ctx, parsedToken.Claims.ID); errors.Is(err, redis.Nil) {
		return parsedToken, jwt.ErrTokenExpired
	} else if err != nil {
		return parsedToken, statuserr.InternalError(err)
	}
	return parsedToken, nil
}
