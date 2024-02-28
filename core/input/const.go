package input

import "fmt"

var (
	errNonInteractiveTerminal = fmt.Errorf("non-interactive terminal")
	defaultPwdPolicy          *PasswordPolicy
)

const (
	pwdSpecialCharset         = "!@#$%^&*?"
	pwdDefaultMinLength       = 10
	pwdDefaultMinUppercase    = 1
	pwdDefaultMinLowercase    = 1
	pwdDefaultMinDigits       = 1
	pwdDefaultMinSpecialChars = 1
)
