package models

type UserCredentials struct {
	Name string `json:"name"`
}

type Signup struct {
	Name                 string `json:"name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type TotpRequired struct {
	VerificationToken string `json:"verificationToken"`
}

type EnableTotp struct {
	VerificationToken string `json:"verificationToken"`
}

type EnableTotpVerify struct {
	Code              int    `json:"code"`
	VerificationToken string `json:"verificationToken"`
}

type RecoveryCodes struct {
	RecoveryCodes []string `json:"recoveryCodes"`
}

type Totp struct {
	Code int `json:"code"`
}

type RecoveryCode struct {
	RecoveryCode      string `json:"recoveryCode"`
	VerificationToken string `json:"verificationToken"`
}
