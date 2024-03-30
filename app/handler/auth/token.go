package auth

import (
	"context"
	"github.com/dstgo/contrib/ginx/resp/errs"
	"github.com/dstgo/contrib/util/jwts"
	"github.com/dstgo/maxwell/app/conf"
	"github.com/dstgo/maxwell/app/data/cache"
	"github.com/dstgo/maxwell/app/types/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// TokenHandler is responsible for maintaining authentication tokens
type TokenHandler struct {
	method       jwt.SigningMethod
	accessCache  cache.TokenCache
	refreshCache cache.TokenCache
	jwtConf      conf.JwtConf
}

func (t *TokenHandler) Issue(ctx context.Context, payload auth.TokenPayload, refresh bool) (auth.TokenPair, error) {
	now := time.Now()
	var tokenPair auth.TokenPair

	// issue access token
	accessToken, err := t.newToken(now, t.jwtConf.Access.Key, payload)
	if err != nil {
		return tokenPair, err
	}

	// store into the cache
	if err := t.accessCache.Set(ctx, accessToken.Claims.ID, accessToken.Claims.ID, t.jwtConf.Access.Expire); err != nil {
		return auth.TokenPair{}, err
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
		return pair, errs.Internal(err)
	}
	// extend lifetime of access token
	ttl += t.jwtConf.Access.Expire
	if err := t.accessCache.Set(ctx, newAccess.Claims.ID, newAccess.Claims.ID, ttl); err != nil {
		return pair, errs.Internal(err)
	}

	// update association
	if err := t.refreshCache.Set(ctx, refresh.Claims.ID, newAccess.Claims.ID, -1); err != nil {
		return pair, errs.Internal(err)
	}

	return pair, nil
}

func (t *TokenHandler) VerifyAccess(ctx context.Context, token string) (auth.Token, error) {
	return t.verify(ctx, t.jwtConf.Access.Key, token)
}

func (t *TokenHandler) VerifyRefresh(ctx context.Context, token string) (auth.Token, error) {
	return t.verify(ctx, t.jwtConf.Refresh.Key, token)
}

func (t *TokenHandler) newToken(now time.Time, key string, payload auth.TokenPayload) (auth.Token, error) {
	// create the token claims
	claims := auth.TokenClaims{
		TokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.jwtConf.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.jwtConf.Access.Expire)),
			ID:        uuid.NewString(),
		},
	}

	// issue the token
	token, err := jwts.NewJwt(key, t.method, claims)
	if err != nil {
		return auth.Token{}, err
	}

	return auth.Token{
		Token:       token.Token,
		Claims:      claims,
		TokenString: token.SignedJwt,
	}, err
}

func (t *TokenHandler) parse(token, secret string) (auth.Token, error) {
	parseJwt, err := jwts.ParseJwt(token, secret, t.method, &auth.TokenClaims{})
	if err == nil || errors.Is(err, jwt.ErrTokenExpired) {
		return auth.Token{
			Token:       parseJwt.Token,
			Claims:      *parseJwt.Claims.(*auth.TokenClaims),
			TokenString: parseJwt.SignedJwt,
		}, nil
	} else {
		return auth.Token{}, err
	}
}

func (t *TokenHandler) verify(ctx context.Context, key, token string) (auth.Token, error) {
	parsedToken, err := t.parse(token, key)
	if err != nil {
		return parsedToken, err
	}
	// check if exists in cache
	if _, err := t.accessCache.Get(ctx, parsedToken.Claims.ID); errors.Is(err, redis.Nil) {
		return parsedToken, jwt.ErrTokenExpired
	} else if err != nil {
		return parsedToken, errs.Internal(err)
	}
	return parsedToken, nil
}
