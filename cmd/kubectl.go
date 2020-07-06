package cmd

import (
	"log"
	"os/exec"
)

// IKubectl Interface for running kubecctl commands
type IKubectl interface {
	runKubectl(failSilently bool, kubectlOptions ...string)
}

// Kubectl Runs kubectl commands
type Kubectl struct {
}

func (k *Kubectl) runKubectl(failSilently bool, kubectlOptions ...string) {
	cmd := exec.Command("kubectl", kubectlOptions...)
	stdout, err := cmd.Output()
	if err != nil && !failSilently {
		log.Fatalf("Test entrypoint returned non-zero exit code: %s", err)
	}
	log.Println(string(stdout))
}
