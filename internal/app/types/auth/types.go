package auth

type LoginOption struct {
	// username or email
	Username string `json:"username"`
	// user password
	Password string `json:"password"`
	// remember user or not
	Remember bool `json:"remember"`
}

type RegisterOption struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

type ResetPasswordOption struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
	Code        string `json:"code"`
}

type RefreshTokenOption struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
