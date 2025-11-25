package vaults

import (
	"errors"
	"regexp"
	"strings"

	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
)

const (
	vaultFileNameRawExt                        = config.AppNameLowerCase
	vaultFileNameDesiredExt                    = "." + vaultFileNameRawExt + ".yaml"
	VaultKey                    crypto.KeyType = 'V'
	vaultIdLength                              = 30
	secretNamePattern                          = `([\w]+)`
	vaultNamePattern                           = `([a-zA-Z0-9_-]+)`
	vaultNamePatternPlaceholder                = "VAULTNAME"
	secretRefPatternBase                       = `\{\{\s*(SLV|slv)\.` + vaultNamePatternPlaceholder + `\.` + secretNamePattern + `\s*\}\}`

	k8sApiVersion           = config.K8SLVGroup + "/" + config.K8SLVVersion
	k8sKind                 = config.K8SLVKind
	k8sVaultSpecField       = config.K8SLVVaultField
	k8sVersionAnnotationKey = config.K8SLVAnnotationVersionKey
)

var (
	secretNameRegex                = regexp.MustCompile(secretNamePattern)
	unsupportedSecretNameCharRegex = regexp.MustCompile(`[^\w]`)
	secretRefRegex                 = regexp.MustCompile(strings.ReplaceAll(secretRefPatternBase, vaultNamePatternPlaceholder, vaultNamePattern))

	errVaultDirPathCreation         = errors.New("error in creating a new vault directory path")
	errVaultNotAccessible           = errors.New("vault is not accessible by the environment")
	errVaultLocked                  = errors.New("the vault is currently locked")
	errVaultExists                  = errors.New("vault exists already")
	errVaultNotFound                = errors.New("vault not found")
	errVaultCannotBeSharedWithVault = errors.New("vault cannot be shared with another vault")
	errInvalidVaultItemName         = errors.New("invalid name format [name must start with a letter and can only contain letters, numbers and underscores]")
	errVaultItemExistsAlready       = errors.New("item exists already for the given name")
	errVaultItemNotFound            = errors.New("no item found for the given name")
	errVaultPublicKeyNotFound       = errors.New("vault public key not found")
	errInvalidReferenceFormat       = errors.New("invalid reference format. references must follow the pattern {{" + config.AppNameUpperCase + ".<vault_name>.<secret_name>}} to allow dereferencing")
	errInvalidImportDataFormat      = errors.New("invalid import data format - expected a map of string to string [secretName: secretValue] in YAML/JSON/ENV format")
	errK8sNameRequired              = errors.New("k8s resource name is required for a k8s compatible SLV vault")
	errVaultWrappedKeysNotFound     = errors.New("vault wrapped keys not found - vault will be inaccessible by any environment")
	errVaultNotWritable             = errors.New("vault is not writable")
)
