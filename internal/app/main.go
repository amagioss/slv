package main

import (
	"os"
	"strings"

	"slv.sh/slv/internal/cli"
	"slv.sh/slv/internal/k8s/job"
	"slv.sh/slv/internal/k8s/operator"
)

func main() {
	switch strings.ToLower(os.Getenv("SLV_MODE")) {
	case "k8s_operator", "k8s-operator":
		operator.Run()
	case "k8s_job", "k8s-job":
		job.Run()
	default:
		cli.Run()
	}
}
