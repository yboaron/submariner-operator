/*
SPDX-License-Identifier: Apache-2.0

Copyright Contributors to the Submariner project.

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

package subctl

import (
	"github.com/spf13/cobra"
	"github.com/submariner-io/submariner-operator/internal/execute"
	"github.com/submariner-io/submariner-operator/pkg/diagnose"
)

var (
	verboseOutput      bool
)

var (
	diagnoseCmd = &cobra.Command{
		Use:   "diagnose",
		Short: "Run diagnostic checks on the Submariner deployment and report any issues",
		Long:  "This command runs various diagnostic checks on the Submariner deployment and reports any issues",
    }
    cniCmd = &cobra.Command{
		Use:   "cni",
		Short: "Check the CNI network plugin",
		Long:  "This command checks if the detected CNI network plugin is supported by Submariner.",
		Run: func(command *cobra.Command, args []string) {
			execute.OnMultiCluster(restConfigProducer, diagnose.CheckCNIConfig)
		},
    }
	connectionsCmd = &cobra.Command{
		Use:   "connections",
		Short: "Check the Gateway connections",
		Long:  "This command checks that the Gateway connections to other clusters are all established",
		Run: func(command *cobra.Command, args []string) {
			execute.OnMultiCluster(restConfigProducer, diagnose.CheckConnections)
		},
	}
	deploymentCmd = &cobra.Command{
		Use:   "deployment",
		Short: "Check the Submariner deployment",
		Long:  "This command checks that the Submariner components are properly deployed and running with no overlapping CIDRs.",
		Run: func(command *cobra.Command, args []string) {
			execute.OnMultiCluster(restConfigProducer, diagnose.CheckDeployments)
		},
	}
	versionCmd = &cobra.Command{
		Use:   "k8s-version",
		Short: "Check the Kubernetes version",
		Long:  "This command checks if Submariner can be deployed on the Kubernetes version.",
		Run: func(command *cobra.Command, args []string) {
			execute.OnMultiCluster(restConfigProducer, diagnose.CheckK8sVersion)
		},
	}
	kpModeCmd = &cobra.Command{
		Use:   "kube-proxy-mode",
		Short: "Check the kube-proxy mode",
		Long:  "This command checks if the kube-proxy mode is supported by Submariner.",
		Run: func(command *cobra.Command, args []string) {
			execute.OnMultiCluster(restConfigProducer, diagnose.CheckKubeProxyMode)
		},
	}
	allCmd = &cobra.Command{
		Use:   "all",
		Short: "Run all diagnostic checks (except those requiring two kubecontexts)",
		Long:  "This command runs all diagnostic checks (except those requiring two kubecontexts) and reports any issues",
		Run: func(command *cobra.Command, args []string) {
			execute.OnMultiCluster(restConfigProducer, diagnose.All)
		},
	}
)

func init() {
	restConfigProducer.AddKubeConfigFlag(diagnoseCmd)
	restConfigProducer.AddInClusterConfigFlag(diagnoseCmd)
	addNamespaceFlag(kpModeCmd)
	rootCmd.AddCommand(diagnoseCmd)
	diagnoseCmd.AddCommand(cniCmd)
	diagnoseCmd.AddCommand(connectionsCmd)
	diagnoseCmd.AddCommand(deploymentCmd)
	diagnoseCmd.AddCommand(versionCmd)
	diagnoseCmd.AddCommand(kpModeCmd)
	diagnoseCmd.AddCommand(allCmd)
}

func addVerboseFlag(command *cobra.Command) {
	command.Flags().BoolVar(&verboseOutput, "verbose", false, "produce verbose output")
}

func addNamespaceFlag(command *cobra.Command) {
	command.Flags().StringVar(&diagnose.KubeProxyPodNamespace, "namespace", "default",
		"namespace in which validation pods should be deployed")
}
