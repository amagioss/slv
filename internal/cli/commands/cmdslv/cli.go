package cmdslv

import "os"

func RunCLI() {
	if err := slvCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
