# Kubernetes Native Integration Test Tool

## Description

A tool used to automated integration tests of services in kubernetes namespaces.
Offers reliability in that tests are ran against services in a more production like environment.

Short explanation:

A test pod is deployed into the target namespace containing the user made test scripts, those scripts are ran and the logs are sent to stdout. The test resources are then removed.

A detailed workflow:
- The given pod template is applied to the target namespace
- The pod is checked to be in the Ready state before continuing
- The test project or scripts are copied from the local machine to the remote pod
- The tests are ran with the given entrypoint and logs are streamed back to standard out
- The exit status is checked, if set to silent a non-zero exit code from the entrypoint will not cause the command to fail
- If a timeout value is given the command will pause for the specified amount of time
- After the timeout the pod resource deleted

## Usage

knitt run --entrypoint COMMAND --pod-template PATH_TO_TEMPLATE [flags]

Flags:

-e, --entry-point string        
Command to run to start test

-s, --fail-silently             
If set to true non-zero exit code of entrypoint will not cause the command to fail

-h, --help                      
help for run

-d, --local-directory string    
Local directory containing test project or scripts (default ".")

-n, --namespace string          
Kubernetes namespace to deploy test pod to (default "default")

-p, --pod-template string       
Path to yaml file containing pod template

-r, --remote-directory string   
Remote directory where test scripts will be copied to (default "/home")

-t, --timeout string            
Time that the test pod should be left alive for after tests have completed. Examples: 1h, 60s, 30m, 1h10m10s (default "0s")

