package cmdprofile

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
)

var (
	profileCmd          *cobra.Command
	profileNewCmd       *cobra.Command
	profileListCmd      *cobra.Command
	profileSetActiveCmd *cobra.Command
	profileDelCmd       *cobra.Command
	profileSyncCmd      *cobra.Command
)

var (
	profileNameFlag = utils.FlagDef{
		Name:      "name",
		Shorthand: "n",
		Usage:     "Profile name",
	}

	profileReadOnlyFlag = utils.FlagDef{
		Name:  "read-only",
		Usage: "Set profile as read-only",
	}

	profileSyncInterval = utils.FlagDef{
		Name:  "sync-interval",
		Usage: "Profile sync interval",
	}
)
