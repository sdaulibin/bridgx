package tencent

import (
	"errors"
	"fmt"
	"github.com/galaxy-future/BridgX/internal/logs"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/galaxy-future/BridgX/pkg/cloud"

	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func (p *TencentCloud) CreateLoadBalancer(req cloud.CreateLoadBalancerRequest) (cloud.CreateLoadBalancerResponse, error) {
	if req.LoadBalancerName == "" {
		return cloud.CreateLoadBalancerResponse{}, fmt.Errorf("loadBalancer name is empty")
	}
	var (
		loadBalancerType = "INTERNAL"
		subnetId         = "subnet-8up9sesx"
	)
	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerName = &req.LoadBalancerName
	request.LoadBalancerType = &loadBalancerType
	request.SubnetId = &subnetId
	response, err := p.clbClient.CreateLoadBalancer(request)
	if err != nil {
		return cloud.CreateLoadBalancerResponse{}, err
	}
	if response.Response.LoadBalancerIds == nil || len(response.Response.LoadBalancerIds) == 0 {
		return cloud.CreateLoadBalancerResponse{}, fmt.Errorf("loadBalancer ids length 0 or nil")
	}
	return cloud.CreateLoadBalancerResponse{
		LoadBalancerId: *response.Response.LoadBalancerIds[0],
	}, nil
}

func (p *TencentCloud) CreateListener(req cloud.CreateListenerRequest) error {
	if len(req.PortList) == 0 {
		return errors.New("listener ports empty")
	}
	request := clb.NewCreateListenerRequest()
	request.LoadBalancerId = &req.LoadBalancerId
	request.Protocol = (*string)(&req.Protocol)
	ports := make([]*int64, 0)
	for _, v := range req.PortList {
		ports = append(ports, common.Int64Ptr(int64(v)))
	}
	request.Ports = ports
	response, err := p.clbClient.CreateListener(request)
	if err != nil {
		return err
	}
	if response.Response.ListenerIds == nil || len(response.Response.ListenerIds) == 0 {
		return fmt.Errorf("Listener ids length 0 or nil")
	}
	return nil
}

func (p *TencentCloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	if len(req.BackendServerList) == 0 {
		return errors.New("register backend server list empty")
	}
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = &req.LoadBalancerId
	targets := make([]*clb.Target, 0)
	for _, server := range req.BackendServerList {
		target := &clb.Target{}
		target.Port = common.Int64Ptr(int64(server.Port))
		target.Weight = common.Int64Ptr(int64(server.Weight))
		target.InstanceId = common.StringPtr(server.ServerId)
		targets = append(targets, target)
	}
	request.Targets = targets
	_, err := p.clbClient.RegisterTargets(request)
	if err != nil {
		logs.Logger.Errorf(err.Error())
		return err
	}
	return nil
}

func (p *TencentCloud) DeregisterBackendServer(req cloud.DeregisterBackendServerRequest) error {
	if len(req.BackendServerList) == 0 {
		return errors.New("deregister backend server list empty")
	}
	request := clb.NewDeregisterTargetsRequest()
	request.LoadBalancerId = &req.LoadBalancerId
	targets := make([]*clb.Target, 0)
	for _, server := range req.BackendServerList {
		target := &clb.Target{}
		target.Port = common.Int64Ptr(int64(server.Port))
		target.Weight = common.Int64Ptr(int64(server.Weight))
		target.InstanceId = common.StringPtr(server.ServerId)
		targets = append(targets, target)
	}
	request.Targets = targets
	_, err := p.clbClient.DeregisterTargets(request)
	if err != nil {
		logs.Logger.Errorf(err.Error())
		return err
	}
	return nil
}

func (p *TencentCloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *TencentCloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	return errors.New("do not use this api")
}
