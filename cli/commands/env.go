package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shibme/slv/configs"
	"github.com/shibme/slv/environment"
	"github.com/spf13/cobra"
)

func EnvCommand() *cobra.Command {
	env := &cobra.Command{
		Use:   "env",
		Short: "Environment operations",
		Long:  `Environment operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	env.AddCommand(newEnvCommand())
	env.AddCommand(listConfigEnvs())
	env.AddCommand(addEnvToConfig())
	return env
}

func newEnvCommand() *cobra.Command {
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

func addEnvToConfig() *cobra.Command {
	addEnv := &cobra.Command{
		Use:   "add",
		Short: "Adds an environment to a config",
		Run: func(cmd *cobra.Command, args []string) {
			envdef := cmd.Flag("envdef").Value.String()
			configName := cmd.Flag("config").Value.String()
			var cfg *configs.Config
			var err error
			if configName != "" {
				cfg, err = configs.GetConfig(configName)
			} else {
				cfg, err = configs.GetDefaultConfig()
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = cfg.AddEnv(envdef)
			if err != nil {
				PrintErrorAndExit(err)
			}
		},
	}

	// Adding the flags
	addEnv.Flags().StringP("envdef", "e", "", "Environment defintion to be added")
	addEnv.Flags().StringP("config", "c", "", "Name of the config to add the environment to")

	// Marking the flags as required
	addEnv.MarkFlagRequired("envdef")
	return addEnv
}

func listConfigEnvs() *cobra.Command {
	listEnv := &cobra.Command{
		Use:   "list",
		Short: "Lists environments from config",
		Run: func(cmd *cobra.Command, args []string) {
			configName := cmd.Flag("config").Value.String()
			var cfg *configs.Config
			var err error
			if configName != "" {
				cfg, err = configs.GetConfig(configName)
			} else {
				cfg, err = configs.GetDefaultConfig()
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			envManifest, err := cfg.GetEnvManifest()
			if err != nil {
				PrintErrorAndExit(err)
			}
			query := cmd.Flag("search").Value.String()
			var envs []*environment.Environment
			if query != "" {
				envs = envManifest.SearchEnv(query)
			} else {
				envs = envManifest.ListEnv()
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			for _, env := range envs {
				fmt.Fprintln(w, env.Id()+":")
				fmt.Fprintln(w, "Public Key:\t", env.PublicKey)
				fmt.Fprintln(w, "Name:\t", env.Name)
				fmt.Fprintln(w, "Email:\t", env.Email)
				fmt.Fprintln(w, "Tags:\t", env.Tags)
				fmt.Fprintln(w)
			}
			w.Flush()
			os.Exit(0)

		},
	}

	// Adding the flags
	listEnv.Flags().StringP("config", "c", "", "Environment defintion to be added")
	listEnv.Flags().StringP("search", "s", "", "Search query to lookup envionments")

	return listEnv
}
