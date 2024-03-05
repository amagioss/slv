package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv"
	"savesecrets.org/slv/cli/internal/commands/cmdenv"
	"savesecrets.org/slv/cli/internal/commands/utils"
)

func vaultShareCommand() *cobra.Command {
	if vaultShareCmd != nil {
		return vaultShareCmd
	}
	vaultShareCmd = &cobra.Command{
		Use:   "share",
		Short: "Shares a vault with another environment or group",
		Run: func(cmd *cobra.Command, args []string) {
			envSecretKey, err := slv.GetSecretKey()
			if err != nil {
				utils.ExitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			query := cmd.Flag(cmdenv.EnvSearchFlag.Name).Value.String()
			shareWithSelf, _ := cmd.Flags().GetBool(cmdenv.EnvSelfFlag.Name)
			publicKeys, _, err := getPublicKeys(publicKeyStrings, query, shareWithSelf)
			if err != nil {
				utils.ExitOnError(err)
			}
			vault, err := getVault(vaultFile)
			if err == nil {
				err = vault.Unlock(*envSecretKey)
				if err == nil {
					for _, publicKey := range publicKeys {
						if _, err = vault.Share(publicKey); err != nil {
							break
						}
					}
					if err == nil {
						fmt.Println("Shared vault:", color.GreenString(vaultFile))
						utils.SafeExit()
					}
				}
			}
			utils.ExitOnError(err)
		},
	}
	vaultShareCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.Name, vaultAccessPublicKeysFlag.Shorthand, []string{}, vaultAccessPublicKeysFlag.Usage)
	vaultShareCmd.Flags().StringP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, "", cmdenv.EnvSearchFlag.Usage)
	vaultShareCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
	return vaultShareCmd
}
