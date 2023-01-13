package tencent

import (
	"errors"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *TencentCloud) CreateLoadBalancer(req cloud.CreateLoadBalancerRequest) (cloud.CreateLoadBalancerResponse, error) {
	// you can call api like this
	// p.clbClient.CreateLoadBalancer()
	// TODO implement me
	return cloud.CreateLoadBalancerResponse{}, errors.New("implement me")
}

func (p *TencentCloud) CreateListener(req cloud.CreateListenerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *TencentCloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *TencentCloud) DeregisterBackendServer(req cloud.DeregisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *TencentCloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *TencentCloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	return errors.New("do not use this api")
}
