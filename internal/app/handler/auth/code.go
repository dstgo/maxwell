package auth

import (
	"errors"
	"github.com/dstgo/maxwell/internal/app/handler/email"
	"github.com/dstgo/maxwell/internal/app/types/auth"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"math/rand/v2"
	"time"
)

func NewVerifyCodeHandler(client *redis.Client, sender *email.Handler) *VerifyCodeHandler {
	return &VerifyCodeHandler{
		cache:  client,
		sender: sender,
	}
}

type VerifyCodeHandler struct {
	cache  *redis.Client
	sender *email.Handler
}

// SendVerifyCodeEmail send a verify code email to the specified address
func (v *VerifyCodeHandler) SendVerifyCodeEmail(ctx context.Context, to, usage string, ttl, retryttl time.Duration) error {

	// check retry ttl
	retryRes, err := v.cache.Get(ctx, to).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		return statuserr.InternalError(err)
	} else if err == nil && retryRes != "" {
		return auth.ErrVerifyCodeRetryLater
	}

	// set retry ttl
	if err := v.cache.Set(ctx, to, "", retryttl).Err(); err != nil {
		return statuserr.InternalError(err)
	}

	// generate verify code
	verifyCode := NewVerifyCode(8)
	codeKey := usage + ":" + verifyCode

	// check verify code if is repeated
	for i := 0; i < 10; i++ {
		_, err := v.cache.Get(ctx, codeKey).Result()
		if errors.Is(err, redis.Nil) {
			break
		} else if err != nil {
			return statuserr.InternalError(err)
		} else { // repeat,
			// regenerate a new one
			verifyCode = NewVerifyCode(8)
		}
	}
	// set verify code
	if _, err := v.cache.Set(ctx, codeKey, to, ttl).Result(); err != nil {
		return statuserr.InternalError(err)
	}

	// send email
	err = v.sender.SendHermesEmail(ctx, "you are applying for verification code", []string{to}, email.ConfirmCodeMail(usage, to, verifyCode, ttl))
	if err != nil {
		v.cache.Del(ctx, codeKey)
		return err
	}

	return nil
}

// CheckVerifyCode check verify code if is valid
func (v *VerifyCodeHandler) CheckVerifyCode(ctx context.Context, to, usage, code string) error {
	codeKey := usage + ":" + code
	getTo, err := v.cache.Get(ctx, codeKey).Result()
	if errors.Is(err, redis.Nil) || getTo != to {
		return auth.ErrVerifyCodeInvalid
	} else if err != nil {
		return statuserr.InternalError(err)
	}
	return nil
}

// NewVerifyCode returns a verification code with specified length.
// the bigger n is, the conflicts will be less, the recommended n is 8.
func NewVerifyCode(n int) string {
	code := make([]byte, n)
	for i, _ := range code {
		if rand.Int()%2 == 1 {
			code[i] = byte(rand.Int() % 10)
		} else {
			code[i] = 'A' + byte(rand.Int()%26)
		}
	}
	return string(code)
}
