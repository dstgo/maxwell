package auth

import (
	"context"
	"github.com/dstgo/maxwell/server/conf"
	"github.com/dstgo/maxwell/server/data/cache"
	"github.com/dstgo/maxwell/server/types/auth"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/jwtx"
	"github.com/ginx-contribs/str2bytes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewTokenHandler(jwtConf conf.Jwt, client *redis.Client) *TokenHandler {
	return &TokenHandler{
		method:       jwt.SigningMethodHS256,
		accessCache:  cache.NewRedisTokenCache("access", client),
		refreshCache: cache.NewRedisTokenCache("refresh", client),
		JwtConf:      jwtConf,
	}
}

// TokenHandler is responsible for maintaining authentication tokens
type TokenHandler struct {
	method       jwt.SigningMethod
	accessCache  cache.TokenCache
	refreshCache cache.TokenCache
	JwtConf      conf.Jwt
}

func (t *TokenHandler) Issue(ctx context.Context, payload auth.TokenPayload, refresh bool) (auth.TokenPair, error) {
	now := time.Now()
	var tokenPair auth.TokenPair

	// issue access token
	accessToken, err := t.newToken(now, t.JwtConf.Access.Key, payload)
	if err != nil {
		return tokenPair, err
	}

	// consider network latency
	latency := time.Second * 10

	ttl := t.JwtConf.Access.Expire + t.JwtConf.Access.Delay + latency
	// store into the cache
	if err := t.accessCache.Set(ctx, accessToken.Claims.ID, accessToken.Claims.ID, ttl); err != nil {
		return auth.TokenPair{}, err
	}

	tokenPair.AccessToken = accessToken
	// no need to refresh the token
	if !refresh {
		return tokenPair, nil
	}

	// issue refresh token
	refreshToken, err := t.newToken(now, t.JwtConf.Refresh.Key, payload)
	if err != nil {
		return tokenPair, err
	}

	// associated with access token
	if err := t.refreshCache.Set(ctx, refreshToken.Claims.ID, accessToken.Claims.ID, t.JwtConf.Refresh.Expire); err != nil {
		return tokenPair, nil
	}
	tokenPair.RefreshToken = refreshToken

	return tokenPair, nil
}

// Refresh refreshes the access token lifetime with the given refresh token
func (t *TokenHandler) Refresh(ctx context.Context, accessToken string, refreshToken string) (auth.TokenPair, error) {
	now := time.Now()
	var pair auth.TokenPair
	// return directly if refresh token is expired
	refresh, err := t.VerifyRefresh(ctx, refreshToken)
	if err != nil {
		return pair, err
	}
	pair.RefreshToken = refresh

	// parse access token
	access, err := t.VerifyAccess(ctx, accessToken, now)
	if errors.Is(err, jwt.ErrTokenExpired) {
		// return if over the delay time
		if access.Claims.ExpiresAt.Add(t.JwtConf.Access.Delay).Sub(now) < 0 {
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
	newAccess, err := t.newToken(now, t.JwtConf.Access.Key, access.Claims.TokenPayload)
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
	ttl += t.JwtConf.Access.Expire
	if err := t.accessCache.Set(ctx, newAccess.Claims.ID, newAccess.Claims.ID, ttl); err != nil {
		return pair, statuserr.InternalError(err)
	}

	// update association
	if err := t.refreshCache.Set(ctx, refresh.Claims.ID, newAccess.Claims.ID, -1); err != nil {
		return pair, statuserr.InternalError(err)
	}

	return pair, nil
}

// VerifyAccess verifies the access token if is valid and parses the payload in the token.
func (t *TokenHandler) VerifyAccess(ctx context.Context, token string, now time.Time) (auth.Token, error) {
	parsedToken, err := t.parse(token, t.JwtConf.Access.Key)
	if errors.Is(err, jwt.ErrTokenExpired) {
		// check if token needs to be refreshed
		if parsedToken.Claims.Remember && parsedToken.Claims.ExpiresAt.Add(t.JwtConf.Access.Delay).Sub(now) > 0 {
			return parsedToken, auth.ErrTokenNeedsRefresh
		}
		return parsedToken, err
	} else if err != nil {
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

// VerifyRefresh verifies the refresh token if is valid.
func (t *TokenHandler) VerifyRefresh(ctx context.Context, token string) (auth.Token, error) {
	parsedToken, err := t.parse(token, t.JwtConf.Refresh.Key)
	if err != nil {
		return parsedToken, err
	}
	// check if exists in cache
	if _, err := t.refreshCache.Get(ctx, parsedToken.Claims.ID); errors.Is(err, redis.Nil) {
		return parsedToken, jwt.ErrTokenExpired
	} else if err != nil {
		return parsedToken, statuserr.InternalError(err)
	}
	return parsedToken, nil
}

func (t *TokenHandler) newToken(now time.Time, key string, payload auth.TokenPayload) (auth.Token, error) {
	// create the token claims
	claims := auth.TokenClaims{
		TokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.JwtConf.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.JwtConf.Access.Expire)),
			ID:        uuid.NewString(),
		},
	}

	// issue the token
	token, err := jwtx.IssueWithClaims(str2bytes.Str2Bytes(key), t.method, claims)
	if err != nil {
		return auth.Token{}, err
	}

	return auth.Token{
		Token:       token.Token,
		Claims:      claims,
		TokenString: token.SignedString,
	}, err
}

func (t *TokenHandler) parse(token, secret string) (auth.Token, error) {
	parseJwt, err := jwtx.VerifyWithClaims(token, str2bytes.Str2Bytes(secret), t.method, &auth.TokenClaims{})
	if err == nil || errors.Is(err, jwt.ErrTokenExpired) {
		return auth.Token{
			Token:       parseJwt.Token,
			Claims:      *parseJwt.Claims.(*auth.TokenClaims),
			TokenString: parseJwt.SignedString,
		}, nil
	} else {
		return auth.Token{}, err
	}
}
