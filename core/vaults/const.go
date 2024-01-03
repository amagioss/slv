package vaults

import (
	"errors"
	"regexp"
	"strings"

	"github.com/amagimedia/slv/core/crypto"
)

const (
	vaultFileExtension                  = ".slv"
	VaultKey             crypto.KeyType = 'V'
	secretNamePattern                   = "[a-zA-Z]([a-zA-Z0-9_]*[a-zA-Z0-9])?"
	secretRefAbbrev                     = "VSR" // VSR = Vault Secret Reference
	secretRefPatternBase                = `\{\{\s*VAULTID\.` + secretNamePattern + `\s*\}\}`
)

var (
	secretNameRegex = regexp.MustCompile(secretNamePattern)
	secretRefRegex  = regexp.MustCompile(strings.ReplaceAll(secretRefPatternBase, "VAULTID", "SLV_VPK_[A-Za-z0-9]+"))

	errInvalidVaultFileName         = errors.New("invalid vault file name [vault file name must end in " + vaultFileExtension + "]")
	errVaultDirPathCreation         = errors.New("error in creating a new vault directory path")
	errVaultNotAccessible           = errors.New("vault is not accessible using the given environment key")
	errVaultLocked                  = errors.New("the vault is currently locked")
	errVaultExists                  = errors.New("vault exists already")
	errVaultNotFound                = errors.New("vault not found")
	errVaultCannotBeSharedWithVault = errors.New("vault cannot be shared with another vault")
	errInvalidSecretName            = errors.New("invalid secret name format [secret name must start with a letter and can only contain letters, numbers and underscores]")
	errVaultSecretExistsAlready     = errors.New("secret exists already for the given name")
	errVaultSecretNotFound          = errors.New("no secret found for the given name")
	errVaultPublicKeyNotFound       = errors.New("vault public key not found")
	errInvalidReferenceFormat       = errors.New("invalid reference format. references must follow the pattern {{SLV_VSR_VAULTID.secretName}} to allow dereferencing")
)
