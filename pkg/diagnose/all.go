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

package diagnose

import (
	"context"
	"fmt"
	"github.com/submariner-io/submariner-operator/internal/constants"
	"github.com/submariner-io/submariner-operator/internal/execute"

	"github.com/pkg/errors"
	"github.com/submariner-io/submariner-operator/internal/cli"
	"github.com/submariner-io/submariner-operator/pkg/subctl/cmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func All(cluster *execute.Cluster) bool {
	success := CheckK8sVersion(cluster)

	fmt.Println()

	status := cli.NewReporter()
	if cluster.Submariner == nil {
		status.Warning(constants.SubmMissingMessage)

		return success
	}

	success = CheckCNIConfig(cluster) && success

	fmt.Println()

	success = CheckConnections(cluster) && success

	fmt.Println()

	success = CheckDeployments(cluster) && success

	fmt.Println()

	success = CheckKubeProxyMode(cluster) && success

	fmt.Println()

	//	Skip for now until firewall commands are migrated
	/*
	success = checkFirewallMetricsConfig(cluster) && success

	fmt.Println()

	success = checkVxLANConfig(cluster) && success

	fmt.Println()
   */
	fmt.Printf("Skipping inter-cluster firewall check as it requires two kubeconfigs." +
		" Please run \"subctl diagnose firewall inter-cluster\" command manually.\n")

	return success
}

func getNumNodesOfCluster(cluster *cmd.Cluster) (int, error) {
	nodes, err := cluster.KubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "error listing Nodes")
	}

	return len(nodes.Items), nil
}

func isClusterSingleNode(cluster *cmd.Cluster, status *cli.Status) bool {
	numNodesOfCluster, err := getNumNodesOfCluster(cluster)
	if err != nil {
		status.EndWithFailure("Error listing the number of nodes of the cluster: %v", err)
		return true
	}

	if numNodesOfCluster == 1 {
		status.EndWithSuccess("Skipping this check as it's a single node cluster.")
		return true
	}

	return false
}
