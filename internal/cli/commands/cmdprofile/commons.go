package cmdprofile

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
)

var (
	profileCmd           *cobra.Command
	profileNewCmd        *cobra.Command
	profileListCmd       *cobra.Command
	profileSetCurrentCmd *cobra.Command
	profileDelCmd        *cobra.Command
	profilePullCmd       *cobra.Command
	profilePushCmd       *cobra.Command
)

var (
	profileNameFlag = utils.FlagDef{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Profile name",
	}

	profileGitURI = utils.FlagDef{
		Name:  "git",
		Usage: "Git URI to clone the profile from",
	}

	profileGitBranch = utils.FlagDef{
		Name:  "git-branch",
		Usage: "Git branch corresponding to the git URI",
	}
)
