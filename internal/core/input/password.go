package input

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type PasswordPolicy struct {
	MinLength       uint8
	MinUppercase    uint8
	MinLowercase    uint8
	MinDigits       uint8
	MinSpecialChars uint8
}

func DefaultPasswordPolicy() *PasswordPolicy {
	if defaultPwdPolicy == nil {
		defaultPwdPolicyMutex.Lock()
		defer defaultPwdPolicyMutex.Unlock()
		if defaultPwdPolicy == nil {
			defaultPwdPolicy = &PasswordPolicy{
				MinLength:       pwdDefaultMinLength,
				MinUppercase:    pwdDefaultMinUppercase,
				MinLowercase:    pwdDefaultMinLowercase,
				MinDigits:       pwdDefaultMinDigits,
				MinSpecialChars: pwdDefaultMinSpecialChars,
			}
		}
	}
	return defaultPwdPolicy
}

func (policy *PasswordPolicy) error() error {
	return fmt.Errorf("password must be at least %d characters long, contain at least %d uppercase, %d lowercase, %d digits, and %d special characters",
		policy.MinLength, policy.MinUppercase, policy.MinLowercase, policy.MinDigits, policy.MinSpecialChars)
}

func (policy *PasswordPolicy) Validate(password string) error {
	if len(password) < int(policy.MinLength) {
		return policy.error()
	}
	var (
		upp, low, num, sym uint8
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp++
		case unicode.IsLower(char):
			low++
		case unicode.IsNumber(char):
			num++
		case strings.ContainsRune(pwdSpecialCharset, char):
			sym++
		default:
			return policy.error()
		}
	}
	if !(upp >= policy.MinUppercase && low >= policy.MinLowercase && num >= policy.MinDigits && sym >= policy.MinSpecialChars) {
		return policy.error()
	}
	return nil
}

func NewPasswordFromUser(policy *PasswordPolicy) ([]byte, error) {
	password, err := GetHiddenInput("Enter a Password: ")
	if err != nil {
		return nil, err
	}
	if policy != nil {
		if err = policy.Validate(string(password)); err != nil {
			return nil, err
		}
	}
	if confirmPassword, err := GetHiddenInput("Confirm Password: "); err != nil {
		return nil, err
	} else if !bytes.Equal(password, confirmPassword) {
		return nil, fmt.Errorf("passwords do not match")
	}
	return password, nil
}
