package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shibme/slv/core/configs"
	"github.com/shibme/slv/core/environments"
	"github.com/spf13/cobra"
)

func envCommand() *cobra.Command {
	if envCmd != nil {
		return envCmd
	}
	envCmd = &cobra.Command{
		Use:   "env",
		Short: "Environment operations",
		Long:  `Environment operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	envCmd.AddCommand(envNewCommand())
	envCmd.AddCommand(envListCommand())
	envCmd.AddCommand(envAddCommand())
	envCmd.AddCommand(envRootInitCommand())
	return envCmd
}

func envNewCommand() *cobra.Command {
	if envNewCmd != nil {
		return envNewCmd
	}
	envNewCmd = &cobra.Command{
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
			env, privKey, _ := environments.New(name, email, environments.SERVICE)
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

			// Adding env to a specified config
			configName := cmd.Flag("add-to-config").Value.String()
			var cfg *configs.Config
			if configName != "" {
				cfg, err = configs.GetConfig(configName)
			} else {
				cfg, err = configs.GetDefaultConfig()
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = cfg.AddEnv(envDef)
			if err != nil {
				PrintErrorAndExit(err)
			}
			os.Exit(0)
		},
	}
	envNewCmd.Flags().StringP("name", "n", "", "Name of the environment")
	envNewCmd.Flags().StringP("email", "e", "", "Email for the environment")
	envNewCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags for the environment")
	envNewCmd.Flags().StringP("add-to-config", "c", "", "Config to add the environment to")
	envNewCmd.MarkFlagRequired("name")
	return envNewCmd
}

func envAddCommand() *cobra.Command {
	if envAddCmd != nil {
		return envAddCmd
	}
	envAddCmd = &cobra.Command{
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
	envAddCmd.Flags().StringP("config", "c", "", "Name of the config to add the environment to")
	envAddCmd.Flags().StringP("envdef", "e", "", "Environment defintion to be added")
	envAddCmd.MarkFlagRequired("envdef")
	return envAddCmd
}

func envRootInitCommand() *cobra.Command {
	if envRootInitCmd != nil {
		return envRootInitCmd
	}
	envRootInitCmd = &cobra.Command{
		Use:   "initroot",
		Short: "Initializes the root environment in a config",
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
			privKey, err := cfg.InitRoot()
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println("Root environment initialized with secret key:", privKey)
		},
	}
	envRootInitCmd.Flags().StringP("config", "c", "", "Name of the config to initialize root environment")
	return envRootInitCmd
}

func envListCommand() *cobra.Command {
	if envListCmd != nil {
		return envListCmd
	}
	envListCmd = &cobra.Command{
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
			var envs []*environments.Environment
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
	envListCmd.Flags().StringP("config", "c", "", "Environment defintion to be added")
	envListCmd.Flags().StringP("search", "s", "", "Search query to lookup envionments")
	return envListCmd
}
