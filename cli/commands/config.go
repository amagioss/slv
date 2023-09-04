package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/configs"
	"github.com/spf13/cobra"
)

func configCommand() *cobra.Command {
	if configCmd != nil {
		return configCmd
	}
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage configs",
		Long:  `Manage configs in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	configCmd.AddCommand(configNewCommand())
	configCmd.AddCommand(configSetCommand())
	configCmd.AddCommand(configListCommand())
	return configCmd
}

func configNewCommand() *cobra.Command {
	if configNewCmd != nil {
		return configNewCmd
	}
	configNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new config",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			err := configs.New(name)
			if err == nil {
				fmt.Println("Created config: ", color.GreenString(name))
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		},
	}
	configNewCmd.Flags().StringP("name", "n", "", "Name for the config")
	configNewCmd.MarkFlagRequired("name")
	return configNewCmd
}

func configListCommand() *cobra.Command {
	if configListCmd != nil {
		return configListCmd
	}
	configListCmd = &cobra.Command{
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
						fmt.Println(color.GreenString(configName), "[*]")
					} else {
						fmt.Println(configName)
					}
				}
			}
		},
	}
	return configListCmd
}

func configSetCommand() *cobra.Command {
	if configSetCmd != nil {
		return configSetCmd
	}
	configSetCmd = &cobra.Command{
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
					fmt.Printf("Successfully set %s as default config\n", color.GreenString(name))
					os.Exit(0)
				}
			}
			fmt.Printf("Config %s not found\n", name)
			os.Exit(1)
		},
	}
	configSetCmd.Flags().StringP("name", "n", "", "Name of the config to be set as default")
	configSetCmd.MarkFlagRequired("name")
	return configSetCmd
}
