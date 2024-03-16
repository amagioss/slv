package cmdvault

import "github.com/spf13/cobra"

func vaultSecretsCommand() *cobra.Command {
	if vaultSecretsCmd != nil {
		return vaultSecretsCmd
	}
	vaultSecretsCmd = &cobra.Command{
		Use:     "secret",
		Aliases: []string{"secrets"},
		Short:   "Manage secrets in the vault",
		Long:    `Add, remove, and get secrets from the vault`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	vaultSecretsCmd.AddCommand(vaultPutCommand())
	vaultSecretsCmd.AddCommand(vaultGetCommand())
	return vaultSecretsCmd
}
