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
		loadBalancerType = "OPEN"
	)
	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerName = &req.LoadBalancerName
	request.LoadBalancerType = &loadBalancerType
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
		return errors.New("Listener ids length 0 or nil")
	}
	return nil
}

func (p *TencentCloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	if len(req.BackendServerList) == 0 {
		return errors.New("register backend server list empty")
	}
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = &req.LoadBalancerId
	request.ListenerId = &req.ListenerId
	request.Targets = createTargets(req.BackendServerList)
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
	request.ListenerId = &req.ListenerId
	request.Targets = createTargets(req.BackendServerList)
	_, err := p.clbClient.DeregisterTargets(request)
	if err != nil {
		logs.Logger.Errorf(err.Error())
		return err
	}
	return nil
}

func (p *TencentCloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	if len(req.BackendServerList) == 0 {
		return errors.New("update backend server list empty")
	}
	request := clb.NewModifyTargetWeightRequest()
	request.LoadBalancerId = &req.LoadBalancerId
	request.ListenerId = &req.ListenerId
	request.Targets = createTargets(req.BackendServerList)
	_, err := p.clbClient.ModifyTargetWeight(request)
	if err != nil {
		logs.Logger.Errorf(err.Error())
		return err
	}
	return nil
}

func (p *TencentCloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	return errors.New("do not use this api")
}

func createTargets(serverList []cloud.BackendServerItem) []*clb.Target {
	targets := make([]*clb.Target, 0)
	for _, server := range serverList {
		target := &clb.Target{}
		target.Port = common.Int64Ptr(int64(server.Port))
		target.Weight = common.Int64Ptr(int64(server.Weight))
		target.InstanceId = common.StringPtr(server.ServerId)
		targets = append(targets, target)
	}
	return targets
}

func (p *TencentCloud) CreateListenerRules(req cloud.CreateListenerRuleRequest) error {
	if len(req.ListenerRuleList) == 0 {
		return errors.New("create listener rules list empty")
	}
	request := clb.NewCreateRuleRequest()
	request.LoadBalancerId = &req.LoadBalancerId
	request.ListenerId = &req.ListenerId
	rules := make([]*clb.RuleInput, 0)
	for _, v := range req.ListenerRuleList {
		rule := &clb.RuleInput{}
		rule.Domain = common.StringPtr(v.Domain)
		rule.Url = common.StringPtr(v.Url)
		rules = append(rules, rule)
	}
	request.Rules = rules
	_, err := p.clbClient.CreateRule(request)
	if err != nil {
		logs.Logger.Errorf(err.Error())
		return err
	}
	return nil
}
