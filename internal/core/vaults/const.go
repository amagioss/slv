package vaults

import (
	"errors"
	"regexp"
	"strings"

	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
)

const (
	vaultFileNameEnding                 = config.AppNameLowerCase
	VaultKey             crypto.KeyType = 'V'
	vaultIdLength                       = 30
	secretNamePattern                   = "[a-zA-Z]([a-zA-Z0-9_]*[a-zA-Z0-9])?"
	secretRefPatternBase                = `\{\{\s*VAULTID\.` + secretNamePattern + `\s*\}\}`
	vaultIdAbbrev                       = "VID"

	k8sApiVersion = config.K8SLVGroup + "/" + config.K8SLVVersion
	k8sKind       = config.K8SLVKind
	k8sVaultField = config.K8SLVVaultField
)

var (
	secretNameRegex = regexp.MustCompile(secretNamePattern)
	secretRefRegex  = regexp.MustCompile(strings.ReplaceAll(secretRefPatternBase, "VAULTID", config.AppNameUpperCase+"_"+vaultIdAbbrev+"_[A-Za-z0-9]+"))

	errVaultVersionNotRecognized    = errors.New("vault version not recognized")
	errGeneratingVaultId            = errors.New("error in generating a new vault id")
	errInvalidVaultFileName         = errors.New("invalid vault file name [vault file name must end with " + vaultFileNameEnding + ".yml or " + vaultFileNameEnding + ".yaml]")
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
	errInvalidReferenceFormat       = errors.New("invalid reference format. references must follow the pattern {{" + config.AppNameUpperCase + "_" + vaultIdAbbrev + "_ABCXYZ.secretName}} to allow dereferencing")
	errInvalidImportDataFormat      = errors.New("invalid import data format - expected a map of string to string [secretName: secretValue] in YAML/JSON format")
	errK8sNameRequired              = errors.New("k8s resource name is required for k8s compatible SLV vault")
	errVaultWrappedKeysNotFound     = errors.New("vault wrapped keys not found - vault will be inaccessible by any environment")
)
