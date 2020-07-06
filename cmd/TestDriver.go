package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"gopkg.in/yaml.v2"
)

// ITestDriver Interface for intiating test steps
type ITestDriver interface {
	createTestPod()
	copyArtifactsToTestPod()
	runTest()
	cleanUpResources()
}

// TestDriver Initiates test steps
type TestDriver struct {
	podTemplatePath string
	namespace       string
	localDirectory  string
	remoteDirectory string
	entryPoint      string
	failSilently    bool
	timeout         string
}

func (t *TestDriver) createTestPod() {
	log.Println("Applying objects...")
	runKubectl(false, "apply", "-f", t.podTemplatePath, "-n", t.namespace)
	runKubectl(false, "wait", "--for=condition=Ready", fmt.Sprintf("pod/%s", getPodName(t.podTemplatePath)), "-n", t.namespace)
}

func (t *TestDriver) copyArtifactsToTestPod() {
	runKubectl(false, "cp", t.localDirectory, fmt.Sprintf("%s:%s", getPodName(t.podTemplatePath), t.remoteDirectory), "-n", t.namespace)
}

func (t *TestDriver) runTest() {
	fmt.Println("Running test...")
	runKubectl(true, "exec", "-it", "-n", t.namespace, getPodName(t.podTemplatePath), "--", "/bin/sh", "-c", t.entryPoint)
}

func (t *TestDriver) cleanUpResources() {
	log.Printf("Testpod will be deleted in: %s", t.timeout)

	// Leave time for test pod to exist for inspection and testing before being deleted
	duration, _ := time.ParseDuration(t.timeout)
	time.Sleep(duration)

	log.Println("Cleaning up test pod...")
	runKubectl(false, "delete", "pod", getPodName(t.podTemplatePath), "-n", t.namespace)
}

func runKubectl(failSilently bool, kubectlOptions ...string) {
	cmd := exec.Command("kubectl", kubectlOptions...)
	stdout, err := cmd.Output()
	if err != nil && !failSilently {
		log.Fatalf("Test entrypoint returned non-zero exit code: %s", err)
	}
	log.Println(string(stdout))
}

type podtemplate struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}

func getPodName(podTemplatePath string) string {
	yamlFile, _ := ioutil.ReadFile(podTemplatePath) //handle err in future
	pod := &podtemplate{}
	yaml.Unmarshal(yamlFile, pod)
	return pod.Metadata.Name
}
