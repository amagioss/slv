package commands

type FlagDef struct {
	name      string
	shorthand string
	usage     string
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

var envAddFlag = FlagDef{
	name:      "add",
	shorthand: "a",
	usage:     "Adds environment to default profile",
}

var envSearchFlag = FlagDef{
	name:      "search",
	shorthand: "s",
	usage:     "Searches for environments",
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
	name:      "enable-hash",
	shorthand: "h",
	usage:     "Enable hashing for secrets to verify secret rotation [Not recommended, though it may be resilient from attacks]",
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

var secretRefFileFlag = FlagDef{
	name:      "file",
	shorthand: "f",
	usage:     "Path to the YAML/JSON file",
}

var secretRefPreviewOnlyFlag = FlagDef{
	name:      "preview",
	shorthand: "p",
	usage:     "Preview only mode",
}
