package tencent

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"testing"
)

var (
	ak = ""
	sk = ""
)

func TestTencentCloud_CreateLoadBalancer(t *testing.T) {
	type args struct {
		request cloud.CreateLoadBalancerRequest
	}

	tests := []struct {
		name    string
		args    args
		want    cloud.CreateLoadBalancerResponse
		wantErr bool
	}{
		{
			name: "测试新建lb",
			args: args{
				request: cloud.CreateLoadBalancerRequest{
					LoadBalancerName: "test001",
				},
			},
			wantErr: false,
		},
	}
	p, _ := New(ak, sk, "ap-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.CreateLoadBalancer(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestTencentCloud CreateLoadBalancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("TestTencentCloud CreateLoadBalancer() got = %v", got)
		})
	}
}

func TestTencentCloud_CreateListener(t *testing.T) {
	type args struct {
		req cloud.CreateListenerRequest
	}

	tests := []struct {
		name    string
		agrs    args
		wantErr bool
	}{
		{
			name: "新建listener",
			agrs: args{req: cloud.CreateListenerRequest{
				LoadBalancerId: "lb-hccgj54f",
				Protocol:       cloud.ProtocolTypeHTTP,
				PortList:       []int{8088, 8089},
			}},
		},
	}
	p, _ := New(ak, sk, "ap-beijing")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.CreateListener(tt.agrs.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestTencentCloud CreateListener() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
