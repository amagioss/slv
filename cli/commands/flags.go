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
		usage:     "Path to the vault file [Must end with .vault.slv]",
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
