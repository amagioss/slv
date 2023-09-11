package groups

import "github.com/shibme/slv/core/crypto"

type Group struct {
	Id          string               `yaml:"id"`
	Name        string               `yaml:"name"`
	Description string               `yaml:"description"`
	Email       string               `yaml:"email"`
	Access      []*crypto.WrappedKey `yaml:"access"`
}
