package cmdvault

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"oss.amagi.com/slv"
	"oss.amagi.com/slv/cli/internal/commands/utils"
)

func vaultShellCommand() *cobra.Command {
	if vaultShellCmd != nil {
		return vaultShellCmd
	}
	vaultShellCmd = &cobra.Command{
		Use:     "shell",
		Aliases: []string{"venv", "sh", "getshell", "vitualenv"},
		Short:   "Opens a shell with secrets loaded as environment variables",
		Run: func(cmd *cobra.Command, args []string) {
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
			shell := os.Getenv("SHELL")
			if shell == "" {
				utils.ExitOnErrorWithMessage("Not a supported shell")
			}
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
			if err = slvShell.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
						os.Exit(status.ExitStatus())
					}
				} else {
					utils.ExitOnError(err)
				}
			}
		},
	}
	vaultShellCmd.Flags().StringP(secretNamePrefixFlag.Name, secretNamePrefixFlag.Shorthand, "", secretNamePrefixFlag.Usage)
	return vaultShellCmd
}
