package input

import "sync"

var (
	defaultPwdPolicyMutex sync.Mutex
	defaultPwdPolicy      *PasswordPolicy
)

const (
	pwdSpecialCharset         = `!@#$%^&*_-+=?:;,.|\/[](){}<>`
	pwdDefaultMinLength       = 10
	pwdDefaultMinUppercase    = 1
	pwdDefaultMinLowercase    = 1
	pwdDefaultMinDigits       = 1
	pwdDefaultMinSpecialChars = 1
)
