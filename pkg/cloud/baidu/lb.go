package baidu

import (
	"errors"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *BaiduCloud) CreateLoadBalancer(req cloud.CreateLoadBalancerRequest) (cloud.CreateLoadBalancerResponse, error) {
	// you can call api like this
	// p.blbClient.CreateLoadBalancer()
	// TODO implement me
	return cloud.CreateLoadBalancerResponse{}, errors.New("implement me")
}

func (p *BaiduCloud) CreateListener(req cloud.CreateListenerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *BaiduCloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *BaiduCloud) DeregisterBackendServer(req cloud.DeregisterBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *BaiduCloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *BaiduCloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	return errors.New("do not use this api")
}
