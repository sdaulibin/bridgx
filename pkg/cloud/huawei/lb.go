package huawei

import (
	"errors"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *HuaweiCloud) CreateLoadBalancer(req cloud.CreateLoadBalancerRequest) (cloud.CreateLoadBalancerResponse, error) {
	// you can call api like this
	// p.elbClient.CreateLoadbalancer()

	// TODO implement me
	return cloud.CreateLoadBalancerResponse{}, errors.New("implement me")
}

func (p *HuaweiCloud) CreateListener(req cloud.CreateListenerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *HuaweiCloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *HuaweiCloud) DeregisterBackendServer(req cloud.DeregisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *HuaweiCloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *HuaweiCloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	return errors.New("do not use this api")
}
