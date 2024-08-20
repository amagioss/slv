package main

import (
	"os"
	"strings"

	"oss.amagi.com/slv/internal/cli"
	"oss.amagi.com/slv/internal/k8s/job"
	"oss.amagi.com/slv/internal/k8s/operator"
)

const (
	slvModeK8sOperator = "K8S_OPERATOR"
	slvModeK8sJob      = "K8S_JOB"
)

func main() {
	switch strings.ToUpper(os.Getenv("SLV_MODE")) {
	case slvModeK8sOperator:
		operator.Run()
	case slvModeK8sJob:
		job.Run()
	default:
		cli.Run()
	}
}
