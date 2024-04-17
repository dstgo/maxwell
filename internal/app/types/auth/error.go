package auth

import (
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

var (
	ErrUserNotFund          = statuserr.Errorf("user not found").SetCode(1_400_000)
	ErrUserAlreadyExists    = statuserr.Errorf("user already exists").SetCode(1_400_002)
	ErrPasswordMismatch     = statuserr.Errorf("password mismatch").SetCode(1_400_004)
	ErrPasswordInconsistent = statuserr.Errorf("password inconsistent").SetCode(1_400_008)
	ErrEmailAlreadyUsed     = statuserr.Errorf("email already used by other").SetCode(1_400_016)

	ErrVerifyCodeRetryLater = statuserr.Errorf("retry applying for verify code later").SetCode(1_400_032)
	ErrVerifyCodeInvalid    = statuserr.Errorf("verify code invalid").SetCode(1_400_033)
)
