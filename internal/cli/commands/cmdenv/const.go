package cmdenv

import (
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
)

var (
	// Commands
	envCmd           *cobra.Command
	envNewCmd        *cobra.Command
	envNewServiceCmd *cobra.Command
	envNewUserCmd    *cobra.Command
	envAddCmd        *cobra.Command
	envListSearchCmd *cobra.Command
	envDeleteCmd     *cobra.Command
	envSelfCmd       *cobra.Command
	envSelfSetCmd    *cobra.Command
)

var (
	// Flags

	envNameFlag = utils.FlagDef{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Environment name",
	}

	envEmailFlag = utils.FlagDef{
		Name:      "email",
		Shorthand: "e",
		Usage:     "Environment email",
	}

	envTagsFlag = utils.FlagDef{
		Name:      "tags",
		Shorthand: "t",
		Usage:     "Environment tags",
	}

	envAddFlag = utils.FlagDef{
		Name:  "add",
		Usage: "Adds environment to default profile",
	}

	envSetRootFlag = utils.FlagDef{
		Name:  "root",
		Usage: "Set the given environment as root",
	}

	EnvSearchFlag = utils.FlagDef{
		Name:      "search",
		Shorthand: "s",
		Usage:     "Searches query to filter environments",
	}

	EnvSelfFlag = utils.FlagDef{
		Name:  "self",
		Usage: "Shares with the environment configured environment as self",
	}

	envDefFlag = utils.FlagDef{
		Name:      "env-def",
		Shorthand: "e",
		Usage:     "Environment definition",
	}

	// Provider Flags
	awsARNFlag = utils.FlagDef{
		Name:  "arn",
		Usage: "ARN for the AWS KMS key",
	}

	gcpKmsResNameFlag = utils.FlagDef{
		Name:  "resource-name",
		Usage: "GCP KMS resource name",
	}

	kmsRSAPublicKey = utils.FlagDef{
		Name:  "rsa-pubkey",
		Usage: "KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding)",
	}
)
