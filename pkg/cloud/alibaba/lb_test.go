package alibaba

import (
	"testing"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func TestAlibabaCloud_CreateLoadBalancer(t *testing.T) {
	type args struct {
		req cloud.CreateLoadBalancerRequest
	}

	tests := []struct {
		name    string
		args    args
		want    cloud.CreateLoadBalancerResponse
		wantErr bool
	}{
		{
			name:    "常规",
			args:    args{req: cloud.CreateLoadBalancerRequest{LoadBalancerName: "api-call"}},
			wantErr: false,
		},
	}

	p, _ := New(_AK, _SK, "cn-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.CreateLoadBalancer(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLoadBalancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("CreateLoadBalancer() got = %v", got)
		})
	}
}

func TestAlibabaCloud_RegisterBackendServer(t *testing.T) {
	type args struct {
		req cloud.RegisterBackendServerRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "添加后端服务",
			args: args{
				req: cloud.RegisterBackendServerRequest{
					LoadBalancerId: "lb-2ze087rpg1xt1rhhyh15v",
					BackendServerList: []cloud.BackendServerItem{{
						ServerId:    "i-2ze4ati3ddzq54gs2blv",
						Port:        80,
						Weight:      10,
						Description: "测试",
					}},
				}},
		},
	}
	p, _ := New(_AK, _SK, "cn-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.RegisterBackendServer(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterBackendServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAlibabaCloud_CreateListener(t *testing.T) {
	type args struct {
		req cloud.CreateListenerRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TCP",
			args: args{
				cloud.CreateListenerRequest{
					LoadBalancerId:    "lb-2ze087rpg1xt1rhhyh15v",
					Protocol:          cloud.ProtocolTypeTCP,
					PortList:          []int{8081, 8082},
					BackendServerPort: 80,
				},
			},
		},
		{
			name: "UDP",
			args: args{
				cloud.CreateListenerRequest{
					LoadBalancerId:    "lb-2ze087rpg1xt1rhhyh15v",
					Protocol:          cloud.ProtocolTypeUDP,
					PortList:          []int{8083, 8084},
					BackendServerPort: 80,
				},
			},
		},
	}
	p, _ := New(_AK, _SK, "cn-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.CreateListener(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateListener() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlibabaCloud_DeregisterBackendServer(t *testing.T) {
	type args struct {
		req cloud.DeregisterBackendServerRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "常规",
			args: args{req: cloud.DeregisterBackendServerRequest{
				LoadBalancerId: "lb-2ze087rpg1xt1rhhyh15v",
				BackendServerList: []cloud.BackendServerItem{
					{
						ServerId:    "i-2ze4ati3ddzq54gs2blv",
						Port:        0,
						Weight:      0,
						Description: "",
					},
				},
			}},
		},
	}
	p, _ := New(_AK, _SK, "cn-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.DeregisterBackendServer(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("DeregisterBackendServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlibabaCloud_UpdateBackendServer(t *testing.T) {
	type args struct {
		req cloud.UpdateBackendServerRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "修改权重",
			args: args{req: cloud.UpdateBackendServerRequest{
				LoadBalancerId: "lb-2ze087rpg1xt1rhhyh15v",
				BackendServerList: []cloud.BackendServerItem{
					{
						ServerId:    "i-2ze4ati3ddzq54gs2blv",
						Port:        80,
						Weight:      99,
						Description: "",
					},
				},
			}},
		},
	}
	p, _ := New(_AK, _SK, "cn-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.UpdateBackendServer(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBackendServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlibabaCloud_StartLoadBalancerListener(t *testing.T) {
	type args struct {
		req cloud.StartLoadBalancerListenerRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "常规",
			args: args{req: cloud.StartLoadBalancerListenerRequest{
				LoadBalancerId: "lb-2ze087rpg1xt1rhhyh15v",
				ListenerPort:   8084,
				Protocol:       cloud.ProtocolTypeUDP,
			}},
		},
	}
	p, _ := New(_AK, _SK, "cn-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.StartLoadBalancerListener(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("StartLoadBalancerListener() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
