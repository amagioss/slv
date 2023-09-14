package vaults

import (
	"errors"

	"github.com/shibme/slv/core/crypto"
)

const (
	vaultFileExtension                 = ".vault.slv"
	VaultKey            crypto.KeyType = 'V'
	directRefAbbrev                    = "DRS" // DRS = Direct Referenced Secret
	autoRefAbbrev                      = "ARS" // ARS = Auto Referenced Secret
	autoReferenceLength                = 16
)

var ErrInvalidVaultFileName = errors.New("invalid vault file name")
var ErrVaultDirPathCreation = errors.New("error in creating a new vault directory path")
var ErrVaultNotAccessible = errors.New("vault is not accessible using the given environment key")
var ErrVaultLocked = errors.New("the vault is currently locked")
var ErrVaultExists = errors.New("vault exists already")
var ErrVaultNotFound = errors.New("vault not found")
var ErrVaultCannotBeSharedWithVault = errors.New("vault cannot be shared with another vault")
var ErrVaultSecretNotFound = errors.New("no secret found for the given name")
var ErrMaximumReferenceAttemptsReached = errors.New("maximum reference attempts reached")
var ErrInvalidReferenceFileFormat = errors.New("invalid reference file - only yaml and json are supported")
var ErrInvalidRefActionType = errors.New("invalid reference action type")
var ErrInvalidAutoRefString = errors.New("invalid auto reference string")
var ErrInvalidDirectRefString = errors.New("invalid direct reference string")
