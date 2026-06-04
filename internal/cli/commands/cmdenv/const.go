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
		Usage: "Add the environment to the active profile after creation",
	}

	envSetRootFlag = utils.FlagDef{
		Name:  "root",
		Usage: "Set the environment as the root environment of the active profile",
	}

	EnvSearchFlag = utils.FlagDef{
		Name:      "env-search",
		Shorthand: "s",
		Usage:     "Search query to filter environments by name, email, or tag",
	}

	showEnvDefFlag = utils.FlagDef{
		Name:  "show-env-def",
		Usage: "Include the Environment Definition String (EDS) in the output",
	}

	EnvSelfFlag = utils.FlagDef{
		Name:  "env-self",
		Usage: "Use the self environment (the one registered on this machine)",
	}

	envDefFlag = utils.FlagDef{
		Name:      "env-def",
		Shorthand: "e",
		Usage:     "Environment definition that begins with SLV_EDS_",
	}

	EnvPublicKeysFlag = utils.FlagDef{
		Name:      "env-pubkey",
		Shorthand: "k",
		Usage:     "Public key(s) of environments to grant access to",
	}

	EnvK8sFlag = utils.FlagDef{
		Name:  "env-k8s",
		Usage: "Use the environment registered with the current Kubernetes cluster",
	}
)
