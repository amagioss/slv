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
		name:  "git-uri",
		usage: "Profile git URI",
	}

	profileGitBranch = FlagDef{
		name:  "git-branch",
		usage: "Profile git branch",
	}

	profileEnvDefFlag = FlagDef{
		name:      "env",
		shorthand: "e",
		usage:     "Environment definition",
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

	envSelfFlag = FlagDef{
		name:      "self",
		shorthand: "u",
		usage:     "Creates a user environment for you and registers locally",
	}

	envAddFlag = FlagDef{
		name:      "add",
		shorthand: "a",
		usage:     "Adds environment to default profile",
	}

	envSearchFlag = FlagDef{
		name:      "search-env",
		shorthand: "s",
		usage:     "Searches query to filter environments",
	}

	// KMS Flags

	kmsAWSARNFlag = FlagDef{
		name:  "arn",
		usage: "ARN for the AWS KMS key",
	}

	kmsRSAPublicKey = FlagDef{
		name:  "rsa-pubkey",
		usage: "KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding)",
	}

	// Vault Command Flags

	vaultFileFlag = FlagDef{
		name:      "vault-file",
		shorthand: "v",
		usage:     "Path to the vault file [Must end with .slv.yml or .slv.yaml]",
	}

	vaultAccessPublicKeysFlag = FlagDef{
		name:      "public-keys",
		shorthand: "k",
		usage:     "Public keys of environments or groups that can access the vault",
	}

	vaultEnableHashingFlag = FlagDef{
		name:  "enable-hash",
		usage: "Preserve a partial secret hash for the purpose of validating secret rotation [Not recommended, though it might be resilient from brute-forcing]",
	}

	vaultK8sFlag = FlagDef{
		name:  "k8s",
		usage: "Specify a name for the K8s SLV object if the vault is to be used in a K8s environment",
	}

	// Secret Command Flags

	secretNameFlag = FlagDef{
		name:      "name",
		shorthand: "n",
		usage:     "Name of the secret",
	}

	secretValueFlag = FlagDef{
		name:      "secret",
		shorthand: "s",
		usage:     "Secret to be added to the vault",
	}

	secretForceUpdateFlag = FlagDef{
		name:  "force",
		usage: "Replaces the secret if it exists already",
	}

	secretListFormatFlag = FlagDef{
		name:  "format",
		usage: "List secrets as one of [json, yaml, table, envars]. Defaults to table.",
	}

	secretEncodeBase64Flag = FlagDef{
		name:  "base64",
		usage: "Encode the returned secret as base64",
	}

	secretRefFileFlag = FlagDef{
		name:      "file",
		shorthand: "f",
		usage:     "Path to the YAML/JSON file",
	}

	secretRefTypeFlag = FlagDef{
		name:      "type",
		shorthand: "t",
		usage:     "Data format to be considered for the file to be referenced",
	}

	secretRefPreviewOnlyFlag = FlagDef{
		name:      "preview",
		shorthand: "p",
		usage:     "Preview only mode",
	}
)
