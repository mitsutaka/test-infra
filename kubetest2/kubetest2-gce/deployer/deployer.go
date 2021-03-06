/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package deployer implements the kubetest2 GKE deployer
package deployer

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/klog"
	"k8s.io/test-infra/kubetest2/pkg/types"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/pflag"
)

// Name is the name of the deployer
const Name = "gce"

type deployer struct {
	// generic parts
	commonOptions types.Options

	kubeconfigPath string
	RepoRoot       string `desc:"The path to the root of the local kubernetes/cloud-provider-gcp repo. Necessary to call certain scripts. Defaults to the current directory. If operating in legacy mode, this should be set to the local kubernetes/kubernetes repo."`
	GCPProject     string `desc:"GCP Project to create VMs in. Must be set."`
}

// New implements deployer.New for gce
func New(opts types.Options) (types.Deployer, *pflag.FlagSet) {
	d := &deployer{
		commonOptions:  opts,
		kubeconfigPath: filepath.Join(opts.ArtifactsDir(), "kubetest2-kubeconfig"),
	}

	flagSet, err := gpflag.Parse(d)
	if err != nil {
		klog.Fatalf("couldn't parse flagset for deployer struct: %s", err)
	}

	// register flags and return
	return d, flagSet
}

// assert that New implements types.NewDeployer
var _ types.NewDeployer = New

// assert that deployer implements types.Deployer
var _ types.Deployer = &deployer{}

func (d *deployer) Provider() string {
	return Name
}

func (d *deployer) IsUp() (up bool, err error) { return false, nil }

func (d *deployer) DumpClusterLogs() error { return nil }

func (d *deployer) Kubeconfig() (string, error) {
	_, err := os.Stat(d.kubeconfigPath)
	if err == nil {
		return d.kubeconfigPath, nil
	} else if os.IsNotExist(err) {
		return "", fmt.Errorf("kubeconfig does not exist at: %s", d.kubeconfigPath)
	}
	return "", fmt.Errorf("unknown error when checking for kubeconfig at %s: %s", d.kubeconfigPath, err)
}

func (d *deployer) Down() error { return nil }
