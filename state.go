package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// ensure checks if, depending on the status, the remote or local
// clusters are actually available (in case of local, launches it
//  if this is not the case)
func ensure(status, clocal, cremote string) error {
	switch status {
	case StatusOffline:
		fmt.Printf("Attempting to switch to %v, checking if local cluster is available\n", clocal)
		// TODO(mhausenblas): do a "minikube status" or "minishift status" and if not "Running", start it
	case StatusOnline:
		fmt.Printf("Attempting to switch to %v, checking if remote cluster is available \n", cremote)
		// TODO(mhausenblas): do a "kubectl get --raw /api" and if not ready, warn user
	}
	return nil
}

// capture queries the current state in the active namespace by exporting
// the state of deployments and services as a YAML doc
func capture(withstderr, verbose bool, namespace string) (string, error) {
	yamldoc := "---"
	deploys, err := kubectl(withstderr, verbose, "get", "--namespace="+namespace,
		"deployments", "--output=yaml")
	if err != nil {
		displayerr("Can't cuddle the cluster", err)
		return "", err
	}
	svcs, err := kubectl(withstderr, verbose, "get", "--namespace="+namespace,
		"services", "--output=yaml")
	if err != nil {
		displayerr("Can't cuddle the cluster", err)
		return "", err
	}
	yamldoc = deploys + "---\n" + svcs
	return yamldoc, nil
}

// dump stores a YAML doc in a file in:
// $StateCacheDir/$status/
func dump(status, yamldoc string) (string, error) {
	targetdir := filepath.Join(StateCacheDir, status)
	if _, err := os.Stat(targetdir); os.IsNotExist(err) {
		os.Mkdir(targetdir, os.ModePerm)
	}
	ts := time.Now().UnixNano()
	fn := filepath.Join(targetdir, "latest.yaml")
	err := ioutil.WriteFile(fn, []byte(yamldoc), 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", ts), nil
}

// restorefrom applies resources from the YAML doc at:
// $StateCacheDir/$state/$timestamp_of_last_state_dump
func restorefrom(withstderr, verbose bool, state, tsLast string) error {
	var err error
	statefile := filepath.Join(StateCacheDir, state, "latest.yaml")
	if verbose {
		fmt.Printf("Trying to restore state from %v/latest.yaml@%v\n", state, tsLast)
	}
	if _, err = os.Stat(statefile); !os.IsNotExist(err) {
		res, err := kubectl(withstderr, verbose, "apply", "--filename="+statefile)
		if err != nil {
			displayerr("Can't cuddle the cluster", err)
			return err
		}
		if verbose {
			fmt.Printf("Successfully restored state:\n%v\n", res)
		}
	}
	return err
}

// use switches over to provided context as in:
// `kubectl config use-context minikube`
func use(withstderr, verbose bool, context string) error {
	_, err := kubectl(withstderr, verbose, "config", "use-context", context)
	if err != nil {
		displayerr("Can't cuddle the cluster", err)
	}
	fmt.Printf("Now using context %v\n", context)
	return err
}
