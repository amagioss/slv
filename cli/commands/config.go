package commands

import (
	"fmt"
	"os"

	"github.com/shibme/slv/config"
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
	configCmd.AddCommand(setConfigCommand())
	// configCmd.AddCommand(deleteConfigCommand())
	configCmd.AddCommand(listConfigCommand())
	return configCmd
}

func createConfigCommand() *cobra.Command {
	configCreate := &cobra.Command{
		Use:   "new",
		Short: "Create a new config",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			err := config.New(name)
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
			configNames, err := config.List()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				currentConfigName, err := config.GetCurrentConfigName()
				if err != nil {
					PrintErrorAndExit(err)
					os.Exit(1)
				}
				for _, configName := range configNames {
					if configName == currentConfigName {
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

func setConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "set",
		Short: "Set a config as current config",
		Run: func(cmd *cobra.Command, args []string) {
			configNames, err := config.List()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			name, _ := cmd.Flags().GetString("name")
			for _, configName := range configNames {
				if configName == name {
					config.Set(name)
					fmt.Printf("Set %s as current config\n", name)
					os.Exit(0)
				}
			}
			fmt.Printf("Config %s not found\n", name)
			os.Exit(1)
		},
	}

	// Adding the flags
	command.Flags().StringP("name", "n", "", "Name of the config to be set as current")

	// Marking the flags as required
	command.MarkFlagRequired("name")
	return command
}

func addEnvToConfigCommand() *cobra.Command {
	addEnv := &cobra.Command{
		Use:     "addenv",
		Short:   "Adds an environment to a config",
		Aliases: []string{"envadd"},
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("env")
			err := config.New(name)
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
	addEnv.Flags().StringP("env", "e", "", "Environment defintion to be added")

	// Marking the flags as required
	addEnv.MarkFlagRequired("env")
	return addEnv
}
