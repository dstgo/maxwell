package auth

type LoginOption struct {
	// username or email
	Username string `json:"username" binding:"required"`
	// user password
	Password string `json:"password" binding:"required"`
	// remember user or not
	Remember bool `json:"remember"`
}

type RegisterOption struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"email"`
	Code     string `json:"code" binding:"required"`
}

type ResetPasswordOption struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type RefreshTokenOption struct {
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type VerifyCodeOption struct {
	To    string `json:"to" binding:"email"`
	Usage Usage  `json:"usage" binding:"required"`
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
