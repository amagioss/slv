package cmdenv

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
)

func ShowEnv(env environments.Environment, includeEDS, omitBindingFromEDS bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Public Key:\t", env.PublicKey)
	fmt.Fprintln(w, "Name:\t", env.Name)
	fmt.Fprintln(w, "Email:\t", env.Email)
	fmt.Fprintln(w, "Tags:\t", env.Tags)
	fmt.Fprintln(w, "Type:\t", env.EnvType)
	if env.SecretBinding != "" {
		fmt.Fprintln(w, "Secret Binding:\t", env.SecretBinding)
	}
	if includeEDS {
		if envDef, err := env.ToDefStr(omitBindingFromEDS); err == nil {
			fmt.Fprintln(w, "------------------------------------------------------------")
			fmt.Fprintln(w, "Env Definition:\t", color.CyanString(envDef))
		}
	}
	w.Flush()
}

func envShowCommand() *cobra.Command {
	if envShowCmd == nil {
		envShowCmd = &cobra.Command{
			Use:     "show",
			Aliases: []string{"describe", "view", "display"},
			Short:   "Shows the requested environment",
			Run: func(cmd *cobra.Command, args []string) {
				envdef, _ := cmd.Flags().GetString(envDefFlag.Name)
				if env, err := environments.FromDefStr(envdef); err == nil {
					ShowEnv(*env, true, false)
				} else {
					utils.ExitOnError(err)
				}
			},
		}
		envShowCmd.Flags().StringP(envDefFlag.Name, envDefFlag.Shorthand, "", envDefFlag.Usage)
		envShowCmd.MarkFlagRequired(envDefFlag.Name)
		envShowCmd.AddCommand(envShowRootCommand())
		envShowCmd.AddCommand(envShowSelfCommand())
		envShowCmd.AddCommand(envShowK8sCommand())
	}
	return envShowCmd
}

func envShowRootCommand() *cobra.Command {
	if envShowRootCmd == nil {
		envShowRootCmd = &cobra.Command{
			Use:   "root",
			Short: "Shows the root environment from the active profile",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetActiveProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				env, err := profile.GetRoot()
				if err != nil {
					utils.ExitOnError(err)
				}
				if env == nil {
					fmt.Println("No environment registered as root.")
				} else {
					ShowEnv(*env, true, false)
				}
			},
		}
	}
	return envShowRootCmd
}

func envShowSelfCommand() *cobra.Command {
	if envShowSelfCmd == nil {
		envShowSelfCmd = &cobra.Command{
			Use:     "self",
			Aliases: []string{"me"},
			Short:   "Shows the current user environment if registered in the host",
			Run: func(cmd *cobra.Command, args []string) {
				env := environments.GetSelf()
				if env == nil {
					fmt.Println("No environment registered as self.")
				} else {
					ShowEnv(*env, true, true)
				}
			},
		}
	}
	return envShowSelfCmd
}

func envShowK8sCommand() *cobra.Command {
	if envShowK8sCmd == nil {
		envShowK8sCmd = &cobra.Command{
			Use:     "k8s",
			Aliases: []string{"k8s-cluster"},
			Short:   "Shows the environment registered with the accessible k8s cluster",
			Run: func(cmd *cobra.Command, args []string) {
				name, address, user, err := session.GetK8sClusterInfo()
				if err != nil {
					utils.ExitOnError(err)
				}
				pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
				pk, err := session.GetPublicKeyFromK8s(config.AppNameLowerCase, pq)
				if err != nil {
					utils.ExitOnError(err)
				}
				var env *environments.Environment
				profile, err := profiles.GetActiveProfile()
				if err == nil {
					env, _ = profile.GetEnv(pk)
				}
				if env == nil {
					fmt.Printf("Public Key: %s\n", pk)
				} else {
					ShowEnv(*env, false, false)
				}
				fmt.Println("\nK8s Cluster Info:")
				fmt.Printf("Name   : %s\n", name)
				fmt.Printf("Address: %s\n", address)
				fmt.Printf("User   : %s\n", user)
			},
		}
		envShowK8sCmd.Flags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	}
	return envShowK8sCmd
}
