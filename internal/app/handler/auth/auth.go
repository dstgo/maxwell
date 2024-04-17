package auth

import (
	"crypto/sha512"
	"github.com/dstgo/maxwell/ent"
	"github.com/dstgo/maxwell/internal/app/data/repo"
	"github.com/dstgo/maxwell/internal/app/types/auth"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/str2bytes"
	"golang.org/x/net/context"
)

func NewAuthHandler(userRepo *repo.UserRepo, tokenHandler *TokenHandler, verify *VerifyCodeHandler) *AuthHandler {
	return &AuthHandler{token: tokenHandler, userRepo: userRepo, verifyCode: verify}
}

// AuthHandler is responsible for user authentication
type AuthHandler struct {
	token      *TokenHandler
	userRepo   *repo.UserRepo
	verifyCode *VerifyCodeHandler
}

func (a *AuthHandler) hash(s string) string {
	sum512 := sha512.Sum512(str2bytes.Str2Bytes(s))
	return str2bytes.Bytes2Str(sum512[:])
}

// LoginByPassword user login by password
func (a *AuthHandler) LoginByPassword(ctx context.Context, option auth.LoginOption) (*TokenPair, error) {
	// find user from repository
	queryUser, err := a.userRepo.FindByNameOrMail(ctx, option.Username)
	if ent.IsNotFound(err) {
		return nil, err
	} else if err != nil { // db error
		return nil, statuserr.InternalError(err)
	}

	// check password
	hashPaswd := a.hash(option.Password)
	if queryUser.Password != hashPaswd {
		return nil, auth.ErrPasswordMismatch
	}

	// issue token
	tokenPair, err := a.token.Issue(ctx, TokenPayload{
		Username: queryUser.Username,
		UserId:   queryUser.UID,
	}, option.Remember)

	if err != nil {
		return nil, statuserr.InternalError(err)
	}

	return &tokenPair, nil
}

// RegisterNewUser registers new user and returns it
func (a *AuthHandler) RegisterNewUser(ctx context.Context, option auth.RegisterOption) (*ent.User, error) {

	// check verify code if is valid
	err := a.verifyCode.CheckVerifyCode(ctx, option.Email, "register", option.Code)
	if err != nil {
		return nil, err
	}

	// check username if is duplicate
	userByName, err := a.userRepo.FindByName(ctx, option.Username)
	if !ent.IsNotFound(err) && err != nil {
		return nil, statuserr.InternalError(err)
	} else if userByName.UID != "" {
		return nil, auth.ErrUserAlreadyExists
	}

	// check email if is duplicate
	userByEmail, err := a.userRepo.FindByEmail(ctx, option.Email)
	if !ent.IsNotFound(err) && err != nil {
		return nil, statuserr.InternalError(err)
	} else if userByEmail.Email != "" {
		return nil, auth.ErrEmailAlreadyUsed
	}

	// create new user
	user, err := a.userRepo.CreateNewUser(ctx, option.Username, option.Email, option.Password)
	if err != nil {
		return nil, statuserr.InternalError(err)
	}
	return user, nil
}

// ResetPassword resets specified user password and returns uid
func (a *AuthHandler) ResetPassword(ctx context.Context, option auth.ResetPasswordOption) error {
	if option.Password != option.NewPassword {
		return auth.ErrPasswordInconsistent
	}

	// check verify code if is valid
	err := a.verifyCode.CheckVerifyCode(ctx, option.Email, "reset", option.Code)
	if err != nil {
		return err
	}

	queryUser, err := a.userRepo.FindByEmail(ctx, option.Email)
	if ent.IsNotFound(err) {
		return auth.ErrUserNotFund
	}

	if _, err := a.userRepo.UpdatePassword(ctx, queryUser.Email, a.hash(option.NewPassword)); err != nil {
		return statuserr.InternalError(err)
	}
	return nil
}
