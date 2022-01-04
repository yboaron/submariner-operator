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

package cleanup

import (
	"github.com/spf13/cobra"
	"github.com/submariner-io/submariner-operator/internal/restconfig"
	"github.com/submariner-io/submariner-operator/pkg/cloud/cleanup"
	"github.com/submariner-io/submariner-operator/pkg/subctl/cmd/cloud/aws"
	"github.com/submariner-io/submariner-operator/pkg/subctl/cmd/cloud/gcp"
)

var parentRestConfigProducer *restconfig.Producer

// NewCommand returns a new cobra.Command used to prepare a cloud infrastructure.
func NewCommand(restConfigProducer *restconfig.Producer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Clean up the cloud",
		Long:  `This command cleans up the cloud after Submariner uninstallation.`,
	}
	parentRestConfigProducer = restConfigProducer

	cmd.AddCommand(newAWSCleanupCommand())
	cmd.AddCommand(newGCPCleanupCommand())
	cmd.AddCommand(newGenericCleanupCommand())

	return cmd
}

// NewCommand returns a new cobra.Command used to prepare a cloud infrastructure.
func newAWSCleanupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "Clean up an AWS cloud",
		Long: "This command cleans up an OpenShift installer-provisioned infrastructure (IPI) on AWS-based" +
			" cloud after Submariner uninstallation.",
		Run: cleanup.Aws,
	}

	aws.AddAWSFlags(cmd)

	return cmd
}

// newGCPCleanupCommand returns a new cobra.Command used to prepare a cloud infrastructure.
func newGCPCleanupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcp",
		Short: "Clean up a GCP cloud",
		Long:  "This command cleans up an installer-provisioned infrastructure (IPI) on GCP-based cloud after Submariner uninstallation.",
		Run:   cleanup.GCP,
	}

	gcp.AddGCPFlags(cmd)

	return cmd
}

func newGenericCleanupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generic",
		Short: "Cleans up a cluster after Submariner uninstallation",
		Long:  "This command removes the labels from gateway nodes after Submariner uninstallation.",
		Run:   cleanup.GenericCluster,
	}

	return cmd
}