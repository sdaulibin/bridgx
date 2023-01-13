package aws

import (
	"errors"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *AWSCloud) CreateLoadBalancer(req cloud.CreateLoadBalancerRequest) (cloud.CreateLoadBalancerResponse, error) {
	// you can call api like this
	// p.elbClient.CreateLoadBalancer()
	// TODO implement me
	return cloud.CreateLoadBalancerResponse{}, errors.New("implement me")
}

func (p *AWSCloud) CreateListener(req cloud.CreateListenerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *AWSCloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *AWSCloud) DeregisterBackendServer(req cloud.DeregisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *AWSCloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	return errors.New("do not use this api")
}

func (p *AWSCloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	return errors.New("do not use this api")
}
