package alibaba

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *AlibabaCloud) CreateLoadBalancer(req cloud.CreateLoadBalancerRequest) (cloud.CreateLoadBalancerResponse, error) {
	var (
		addressTypeInternet = "internet"
	)
	request := slb.CreateCreateLoadBalancerRequest()
	request.LoadBalancerName = req.LoadBalancerName
	request.AddressType = addressTypeInternet
	request.LoadBalancerSpec = _loadBalancerSpecS1Small

	response, err := p.slbClient.CreateLoadBalancer(request)
	if err != nil {
		return cloud.CreateLoadBalancerResponse{}, err
	}
	if response.GetHttpStatus() != http.StatusOK {
		return cloud.CreateLoadBalancerResponse{}, fmt.Errorf("http status %v", response.GetHttpStatus())
	}
	return cloud.CreateLoadBalancerResponse{
		LoadBalancerId: response.LoadBalancerId,
	}, nil
}

func (p *AlibabaCloud) RegisterBackendServer(req cloud.RegisterBackendServerRequest) error {
	if len(req.BackendServerList) == 0 {
		return errors.New("target backend server empty")
	}
	var (
		backendServers []BackendServer
		serverTypeEcs  = "ecs"
	)
	request := slb.CreateAddBackendServersRequest()
	request.LoadBalancerId = req.LoadBalancerId
	for _, server := range req.BackendServerList {
		backendServers = append(backendServers, BackendServer{
			ServerId:    server.ServerId,
			Port:        server.Port,
			Weight:      server.Weight,
			Type:        serverTypeEcs,
			Description: server.Description,
		})
	}
	backendServersJson, err := json.Marshal(backendServers)
	if err != nil {
		return err
	}
	request.BackendServers = string(backendServersJson)
	response, err := p.slbClient.AddBackendServers(request)
	if err != nil {
		return err
	}
	if response.GetHttpStatus() != http.StatusOK {
		return fmt.Errorf("http status %v", response.GetHttpStatus())
	}

	return nil
}

func (p *AlibabaCloud) DeregisterBackendServer(req cloud.DeregisterBackendServerRequest) error {
	if len(req.BackendServerList) == 0 {
		return errors.New("target backend server empty")
	}
	var (
		backendServers []BackendServer
		serverTypeEcs  = "ecs"
	)
	request := slb.CreateRemoveBackendServersRequest()
	request.LoadBalancerId = req.LoadBalancerId
	for _, server := range req.BackendServerList {
		backendServers = append(backendServers, BackendServer{
			ServerId:    server.ServerId,
			Port:        server.Port,
			Weight:      server.Weight,
			Type:        serverTypeEcs,
			Description: server.Description,
		})
	}
	backendServersJson, err := json.Marshal(backendServers)
	if err != nil {
		return err
	}
	request.BackendServers = string(backendServersJson)
	response, err := p.slbClient.RemoveBackendServers(request)
	if err != nil {
		return err
	}
	if response.GetHttpStatus() != http.StatusOK {
		return fmt.Errorf("http status %v", response.GetHttpStatus())
	}
	return nil
}

func (p *AlibabaCloud) CreateListener(req cloud.CreateListenerRequest) error {
	if req.Protocol == "" {
		return errors.New("protocol empty")
	}
	if len(req.PortList) == 0 {
		return errors.New("port empty")
	}
	switch req.Protocol {
	case cloud.ProtocolTypeTCP:
		return p.createTCPListener(req)
	case cloud.ProtocolTypeUDP:
		return p.createUDPListener(req)
	default:
		return fmt.Errorf("protocol %s can not handle", req.Protocol)
	}
}

func (p *AlibabaCloud) createTCPListener(req cloud.CreateListenerRequest) error {
	// 对于按流量计费的公网负载均衡实例，可以将带宽峰值设置为-1，即不限制带宽峰值。
	defaultBandwidth := -1
	request := slb.CreateCreateLoadBalancerTCPListenerRequest()
	request.LoadBalancerId = req.LoadBalancerId
	request.Bandwidth = requests.NewInteger(defaultBandwidth)
	request.BackendServerPort = requests.NewInteger(req.BackendServerPort)
	request.VServerGroupId = req.ServerGroupId

	for _, port := range req.PortList {
		request.ListenerPort = requests.NewInteger(port)
		response, err := p.slbClient.CreateLoadBalancerTCPListener(request)
		if err != nil {
			return err
		}
		if response.GetHttpStatus() != http.StatusOK {
			return fmt.Errorf("http status %v", response.GetHttpStatus())
		}
	}
	return nil
}

func (p *AlibabaCloud) createUDPListener(req cloud.CreateListenerRequest) error {
	// 对于按流量计费的公网负载均衡实例，可以将带宽峰值设置为-1，即不限制带宽峰值。
	defaultBandwidth := -1
	request := slb.CreateCreateLoadBalancerUDPListenerRequest()
	request.LoadBalancerId = req.LoadBalancerId
	request.Bandwidth = requests.NewInteger(defaultBandwidth)
	request.BackendServerPort = requests.NewInteger(req.BackendServerPort)
	request.VServerGroupId = req.ServerGroupId
	for _, port := range req.PortList {
		request.ListenerPort = requests.NewInteger(port)
		response, err := p.slbClient.CreateLoadBalancerUDPListener(request)
		if err != nil {
			return err
		}
		if response.GetHttpStatus() != http.StatusOK {
			return fmt.Errorf("http status %v", response.GetHttpStatus())
		}
	}
	return nil
}

func (p *AlibabaCloud) UpdateBackendServer(req cloud.UpdateBackendServerRequest) error {
	if len(req.BackendServerList) == 0 {
		return errors.New("target backend server empty")
	}
	var (
		backendServers []BackendServer
		serverTypeEcs  = "ecs"
	)
	request := slb.CreateSetBackendServersRequest()
	request.LoadBalancerId = req.LoadBalancerId
	for _, server := range req.BackendServerList {
		if server.ServerId == "" {
			return errors.New("ServerId empty")
		}
		i := BackendServer{
			ServerId: server.ServerId,
			Type:     serverTypeEcs,
		}
		if server.Port != 0 {
			i.Port = server.Port
		}
		if server.Weight != 0 {
			i.Weight = server.Weight
		}
		if len(server.Description) > 0 {
			i.Description = server.Description
		}
		backendServers = append(backendServers, i)
	}
	backendServersJson, err := json.Marshal(backendServers)
	if err != nil {
		return err
	}
	request.BackendServers = string(backendServersJson)
	response, err := p.slbClient.SetBackendServers(request)
	if err != nil {
		return err
	}
	if response.GetHttpStatus() != http.StatusOK {
		return fmt.Errorf("http status %v", response.GetHttpStatus())
	}

	return nil
}

func (p *AlibabaCloud) StartLoadBalancerListener(req cloud.StartLoadBalancerListenerRequest) error {
	request := slb.CreateStartLoadBalancerListenerRequest()
	request.LoadBalancerId = req.LoadBalancerId
	request.ListenerPort = requests.NewInteger(req.ListenerPort)
	response, err := p.slbClient.StartLoadBalancerListener(request)
	if err != nil {
		return err
	}
	if response.GetHttpStatus() != http.StatusOK {
		return fmt.Errorf("http status %v", response.GetHttpStatus())
	}
	return nil
}
