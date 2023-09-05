package vaults

import (
	"errors"

	"github.com/shibme/slv/core/crypto"
)

const (
	vaultFileExtension                = ".vault.slv"
	VaultKey           crypto.KeyType = 'V'
	maxRefNameAttempts                = 10

	referencedSecretPreviewVal = "SLV_SS_ENCRYPTED_VALUE"
)

var ErrInvalidVaultFileName = errors.New("invalid vault file name")
var ErrReadingVault = errors.New("error in reading the vault file")
var ErrVaultDirPathCreation = errors.New("error in creating a new vault directory path")
var ErrVaultNotAccessible = errors.New("vault is not accessible using the given environment key")
var ErrVaultLocked = errors.New("the vault is currently locked")
var ErrVaultExists = errors.New("vault exists already")
var ErrVaultNotFound = errors.New("vault not found")
var ErrVaultCannotBeSharedWithVault = errors.New("vault cannot be shared with another vault")
var ErrVaultAlreadySharedWithKey = errors.New("vault already shared with the given key")
var ErrVaultSecretNotFound = errors.New("no secret found for the given name")
var ErrMissingVaultPublicKey = errors.New("missing vault public key")
var ErrMaximumReferenceAttemptsReached = errors.New("maximum reference attempts reached")
