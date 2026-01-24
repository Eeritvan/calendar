package models

type UserCredentials struct {
	Name string `json:"name"`
}

type Signup struct {
	Name                 string `json:"name" validate:"required,max=100"`
	Password             string `json:"password" validate:"required,min=8,max=100"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password,min=8,max=100"`
}

type Login struct {
	Name     string `json:"name" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type TotpRequired struct {
	VerificationToken string `json:"verificationToken"`
}

type EnableTotp struct {
	VerificationToken string `json:"verificationToken"`
}

type EnableTotpVerify struct {
	Code              string `json:"code" validate:"required,numeric,len=6"`
	VerificationToken string `json:"verificationToken" validate:"required,jwt"`
}

type RecoveryCodes struct {
	RecoveryCodes []string `json:"recoveryCodes"`
}

type Totp struct {
	Code string `json:"code" validate:"required,numeric,len=6"`
}

type RecoveryCode struct {
	RecoveryCode      string `json:"recoveryCode" validate:"required,len=14"`
	VerificationToken string `json:"verificationToken" validate:"required,jwt"`
}
