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
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

func execVaultCommand(vaultFile, prefix, command string) {
	shell := false
	if command == "" {
		shell = true
		command = os.Getenv("SHELL")
		if command == "" {
			if runtime.GOOS == "windows" {
				command = "cmd"
			} else {
				utils.ExitOnErrorWithMessage("Not a supported shell")
			}
		}
	}
	commandArr := strings.Fields(command)
	if len(commandArr) == 1 {
		runVaultCommand(shell, vaultFile, prefix, commandArr[0])
	} else {
		runVaultCommand(shell, vaultFile, prefix, commandArr[0], commandArr[1:]...)
	}
}

func runVaultCommand(shell bool, vaultFile, prefix, command string, args ...string) {
	vault, err := vaults.Get(vaultFile)
	if err != nil {
		utils.ExitOnError(err)
	}
	envSecretKey, err := session.GetSecretKey()
	if err != nil {
		utils.ExitOnError(err)
	}
	err = vault.Unlock(envSecretKey)
	if err != nil {
		utils.ExitOnError(err)
	}
	secrets, err := vault.GetAllValues()
	if err != nil {
		utils.ExitOnError(err)
	}
	slvShell := exec.Command(command, args...)
	for _, envar := range os.Environ() {
		if !strings.HasPrefix(envar, "SLV_ENV_SECRET_") {
			slvShell.Env = append(slvShell.Env, envar)
		}
	}
	for name, value := range secrets {
		if prefix != "" {
			name = prefix + name
		}
		slvShell.Env = append(slvShell.Env, name+"="+string(value))
	}
	slvShell.Stdin = os.Stdin
	slvShell.Stdout = os.Stdout
	slvShell.Stderr = os.Stderr
	fullCommand := command
	if len(args) > 0 {
		fullCommand += " " + strings.Join(args, " ")
	}
	if shell {
		fmt.Printf("Initialized %s session with secrets loaded as environment variables from the vault %s.\n",
			config.AppNameUpperCase, color.CyanString(vaultFile))
	} else {
		fmt.Printf("Running command [%s] with secrets loaded as environment variables from the vault %s.\n",
			color.CyanString(fullCommand), color.CyanString(vaultFile))
	}
	if prefix != "" {
		fmt.Printf("Please note that the secret names are prefixed with %s\n", color.CyanString(prefix))
	}
	if err = slvShell.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if shell {
					fmt.Printf(color.RedString("%s session for the vault %s terminated with exit code %d\n"), config.AppNameUpperCase, vaultFile, status.ExitStatus())
				} else {
					fmt.Printf(color.RedString("Command [%s] terminated with exit code %d\n"), fullCommand, status.ExitStatus())
				}
				os.Exit(status.ExitStatus())
			}
		} else {
			utils.ExitOnError(err)
		}
	} else if shell {
		fmt.Printf("%s session for the vault %s ended successfully\n", config.AppNameUpperCase, color.CyanString(vaultFile))
	} else {
		fmt.Printf("Command [%s] ended successfully\n", color.CyanString(fullCommand))
	}
}

func vaultRunCommand() *cobra.Command {
	if vaultRunCmd == nil {
		vaultRunCmd = &cobra.Command{
			Use:     "run",
			Aliases: []string{"shell", "session", "venv", "vitualenv"},
			Short:   "Runs the given command or opens a shell with the vault items loaded as environment variables",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				command := cmd.Flag(vaultShellCommandFlag.Name).Value.String()
				prefix := cmd.Flag(varNamePrefixFlag.Name).Value.String()
				if len(args) > 0 {
					runVaultCommand(false, vaultFile, prefix, args[0], args[1:]...)
				} else {
					execVaultCommand(vaultFile, prefix, command)
				}
			},
		}
		vaultRunCmd.Flags().StringP(varNamePrefixFlag.Name, varNamePrefixFlag.Shorthand, "", varNamePrefixFlag.Usage)
		vaultRunCmd.Flags().StringP(vaultShellCommandFlag.Name, vaultShellCommandFlag.Shorthand, "", vaultShellCommandFlag.Usage)
	}
	return vaultRunCmd
}
