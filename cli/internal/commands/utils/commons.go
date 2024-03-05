package utils

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"savesecrets.org/slv/core/environments"
)

func ExitOnError(err error) {
	fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
	ErroredExit()
}

func ExitOnErrorWithMessage(errMessage string) {
	fmt.Fprintln(os.Stderr, color.RedString(errMessage))
	ErroredExit()
}

func ErroredExit() {
	os.Exit(1)
}

func SafeExit() {
	os.Exit(0)
}

func ShowEnv(env environments.Environment, includeEDS, excludeBindingFromEds bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "ID (Public Key):\t", env.PublicKey)
	fmt.Fprintln(w, "Name:\t", env.Name)
	fmt.Fprintln(w, "Email:\t", env.Email)
	fmt.Fprintln(w, "Tags:\t", env.Tags)
	if env.SecretBinding != "" {
		fmt.Fprintln(w, "Secret Binding:\t", env.SecretBinding)
	}
	if includeEDS {
		secretBinding := env.SecretBinding
		if excludeBindingFromEds {
			env.SecretBinding = ""
		}
		if envDef, err := env.ToEnvDef(); err == nil {
			fmt.Fprintln(w, "------------------------------------------------------------")
			fmt.Fprintln(w, "Env Definition:\t", color.CyanString(envDef))
		}
		env.SecretBinding = secretBinding
	}
	w.Flush()
}
