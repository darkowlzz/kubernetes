/*
Copyright 2017 The Kubernetes Authors.

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

package phases

import (
	"fmt"

	"github.com/spf13/cobra"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmapiext "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	etcdphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/etcd"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/util/normalizer"
)

var (
	etcdLocalLongDesc = fmt.Sprintf(normalizer.LongDesc(`
		Generates the static Pod manifest file for a local, single-node etcd instance and saves it to %s file.
		`+cmdutil.AlphaDisclaimer), kubeadmconstants.GetStaticPodFilepath(kubeadmconstants.Etcd, kubeadmconstants.GetStaticPodDirectory()))

	etcdLocalExample = normalizer.Examples(`
		# Generates the static Pod manifest file for etcd, functionally 
		# equivalent to what generated by kubeadm init.
		kubeadm alpha phase etcd local

		#  Generates the static Pod manifest file for etcd.
		kubeadm alpha phase etcd local --config masterconfiguration.yaml
		`)
)

// NewCmdEtcd return main command for Etcd phase
func NewCmdEtcd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "etcd",
		Short: "Generates static Pod manifest file for etcd.",
		Long:  cmdutil.MacroCommandLongDescription,
	}

	manifestPath := kubeadmconstants.GetStaticPodDirectory()
	cmd.AddCommand(getEtcdSubCommands(manifestPath, "")...)
	return cmd
}

// getEtcdSubCommands returns sub commands for etcd phase
func getEtcdSubCommands(outDir, defaultKubernetesVersion string) []*cobra.Command {

	cfg := &kubeadmapiext.MasterConfiguration{}

	// This is used for unit testing only...
	// If we wouldn't set this to something, the code would dynamically look up the version from the internet
	// By setting this explicitely for tests workarounds that
	if defaultKubernetesVersion != "" {
		cfg.KubernetesVersion = defaultKubernetesVersion
	}

	// Default values for the cobra help text
	legacyscheme.Scheme.Default(cfg)

	var cfgPath string
	var subCmds []*cobra.Command

	properties := struct {
		use      string
		short    string
		long     string
		examples string
		cmdFunc  func(outDir string, cfg *kubeadmapi.MasterConfiguration) error
	}{
		use:      "local",
		short:    "Generates the static Pod manifest file for a local, single-node etcd instance",
		long:     etcdLocalLongDesc,
		examples: etcdLocalExample,
		cmdFunc:  etcdphase.CreateLocalEtcdStaticPodManifestFile,
	}

	// Creates the UX Command
	cmd := &cobra.Command{
		Use:     properties.use,
		Short:   properties.short,
		Long:    properties.long,
		Example: properties.examples,
		Run:     runCmdPhase(properties.cmdFunc, &outDir, &cfgPath, cfg),
	}

	// Add flags to the command
	cmd.Flags().StringVar(&cfg.CertificatesDir, "cert-dir", cfg.CertificatesDir, `The path where certificates are stored`)
	cmd.Flags().StringVar(&cfgPath, "config", cfgPath, "Path to kubeadm config file (WARNING: Usage of a configuration file is experimental)")

	subCmds = append(subCmds, cmd)

	return subCmds
}