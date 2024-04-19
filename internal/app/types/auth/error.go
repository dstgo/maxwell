package auth

import (
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

var (
	ErrUserNotFund          = statuserr.Errorf("user not found").SetCode(1_400_000).SetStatus(status.BadRequest)
	ErrUserAlreadyExists    = statuserr.Errorf("user already exists").SetCode(1_400_002).SetStatus(status.BadRequest)
	ErrPasswordMismatch     = statuserr.Errorf("password mismatch").SetCode(1_400_004).SetStatus(status.BadRequest)
	ErrPasswordInconsistent = statuserr.Errorf("password inconsistent").SetCode(1_400_008).SetStatus(status.BadRequest)
	ErrEmailAlreadyUsed     = statuserr.Errorf("email already used by other").SetCode(1_400_016).SetStatus(status.BadRequest)

	ErrVerifyCodeRetryLater = statuserr.Errorf("retry applying for verify code later").SetCode(1_400_032).SetStatus(status.BadRequest)
	ErrVerifyCodeInvalid    = statuserr.Errorf("invliad verify code").SetCode(1_400_033).SetStatus(status.BadRequest)

	ErrVerifyCodeUsageUnsupported = statuserr.Errorf("verify code usage unsupported").SetCode(1_400_036).SetStatus(status.BadRequest)
)
