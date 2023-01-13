package ecloud

import (
	"errors"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *ECloud) CreateLoadBalancer(req cloud.CreateLoadBalancerRequest) (cloud.CreateLoadBalancerResponse, error) {
	// TODO implement me
	return cloud.CreateLoadBalancerResponse{}, errors.New("implement me")
}

func (p *ECloud) CreateListener(req cloud.CreateListenerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) DeregisterBackendServer(req cloud.DeregisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}
