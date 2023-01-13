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

func TestAlibabaCloud_RegisterTargets(t *testing.T) {
	type args struct {
		req cloud.RegisterBackendServerRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "常规",
			args: args{
				req: cloud.RegisterBackendServerRequest{
					LoadBalancerId: "lb-2zepo6d2xxxxxxxxxxx",
					BackendServerList: []cloud.BackendServerItem{{
						ServerId:    "xxxxx",
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
			name: "常规",
			args: args{
				cloud.CreateListenerRequest{
					LoadBalancerId:    "lb-2ze2w4rosno895x6sf3k3",
					Protocol:          cloud.ProtocolTypeTCP,
					PortList:          []int{80},
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
