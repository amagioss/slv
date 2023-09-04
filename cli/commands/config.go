package commands

import (
	"fmt"
	"os"

	"github.com/shibme/slv/configs"
	"github.com/spf13/cobra"
)

func ConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configs",
		Long:  `Manage configs in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	configCmd.AddCommand(createConfigCommand())
	configCmd.AddCommand(setDefaultConfigCommand())
	// configCmd.AddCommand(deleteConfigCommand())
	configCmd.AddCommand(listConfigCommand())
	configCmd.AddCommand(configEnvCmd())
	return configCmd
}

func createConfigCommand() *cobra.Command {
	configCreate := &cobra.Command{
		Use:   "new",
		Short: "Create a new config",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			err := configs.New(name)
			if err == nil {
				fmt.Println("Created config: ", green, name)
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		},
	}

	// Adding the flags
	configCreate.Flags().StringP("name", "n", "", "Name for the config")

	// Marking the flags as required
	configCreate.MarkFlagRequired("name")
	return configCreate
}

func listConfigCommand() *cobra.Command {
	configList := &cobra.Command{
		Use:   "list",
		Short: "Lists all configs",
		Run: func(cmd *cobra.Command, args []string) {
			configNames, err := configs.List()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				defaultConfigName, _ := configs.GetDefaultConfigName()
				for _, configName := range configNames {
					if configName == defaultConfigName {
						fmt.Println(green + configName + reset + " [*]")
					} else {
						fmt.Println(configName)
					}
				}
			}
		},
	}
	return configList
}

func setDefaultConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "default",
		Aliases: []string{"set-default"},
		Short:   "Set a config as default config",
		Run: func(cmd *cobra.Command, args []string) {
			configNames, err := configs.List()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			name, _ := cmd.Flags().GetString("name")
			for _, configName := range configNames {
				if configName == name {
					configs.SetDefault(name)
					fmt.Printf("Successfully set %s%s%s as default config\n", green, name, reset)
					os.Exit(0)
				}
			}
			fmt.Printf("Config %s not found\n", name)
			os.Exit(1)
		},
	}

	// Adding the flags
	command.Flags().StringP("name", "n", "", "Name of the config to be set as default")

	// Marking the flags as required
	command.MarkFlagRequired("name")
	return command
}

func configEnvCmd() *cobra.Command {
	env := &cobra.Command{
		Use:   "env",
		Short: "Managing environments",
		Long:  `Manage environments in an SLV Config`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	env.AddCommand(addEnvToConfig())
	env.AddCommand(listConfigEnvs())
	return env
}
