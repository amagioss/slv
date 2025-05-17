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
)

var (
	profileNameFlag = utils.FlagDef{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Profile name",
	}

	profileUpdateInterval = utils.FlagDef{
		Name:  "update-interval",
		Usage: "Interval in seconds to check for remote updates",
	}
)
