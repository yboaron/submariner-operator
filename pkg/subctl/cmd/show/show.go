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

package show

import (
	"github.com/spf13/cobra"
	"github.com/submariner-io/submariner-operator/internal/restconfig"
	"github.com/submariner-io/submariner-operator/pkg/subctl/cmd"
)

var (
	// showCmd represents the show command.
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show information about submariner",
		Long:  `This command shows information about some aspect of the submariner deployment in a cluster.`,
	}
	restConfigProducer = restconfig.NewProducer()
)

func init() {
	restConfigProducer.AddKubeConfigFlag(showCmd)
	cmd.AddToRootCommand(showCmd)
}
