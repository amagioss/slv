package cmdvault

import (
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
)

var (
	// Commands
	vaultCmd             *cobra.Command
	vaultNewCmd          *cobra.Command
	vaultToK8sCmd        *cobra.Command
	vaultAccessCmd       *cobra.Command
	vaultAccessAddCmd    *cobra.Command
	vaultAccessRemoveCmd *cobra.Command
	vaultPutCmd          *cobra.Command
	vaultDeleteCmd       *cobra.Command
	vaultGetCmd          *cobra.Command
	vaultShellCmd        *cobra.Command
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

	vaultAccessPublicKeysFlag = utils.FlagDef{
		Name:      "env-pubkey",
		Shorthand: "k",
		Usage:     "Public keys of environments that can access the vault",
	}

	vaultAccessK8sFlag = utils.FlagDef{
		Name:  "env-k8s",
		Usage: "Shares vault access with the accessible k8s cluster",
	}

	vaultEnableHashingFlag = utils.FlagDef{
		Name:  "hash",
		Usage: "Enables hashing by preserving a partial hash of the actual secret for the purpose of validating secret rotation [Not recommended, though it might be difficult to brute-force]",
	}

	vaultK8sNameFlag = utils.FlagDef{
		Name:  "k8s-name",
		Usage: "Name for the K8s SLV resource",
	}

	vaultK8sNamespaceFlag = utils.FlagDef{
		Name:  "k8s-namespace",
		Usage: "Namespace for the K8s SLV resource",
	}

	vaultK8sSecretFlag = utils.FlagDef{
		Name:  "k8s-secret",
		Usage: "A K8s Secret that needs to be transformed to an SLV vault",
	}

	secretNamePrefixFlag = utils.FlagDef{
		Name:  "prefix",
		Usage: "Prefix to set to the secret name while setting it as the environment variable",
	}

	secretNameFlag = utils.FlagDef{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Name of the secret",
	}

	secretValueFlag = utils.FlagDef{
		Name:  "secret",
		Usage: "Secret to be added to the vault",
	}

	vaultImportFileFlag = utils.FlagDef{
		Name:  "file",
		Usage: "Path to the YAML/JSON file to be imported",
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
		Usage: "Path to the YAML/JSON file to be referenced",
	}

	vaultRefTypeFlag = utils.FlagDef{
		Name:  "format",
		Usage: "Data serialization format of the referenced file",
	}

	vaultDerefPathFlag = utils.FlagDef{
		Name:  "path",
		Usage: "Path to a file/directory to dereference secrets",
	}

	secretRefPreviewOnlyFlag = utils.FlagDef{
		Name:  "preview",
		Usage: "Preview only mode",
	}
)
