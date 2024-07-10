package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/profiles"
	k8sutils "oss.amagi.com/slv/internal/k8s/utils"
)

func envK8sCommand() *cobra.Command {
	if envK8sCmd != nil {
		return envK8sCmd
	}
	envK8sCmd = &cobra.Command{
		Use:     "k8s",
		Aliases: []string{"k8s-cluster"},
		Short:   "Shows the environment registered with the current k8s context",
		Run: func(cmd *cobra.Command, args []string) {
			name, address, user, err := k8sutils.GetClusterInfo()
			if err != nil {
				utils.ExitOnError(err)
			}
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			pk, err := k8sutils.GetPublicKeyFromK8s(config.AppNameLowerCase, pq)
			if err != nil {
				utils.ExitOnError(err)
			}
			var env *environments.Environment
			profile, err := profiles.GetDefaultProfile()
			if err == nil {
				env, _ = profile.GetEnv(pk)
			}
			if env == nil {
				fmt.Printf("Public Key: %s\n", pk)
			} else {
				utils.ShowEnv(*env, false, false)
			}
			fmt.Println("\nK8s Cluster Info:")
			fmt.Printf("Name   : %s\n", name)
			fmt.Printf("Address: %s\n", address)
			fmt.Printf("User   : %s\n", user)
		},
	}
	envK8sCmd.Flags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	return envK8sCmd
}
