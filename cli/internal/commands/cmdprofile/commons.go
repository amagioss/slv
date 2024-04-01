package cmdprofile

import (
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/cli/internal/commands/utils"
)

var (
	// Commands
	profileCmd     *cobra.Command
	profileNewCmd  *cobra.Command
	profileListCmd *cobra.Command
	profileSetCmd  *cobra.Command
	profileDelCmd  *cobra.Command
	profilePullCmd *cobra.Command
	profilePushCmd *cobra.Command
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
