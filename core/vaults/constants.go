package vaults

import (
	"errors"
	"regexp"
	"strings"

	"github.com/shibme/slv/core/crypto"
)

const (
	vaultFileExtension                  = ".slv"
	VaultKey             crypto.KeyType = 'V'
	secretNamePattern                   = "[A-Za-z][A-Za-z0-9_]+[A-Za-z0-9]"
	secretRefAbbrev                     = "VSR" // VSR = Vault Secret Reference
	vaultIdPattern                      = "[A-Za-z0-9]+"
	secretRefPatternBase                = `\{\{\s*SLV_VSR_VAULTID\.` + secretNamePattern + `\s*\}\}`
)

var secretNameRegex = regexp.MustCompile(secretNamePattern)
var secretRefRegex = regexp.MustCompile(strings.ReplaceAll(secretRefPatternBase, "VAULTID", vaultIdPattern))

var ErrInvalidVaultFileName = errors.New("invalid vault file name [vault file name must end in " + vaultFileExtension + "]")
var ErrVaultDirPathCreation = errors.New("error in creating a new vault directory path")
var ErrVaultNotAccessible = errors.New("vault is not accessible using the given environment key")
var ErrVaultLocked = errors.New("the vault is currently locked")
var ErrVaultExists = errors.New("vault exists already")
var ErrVaultNotFound = errors.New("vault not found")
var ErrVaultCannotBeSharedWithVault = errors.New("vault cannot be shared with another vault")
var ErrInvalidSecretName = errors.New("invalid secret name format [secret name must start with a letter and can only contain letters, numbers and underscores]")
var ErrVaultSecretExistsAlready = errors.New("secret exists already for the given name")
var ErrVaultSecretNotFound = errors.New("no secret found for the given name")
var ErrInvalidReferenceFormat = errors.New("invalid reference format. references must follow the pattern {{SLV_VSR_VAULTID.secretName}} to allow dereferencing")
