package commands

type FlagDef struct {
	name      string
	shorthand string
	usage     string
}

// Common Flags
var yesFlag = FlagDef{
	name:      "yes",
	shorthand: "y",
	usage:     "Confirm action",
}

var versionFlag = FlagDef{
	name:      "version",
	shorthand: "v",
	usage:     "Shows version info",
}

// Profile Command Flags
var profileNameFlag = FlagDef{
	name:      "name",
	shorthand: "n",
	usage:     "Profile name",
}

var profileEnvDefFlag = FlagDef{
	name:      "env",
	shorthand: "e",
	usage:     "Environment definition",
}

var profileSetRootEnvFlag = FlagDef{
	name:  "root",
	usage: "Set the given environment as root",
}

// Env Command Flags
var envNameFlag = FlagDef{
	name:      "name",
	shorthand: "n",
	usage:     "Environment name",
}

var envEmailFlag = FlagDef{
	name:      "email",
	shorthand: "e",
	usage:     "Environment email",
}

var envTagsFlag = FlagDef{
	name:      "tags",
	shorthand: "t",
	usage:     "Environment tags",
}

var envKMSTypeFlag = FlagDef{
	name:  "kms-type",
	usage: "KMS provider [AWS, GCP, AZURE]",
}

var envKMSIdFlag = FlagDef{
	name:  "kms-id",
	usage: "KMS ID",
}

var envKMSPemFlag = FlagDef{
	name:  "kms-pubkey",
	usage: "KMS public key [RSA 4096] file (pem format)",
}

var envSelfFlag = FlagDef{
	name:      "self",
	shorthand: "u",
	usage:     "Creates a user environment for you and registers locally",
}

var envAddFlag = FlagDef{
	name:      "add",
	shorthand: "a",
	usage:     "Adds environment to default profile",
}

var envSearchFlag = FlagDef{
	name:      "search-env",
	shorthand: "s",
	usage:     "Searches query to filter environments",
}

// KMS Flags

var kmsAWSARNFlag = FlagDef{
	name:  "arn",
	usage: "ARN for the AWS KMS key",
}

var kmsRSAPublicKey = FlagDef{
	name:  "rsa-pubkey",
	usage: "KMS public key [RSA 4096] file (pem format)",
}

// Vault Command Flags

var vaultFileFlag = FlagDef{
	name:      "vault-file",
	shorthand: "v",
	usage:     "Path to the vault file [Must end with .vault.slv]",
}

var vaultAccessPublicKeysFlag = FlagDef{
	name:      "public-keys",
	shorthand: "k",
	usage:     "Public keys of environments or groups that can access the vault",
}

var vaultEnableHashingFlag = FlagDef{
	name:  "enable-hash",
	usage: "Preserve a partial secret hash for the purpose of validating secret rotation [Not recommended, though it might be resilient from brute-forcing]",
}

// Secret Command Flags

var secretNameFlag = FlagDef{
	name:      "name",
	shorthand: "n",
	usage:     "Name of the secret",
}

var secretValueFlag = FlagDef{
	name:      "secret",
	shorthand: "s",
	usage:     "Secret to be added to the vault",
}

var secretForceUpdateFlag = FlagDef{
	name:  "force",
	usage: "Replaces the secret if it exists already",
}

var secretRefFileFlag = FlagDef{
	name:      "file",
	shorthand: "f",
	usage:     "Path to the YAML/JSON file",
}

var secretRefTypeFlag = FlagDef{
	name:      "type",
	shorthand: "t",
	usage:     "Data format to be considered for the file to be referenced",
}

var secretRefPreviewOnlyFlag = FlagDef{
	name:      "preview",
	shorthand: "p",
	usage:     "Preview only mode",
}
