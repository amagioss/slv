package main

import (
	"os"
	"strings"

	"oss.amagi.com/slv/internal/cli"
	"oss.amagi.com/slv/internal/k8s/job"
	"oss.amagi.com/slv/internal/k8s/operator"
)

const (
	slvModeEnvar       = "SLV_MODE"
	slvModeK8sOperator = "k8s_operator"
	slvModeK8sJob      = "k8s_job"
)

func main() {
	slvMode := strings.ToLower(os.Getenv(slvModeEnvar))
	switch slvMode {
	case slvModeK8sOperator:
		operator.Run()
	case slvModeK8sJob:
		job.Run()
	default:
		cli.Run()
	}
}
