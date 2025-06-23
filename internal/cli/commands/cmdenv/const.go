package cmdenv

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
)

var (
	envCmd                     *cobra.Command
	envNewCmd                  *cobra.Command
	envNewServiceCmd           *cobra.Command
	envNewDirectServicetextCmd *cobra.Command
	envNewUserCmd              *cobra.Command
	envAddCmd                  *cobra.Command
	envListCmd                 *cobra.Command
	envDeleteCmd               *cobra.Command
	envSetSelfSCmd             *cobra.Command
	envShowCmd                 *cobra.Command
	envShowRootCmd             *cobra.Command
	envShowSelfCmd             *cobra.Command
	envShowK8sCmd              *cobra.Command
)

var (
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
		Usage: "Adds environment to active profile",
	}

	envSetRootFlag = utils.FlagDef{
		Name:  "root",
		Usage: "Set the given environment as root",
	}

	EnvSearchFlag = utils.FlagDef{
		Name:      "env-search",
		Shorthand: "s",
		Usage:     "Searches query to filter environments",
	}

	showEnvDefFlag = utils.FlagDef{
		Name:  "show-env-def",
		Usage: "Shows the environment definition in the output",
	}

	EnvSelfFlag = utils.FlagDef{
		Name:  "env-self",
		Usage: "References to the self environment (the local environment where the command is executed)",
	}

	envDefFlag = utils.FlagDef{
		Name:      "env-def",
		Shorthand: "e",
		Usage:     "Environment definition that begins with SLV_EDS_",
	}

	EnvPublicKeysFlag = utils.FlagDef{
		Name:      "env-pubkey",
		Shorthand: "k",
		Usage:     "Public keys of environments that can access the vault",
	}

	EnvK8sFlag = utils.FlagDef{
		Name:  "env-k8s",
		Usage: "Shares vault access with the accessible k8s cluster",
	}
)
