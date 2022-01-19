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

package execute

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/submariner-io/submariner-operator/api/submariner/v1alpha1"
	"github.com/submariner-io/submariner-operator/internal/constants"
	"github.com/submariner-io/submariner-operator/internal/restconfig"
	"github.com/submariner-io/submariner-operator/pkg/client"
	operatorClientset "github.com/submariner-io/submariner-operator/pkg/client/clientset/versioned"
	submarinerv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
	subClientsetv1 "github.com/submariner-io/submariner/pkg/client/clientset/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Cluster struct {
	Config     *rest.Config
	Name       string
	KubeClient kubernetes.Interface
	DynClient  dynamic.Interface
	SubmClient subClientsetv1.Interface
	Submariner *v1alpha1.Submariner
}

func OnMultiCluster(restConfigProducer restconfig.Producer, run func(*Cluster) bool) {
	success := true

	for _, config := range restConfigProducer.MustGetForClusters() {
		fmt.Printf("Cluster %q\n", config.ClusterName)

		cluster, errMsg := NewCluster(config.Config, config.ClusterName)
		if cluster == nil {
			success = false

			fmt.Println(errMsg)
			fmt.Println()

			continue
		}

		success = run(cluster) && success

		fmt.Println()
	}

	if !success {
		os.Exit(1)
	}
}

func NewCluster(config *rest.Config, clusterName string) (*Cluster, string) {
	cluster := &Cluster{
		Config: config,
		Name:   clusterName,
	}

	clientProducer, err := client.NewProducerFromRestConfig(config)
	if err != nil {
		return nil, fmt.Sprintf("Error creating client producer: %v", err)
	}

	cluster.KubeClient = clientProducer.ForKubernetes()
	cluster.DynClient = clientProducer.ForDynamic()
	cluster.SubmClient = clientProducer.ForSubmariner()

	cluster.Submariner, err = getSubmarinerResourceWithError(clientProducer.ForOperator())
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, fmt.Sprintf("Error retrieving Submariner resource: %v", err)
	}

	return cluster, ""
}

func getSubmarinerResourceWithError(operatorClient operatorClientset.Interface) (*v1alpha1.Submariner, error) {
	submariner, err := operatorClient.SubmarinerV1alpha1().Submariners(constants.SubmarinerNamespace).
		Get(context.TODO(), constants.SubmarinerName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithMessagef(err, "error retrieving Submariner object %s", constants.SubmarinerName)
	}

	return submariner, nil
}

func (c *Cluster) GetGateways() ([]submarinerv1.Gateway, error) {
	gateways, err := c.SubmClient.SubmarinerV1().Gateways(constants.OperatorNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error listing Gateways")
	}

	return gateways.Items, nil
}
