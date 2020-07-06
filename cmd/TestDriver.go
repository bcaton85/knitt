package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
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
	k               IKubectl
}

func (t *TestDriver) createTestPod() {
	log.Println("Applying objects...")
	t.k.runKubectl(false, "apply", "-f", t.podTemplatePath, "-n", t.namespace)
	t.k.runKubectl(false, "wait", "--for=condition=Ready", fmt.Sprintf("pod/%s", getPodName(t.podTemplatePath)), "-n", t.namespace)
}

func (t *TestDriver) copyArtifactsToTestPod() {
	t.k.runKubectl(false, "cp", t.localDirectory, fmt.Sprintf("%s:%s", getPodName(t.podTemplatePath), t.remoteDirectory), "-n", t.namespace)
}

func (t *TestDriver) runTest() {
	fmt.Println("Running test...")
	t.k.runKubectl(true, "exec", "-it", "-n", t.namespace, getPodName(t.podTemplatePath), "--", "/bin/sh", "-c", t.entryPoint)
}

func (t *TestDriver) cleanUpResources() {
	log.Printf("Testpod will be deleted in: %s", t.timeout)

	// Leave time for test pod to exist for inspection and testing before being deleted
	duration, _ := time.ParseDuration(t.timeout)
	time.Sleep(duration)

	log.Println("Cleaning up test pod...")
	t.k.runKubectl(false, "delete", "pod", getPodName(t.podTemplatePath), "-n", t.namespace)
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
