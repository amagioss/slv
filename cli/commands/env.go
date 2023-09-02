package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shibme/slv/environment"
	"github.com/spf13/cobra"
)

func EnvCmd() *cobra.Command {
	env := &cobra.Command{
		Use:   "env",
		Short: "Managing environments",
		Long:  `Manage environments in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	env.AddCommand(NewEnvCommand())
	return env
}

func NewEnvCommand() *cobra.Command {
	envCreate := &cobra.Command{
		Use:   "new",
		Short: "Creates a service environment",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			tags, err := cmd.Flags().GetStringSlice("tags")
			if err != nil {
				PrintErrorAndExit(err)
				os.Exit(1)
			}
			env, privKey, _ := environment.New(name, email, environment.SERVICE)
			env.AddTags(tags...)
			envDef, err := env.ToEnvDef()
			if err != nil {
				PrintErrorAndExit(err)
				os.Exit(1)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "Secret Key:\t", privKey)
			fmt.Fprintln(w, "Public Key:\t", env.PublicKey)
			fmt.Fprintln(w, "Name:\t", env.Name)
			fmt.Fprintln(w, "Email:\t", env.Email)
			fmt.Fprintln(w, "Tags:\t", env.Tags)
			fmt.Fprintln(w, "Environment Definition:\t", envDef)
			w.Flush()
			os.Exit(0)
		},
	}

	// Adding the flags
	envCreate.Flags().StringP("name", "n", "", "Name of the environment")
	envCreate.Flags().StringP("email", "e", "", "Email for the environment")
	envCreate.Flags().StringSliceP("tags", "t", []string{}, "Tags for the environment")

	// Marking the flags as required
	envCreate.MarkFlagRequired("name")
	envCreate.MarkFlagRequired("email")
	return envCreate
}
