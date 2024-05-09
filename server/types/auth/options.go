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
	// username must be alphanumeric
	Username string `json:"username" binding:"required,alphanum"`
	// user password
	Password string `json:"password" binding:"required"`
	// user email address
	Email string `json:"email" binding:"email"`
	// verification code from verify email
	Code string `json:"code" binding:"required,alphanum"`
}

type ResetPasswordOption struct {
	// user email address
	Email string `json:"email" binding:"email"`
	// new password
	Password string `json:"password" binding:"required"`
	// verification code from verify email
	Code string `json:"code" binding:"required"`
}

type RefreshTokenOption struct {
	// access token
	AccessToken string `json:"accessToken" binding:"required"`
	// refresh token
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type VerifyCodeOption struct {
	// email receiver
	To string `json:"to" binding:"email"`
	// verify code usage: 1-register 2-reset password
	Usage Usage `json:"usage" binding:"required"`
}
