package commands

type FlagDef struct {
	name      string
	shorthand string
	usage     string
}

var (

	// Common Flags
	yesFlag = FlagDef{
		name:      "yes",
		shorthand: "y",
		usage:     "Confirm action",
	}

	versionFlag = FlagDef{
		name:      "version",
		shorthand: "v",
		usage:     "Shows version info",
	}

	// Profile Command Flags
	profileNameFlag = FlagDef{
		name:      "name",
		shorthand: "n",
		usage:     "Profile name",
	}

	profileGitURI = FlagDef{
		name:  "git",
		usage: "Git URI to clone the profile from",
	}

	profileGitBranch = FlagDef{
		name:  "git-branch",
		usage: "Git branch corresponding to the git URI",
	}

	profileSetRootEnvFlag = FlagDef{
		name:  "root",
		usage: "Set the given environment as root",
	}

	// Env Command Flags
	envNameFlag = FlagDef{
		name:      "name",
		shorthand: "n",
		usage:     "Environment name",
	}

	envEmailFlag = FlagDef{
		name:      "email",
		shorthand: "e",
		usage:     "Environment email",
	}

	envTagsFlag = FlagDef{
		name:      "tags",
		shorthand: "t",
		usage:     "Environment tags",
	}

	envAddFlag = FlagDef{
		name:  "add",
		usage: "Adds environment to default profile",
	}

	envSearchFlag = FlagDef{
		name:      "search",
		shorthand: "s",
		usage:     "Searches query to filter environments",
	}

	envSelfFlag = FlagDef{
		name:  "self",
		usage: "Shares with the environment configured environment as self",
	}

	envDefFlag = FlagDef{
		name:      "env-def",
		shorthand: "e",
		usage:     "Environment definition",
	}

	// Provider Flags

	awsARNFlag = FlagDef{
		name:  "arn",
		usage: "ARN for the AWS KMS key",
	}

	gcpKmsResNameFlag = FlagDef{
		name:  "resource-name",
		usage: "GCP KMS resource name",
	}

	kmsRSAPublicKey = FlagDef{
		name:  "rsa-pubkey",
		usage: "KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding)",
	}

	// Vault Command Flags

	vaultFileFlag = FlagDef{
		name:      "vault",
		shorthand: "v",
		usage:     "Path to the vault file [Should end with .slv.yml or .slv.yaml]",
	}

	vaultAccessPublicKeysFlag = FlagDef{
		name:      "pubkey",
		shorthand: "k",
		usage:     "Public keys of environments that can access the vault",
	}

	vaultEnableHashingFlag = FlagDef{
		name:  "hash",
		usage: "Enables hashing by preserving a partial hash of the actual secret for the purpose of validating secret rotation [Not recommended, though it might be difficult to brute-force]",
	}

	vaultK8sFlag = FlagDef{
		name:  "k8s",
		usage: "Specify a name for the K8s SLV object if the vault is to be used in a K8s environment",
	}

	secretNameFlag = FlagDef{
		name:      "name",
		shorthand: "n",
		usage:     "Name of the secret",
	}

	secretValueFlag = FlagDef{
		name:  "secret",
		usage: "Secret to be added to the vault",
	}

	vaultImportFileFlag = FlagDef{
		name:  "file",
		usage: "Path to the YAML/JSON file to be imported",
	}

	secretForceUpdateFlag = FlagDef{
		name:  "force",
		usage: "Replaces the secret if it exists already",
	}

	vaultExportFormatFlag = FlagDef{
		name:  "format",
		usage: "List secrets as one of [json, yaml, envar]. Defaults to envar",
	}

	secretEncodeBase64Flag = FlagDef{
		name:  "base64",
		usage: "Encode the returned secret as base64",
	}

	vaultRefFileFlag = FlagDef{
		name:  "file",
		usage: "Path to the YAML/JSON file to be referenced",
	}

	vaultRefTypeFlag = FlagDef{
		name:  "format",
		usage: "Data serialization format of the referenced file",
	}

	vaultDerefPathFlag = FlagDef{
		name:  "path",
		usage: "Path to a file/directory to dereference secrets",
	}

	secretRefPreviewOnlyFlag = FlagDef{
		name:  "preview",
		usage: "Preview only mode",
	}
)
