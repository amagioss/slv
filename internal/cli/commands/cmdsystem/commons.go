package cmdsystem

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
)

var (
	systemCmd      *cobra.Command
	systemResetCmd *cobra.Command
)

var (
	yesFlag = utils.FlagDef{
		Name:      "yes",
		Shorthand: "y",
		Usage:     "Confirm action",
	}
)
