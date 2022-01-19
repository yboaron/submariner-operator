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
	"github.com/submariner-io/submariner-operator/internal/execute"
	"strings"

	"github.com/submariner-io/submariner-operator/internal/cli"
	"github.com/submariner-io/submariner-operator/pkg/subctl/resource"
)

var KubeProxyPodNamespace string

const (
	kubeProxyIPVSIfaceCommand = "ip a s kube-ipvs0"
	missingInterface          = "ip: can't find device"
)

func CheckKubeProxyMode(cluster *execute.Cluster) bool {
	status := cli.NewReporter()
	status.Start("Checking Submariner support for the kube-proxy mode")

	scheduling := resource.PodScheduling{ScheduleOn: resource.GatewayNode, Networking: resource.HostNetworking}

	podOutput, err := resource.SchedulePodAwaitCompletion(&resource.PodConfig{
		Name:       "query-iface-list",
		ClientSet:  cluster.KubeClient,
		Scheduling: scheduling,
		Namespace:  KubeProxyPodNamespace,
		Command:    kubeProxyIPVSIfaceCommand,
	})
	if err != nil {
		status.Failure("Error spawning the network pod: %v", err)
		status.End()
		return false
	}

	failed := false
	if strings.Contains(podOutput, missingInterface) {
		status.Success("The kube-proxy mode is supported")
	} else {
		status.Failure("The cluster is deployed with kube-proxy ipvs mode which Submariner does not support")
		failed = true
	}

	status.End()

	return !failed
}
