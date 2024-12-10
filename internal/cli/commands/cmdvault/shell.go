package cmdvault

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/secretkey"
)

func execVaultShellCommand(vaultFile, prefix, command string) {
	if command == "" {
		command = os.Getenv("SHELL")
		if command == "" {
			if runtime.GOOS == "windows" {
				command = "cmd"
			} else {
				utils.ExitOnErrorWithMessage("Not a supported shell")
			}
		}
	}
	vault, err := getVault(vaultFile)
	if err != nil {
		utils.ExitOnError(err)
	}
	envSecretKey, err := secretkey.Get()
	if err != nil {
		utils.ExitOnError(err)
	}
	err = vault.Unlock(envSecretKey)
	if err != nil {
		utils.ExitOnError(err)
	}
	secrets, err := vault.List(true)
	if err != nil {
		utils.ExitOnError(err)
	}
	slvShell := exec.Command(command)
	for _, envar := range os.Environ() {
		if !strings.HasPrefix(envar, "SLV_ENV_SECRET_") {
			slvShell.Env = append(slvShell.Env, envar)
		}
	}
	for name, data := range secrets {
		if prefix != "" {
			name = prefix + name
		}
		slvShell.Env = append(slvShell.Env, name+"="+string(data.Value()))
	}
	slvShell.Stdin = os.Stdin
	slvShell.Stdout = os.Stdout
	slvShell.Stderr = os.Stderr
	fmt.Printf("Initializing %s session with secrets loaded into environment variables from the vault %s...\n",
		config.AppNameUpperCase, color.CyanString(vaultFile))
	if prefix != "" {
		fmt.Printf("Please note that the secret names are prefixed with %s\n", color.CyanString(prefix))
	}
	if err = slvShell.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				fmt.Printf(color.RedString("%s session for the vault %s terminated with exit code %d\n"), config.AppNameUpperCase, vaultFile, status.ExitStatus())
				os.Exit(status.ExitStatus())
			}
		} else {
			utils.ExitOnError(err)
		}
	} else {
		fmt.Printf("%s session for the vault %s has been closed\n", config.AppNameUpperCase, color.CyanString(vaultFile))
	}
}

func vaultShellCommand() *cobra.Command {
	if vaultShellCmd != nil {
		return vaultShellCmd
	}
	vaultShellCmd = &cobra.Command{
		Use:     "shell",
		Aliases: []string{"session", "sh", "venv", "vitualenv"},
		Short:   "Opens a shell with the vault items loaded as environment variables [optinally run a command]",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			command := cmd.Flag(vaultShellCommandFlag.Name).Value.String()
			prefix := cmd.Flag(varNamePrefixFlag.Name).Value.String()
			execVaultShellCommand(vaultFile, prefix, command)
		},
	}
	vaultShellCmd.Flags().StringP(varNamePrefixFlag.Name, varNamePrefixFlag.Shorthand, "", varNamePrefixFlag.Usage)
	vaultShellCmd.Flags().StringP(vaultShellCommandFlag.Name, vaultShellCommandFlag.Shorthand, "", vaultShellCommandFlag.Usage)
	return vaultShellCmd
}
