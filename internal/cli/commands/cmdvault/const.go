package cmdvault

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
)

var (
	// Commands
	vaultCmd             *cobra.Command
	vaultNewCmd          *cobra.Command
	vaultUpdateCmd       *cobra.Command
	vaultAccessCmd       *cobra.Command
	vaultAccessAddCmd    *cobra.Command
	vaultAccessRemoveCmd *cobra.Command
	vaultPutCmd          *cobra.Command
	vaultDeleteCmd       *cobra.Command
	vaultGetCmd          *cobra.Command
	vaultRunCmd          *cobra.Command
	vaultRefCmd          *cobra.Command
	vaultDerefCmd        *cobra.Command
)

var (
	// Flags
	vaultFileFlag = utils.FlagDef{
		Name:      "vault",
		Shorthand: "v",
		Usage:     "Path to the vault file [Should end with .slv.yml or .slv.yaml]",
	}

	vaultEnableHashingFlag = utils.FlagDef{
		Name:  "hash",
		Usage: "Enables hashing by preserving a partial hash of the actual secret for the purpose of validating secret rotation [Not recommended, though it might be difficult to brute-force]",
	}

	vaultNameFlag = utils.FlagDef{
		Name:  "name",
		Usage: "Name for the SLV vault",
	}

	vaultK8sNamespaceFlag = utils.FlagDef{
		Name:  "k8s-namespace",
		Usage: "Namespace for the K8s SLV resource",
	}

	vaultK8sSecretTypeFlag = utils.FlagDef{
		Name:  "k8s-secret-type",
		Usage: "Type of the K8s Secret (e.g. Opaque, kubernetes.io/tls, etc.)",
	}

	vaultK8sSecretFlag = utils.FlagDef{
		Name:  "k8s-secret",
		Usage: "A K8s Secret that needs to be transformed to an SLV vault (Use '-' to read from stdin)",
	}

	varNamePrefixFlag = utils.FlagDef{
		Name:  "prefix",
		Usage: "Prefix to set to the secret name while setting it as the environment variable",
	}

	vaultShellCommandFlag = utils.FlagDef{
		Name:      "command",
		Shorthand: "c",
		Usage:     "Command to run in the shell",
	}

	itemNameFlag = utils.FlagDef{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Name of the item",
	}

	itemValueFlag = utils.FlagDef{
		Name:  "value",
		Usage: "Value of the item to be used (Use - to read from stdin)",
	}

	deprecatedSecretFlag = utils.FlagDef{
		Name:  "secret",
		Usage: "Secret to be added to the vault",
	}

	vaultImportFileFlag = utils.FlagDef{
		Name:  "file",
		Usage: "Path to the YAML/JSON/ENV file to be imported",
	}

	plaintextValueFlag = utils.FlagDef{
		Name:  "plaintext",
		Usage: "Indicates that the value will be stored as plaintext (use only for config values that are not sensitive)",
	}

	secretForceUpdateFlag = utils.FlagDef{
		Name:  "force",
		Usage: "Replaces the secret if it exists already",
	}

	vaultExportFormatFlag = utils.FlagDef{
		Name:  "format",
		Usage: "List secrets as one of [json, yaml, envar]",
	}

	valueWithMetadata = utils.FlagDef{
		Name:  "with-metadata",
		Usage: "Returns the vault values with metadata",
	}

	valueEncodeBase64Flag = utils.FlagDef{
		Name:  "base64",
		Usage: "Encode the returned value as base64",
	}

	vaultRefFileFlag = utils.FlagDef{
		Name:  "file",
		Usage: "Path to the YAML/JSON/blob file to be referenced",
	}

	secretSubstitutionPreviewOnlyFlag = utils.FlagDef{
		Name:  "preview",
		Usage: "Enables preview mode (shows the substitution result without writing to the file)",
	}
)
