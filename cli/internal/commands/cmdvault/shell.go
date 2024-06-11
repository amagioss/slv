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
	"oss.amagi.com/slv"
	"oss.amagi.com/slv/cli/internal/commands/utils"
	"oss.amagi.com/slv/core/config"
)

func vaultShellCommand() *cobra.Command {
	if vaultShellCmd != nil {
		return vaultShellCmd
	}
	vaultShellCmd = &cobra.Command{
		Use:     "shell",
		Aliases: []string{"session", "sh", "venv", "vitualenv"},
		Short:   "Opens a shell with secrets loaded as environment variables",
		Run: func(cmd *cobra.Command, args []string) {
			shell := os.Getenv("SHELL")
			if shell == "" {
				if runtime.GOOS == "windows" {
					shell = "cmd"
				} else {
					utils.ExitOnErrorWithMessage("Not a supported shell")
				}
			}
			envSecretKey, err := slv.GetSecretKey()
			if err != nil {
				utils.ExitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			err = vault.Unlock(*envSecretKey)
			if err != nil {
				utils.ExitOnError(err)
			}
			secrets, err := vault.GetAllSecrets()
			if err != nil {
				utils.ExitOnError(err)
			}
			prefix := cmd.Flag(secretNamePrefixFlag.Name).Value.String()
			slvShell := exec.Command(shell)
			for _, envar := range os.Environ() {
				if !strings.HasPrefix(envar, "SLV_ENV_SECRET_") {
					slvShell.Env = append(slvShell.Env, envar)
				}
			}
			for name, secret := range secrets {
				if prefix != "" {
					name = prefix + name
				}
				slvShell.Env = append(slvShell.Env, name+"="+string(secret))
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
		},
	}
	vaultShellCmd.Flags().StringP(secretNamePrefixFlag.Name, secretNamePrefixFlag.Shorthand, "", secretNamePrefixFlag.Usage)
	return vaultShellCmd
}
