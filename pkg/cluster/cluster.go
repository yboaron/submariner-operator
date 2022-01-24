package cluster

import (
	"context"
	"github.com/pkg/errors"
	"github.com/submariner-io/submariner-operator/api/submariner/v1alpha1"
	"github.com/submariner-io/submariner-operator/internal/constants"
	"github.com/submariner-io/submariner-operator/pkg/client"
	submarinerv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Info struct {
	Name            string
	ClientProducer  client.Producer
	submarinerOrErr *interface{}
	Submariner      *v1alpha1.Submariner
}

func New(clusterName string, clientProducer client.Producer) (*Info, error) {
	cluster := &Info{
		Name:   clusterName,
		ClientProducer: clientProducer,
	}

	var err error
	cluster.Submariner, err = cluster.GetSubmariner()
	if err != nil {
		return nil, errors.Wrap(err,"Error retrieving Submariner")
	}

	return cluster, nil
}

func (c *Info) GetSubmariner() (*v1alpha1.Submariner, error) {
	if c.submarinerOrErr == nil {
		c.submarinerOrErr = new(interface{})

		submariner, err := c.ClientProducer.ForOperator().SubmarinerV1alpha1().Submariners(constants.SubmarinerNamespace).
			Get(context.TODO(), constants.SubmarinerName, metav1.GetOptions{})
		if !apierrors.IsNotFound(err) {
			if err != nil {
				*c.submarinerOrErr = errors.Wrap(err, "Error retrieving Submariner resource")
			} else {
				*c.submarinerOrErr = submariner
			}
		}
	}
	if *c.submarinerOrErr == nil {
		return nil, nil
	}
	if err, ok := (*c.submarinerOrErr).(error); ok {
		return nil, err
	}
	return (*c.submarinerOrErr).(*v1alpha1.Submariner), nil
}

func (c *Info) GetGateways() ([]submarinerv1.Gateway, error) {
	gateways, err := c.ClientProducer.ForSubmariner().SubmarinerV1().Gateways(constants.OperatorNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error listing Gateways")
	}

	return gateways.Items, nil
}
