package tencent

import (
	"testing"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

var (
	ak = ""
	sk = ""
)

func TestTencentCloud_CreateLoadBalancer(t *testing.T) {
	type args struct {
		Req cloud.CreateLoadBalancerRequest
	}

	tests := []struct {
		Name    string
		Args    args
		Want    cloud.CreateLoadBalancerResponse
		WantErr bool
	}{
		{
			Name: "测试新建lb",
			Args: args{
				Req: cloud.CreateLoadBalancerRequest{
					LoadBalancerName: "test20230221001",
				},
			},
			WantErr: false,
		},
	}
	p, _ := New(ak, sk, "ap-beijing")
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := p.CreateLoadBalancer(tt.Args.Req)
			if (err != nil) != tt.WantErr {
				t.Errorf("TestTencentCloud CreateLoadBalancer() error = %v, wantErr %v", err, tt.WantErr)
				return
			}
			t.Logf("TestTencentCloud CreateLoadBalancer() got = %v", got)
		})
	}
}

func TestTencentCloud_CreateListener(t *testing.T) {
	type args struct {
		Req cloud.CreateListenerRequest
	}

	tests := []struct {
		Name    string
		Args    args
		WantErr bool
	}{
		{
			Name: "新建listener",
			Args: args{Req: cloud.CreateListenerRequest{
				LoadBalancerId: "lb-r70cytxd",
				Protocol:       cloud.ProtocolTypeHTTP,
				PortList:       []int{8088, 8089},
			}},
		},
	}
	p, _ := New(ak, sk, "ap-beijing")
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := p.CreateListener(tt.Args.Req)
			if (err != nil) != tt.WantErr {
				t.Errorf("TestTencentCloud CreateListener() error = %v, wantErr %v", err, tt.WantErr)
				return
			}
		})
	}
}

func TestTencentCloud_RegisterBackendServer(t *testing.T) {
	type args struct {
		Req cloud.RegisterBackendServerRequest
	}
	tests := []struct {
		Name    string
		Args    args
		WantErr bool
	}{
		{
			Name: "绑定后端服务",
			Args: args{
				Req: cloud.RegisterBackendServerRequest{
					LoadBalancerId: "",
					ListenerId:     "",
					BackendServerList: []cloud.BackendServerItem{
						{
							ServerId: "",
							Port:     8801,
							Weight:   50,
						},
					},
				},
			},
			WantErr: false,
		},
	}
	p, _ := New(ak, sk, "ap-beijing")
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := p.RegisterBackendServer(tt.Args.Req)
			if (err != nil) != tt.WantErr {
				t.Errorf("TestTencentCloud RegisterBackendServer() error = %v, wantErr %v", err, tt.WantErr)
			}
		})
	}
}

func TestTencentCloud_DeregisterBackendServer(t *testing.T) {
	type args struct {
		Req cloud.DeregisterBackendServerRequest
	}
	tests := []struct {
		Name    string
		Args    args
		WantErr bool
	}{
		{
			Name: "解绑后端服务",
			Args: args{
				Req: cloud.DeregisterBackendServerRequest{
					LoadBalancerId: "",
					ListenerId:     "",
					BackendServerList: []cloud.BackendServerItem{
						{
							ServerId: "",
						},
					},
				},
			},
			WantErr: false,
		},
	}
	p, _ := New(ak, sk, "ap-beijing")
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := p.DeregisterBackendServer(tt.Args.Req)
			if (err != nil) != tt.WantErr {
				t.Errorf("TestTencentCloud DeregisterBackendServer() error = %v, wantErr %v", err, tt.WantErr)
			}
		})
	}
}

func TestTencentCloud_UpdateBackendServer(t *testing.T) {
	type args struct {
		Req cloud.UpdateBackendServerRequest
	}
	tests := []struct {
		Name    string
		Args    args
		WantErr bool
	}{
		{
			Name: "修改权重",
			Args: args{
				Req: cloud.UpdateBackendServerRequest{
					LoadBalancerId: "",
					ListenerId:     "",
					BackendServerList: []cloud.BackendServerItem{
						{
							ServerId: "",
							Weight:   100,
						},
					},
				},
			},
			WantErr: false,
		},
	}
	p, _ := New(ak, sk, "ap-beijing")
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := p.UpdateBackendServer(tt.Args.Req)
			if (err != nil) != tt.WantErr {
				t.Errorf("TestTencentCloud UpdateBackendServer() error = %v, wantErr %v", err, tt.WantErr)
			}
		})
	}
}
