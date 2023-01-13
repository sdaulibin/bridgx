package baidu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/baidubce/bce-sdk-go/model"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/baidubce/bce-sdk-go/services/blb"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
)

type BaiduCloud struct {
	ak        string
	sk        string
	vpcClient *vpc.Client
	bccClient *bcc.Client
	eipClient *eip.Client
	bosClient *bos.Client
	blbClient *blb.Client
}

func New(AK, SK, regionId string) (*BaiduCloud, error) {
	ep, ok := _endPoints[strings.ToLower(regionId)]
	if !ok {
		return nil, errors.New("regionId error:" + regionId)
	}

	vpcClient, err := vpc.NewClient(AK, SK, fmt.Sprintf("bcc%s", ep))
	if err != nil {
		return nil, err
	}

	bccClient, err := bcc.NewClient(AK, SK, fmt.Sprintf("bcc%s", ep))
	if err != nil {
		return nil, err
	}
	eipClient, err := eip.NewClient(AK, SK, fmt.Sprintf("eip%s", ep))
	if err != nil {
		return nil, err
	}
	bosClient, err := bos.NewClient(AK, SK, fmt.Sprintf("%s.bcebos.com", regionId))
	if err != nil {
		return nil, err
	}
	blbClient, err := blb.NewClient(AK, SK, fmt.Sprintf("blb%s", ep))
	if err != nil {
		return nil, err
	}
	return &BaiduCloud{
		ak:        AK,
		sk:        SK,
		vpcClient: vpcClient,
		bccClient: bccClient,
		eipClient: eipClient,
		bosClient: bosClient,
		blbClient: blbClient,
	}, nil
}

// BatchCreate DryRun 没有用到
func (p *BaiduCloud) BatchCreate(m cloud.Params, num int) (instanceIds []string, err error) {

	if m.DryRun == true {
		if len(strings.Split(m.Network.SecurityGroup, ",")) != 1 {
			return nil, fmt.Errorf("only one security group is supported")
		}
		return nil, nil
	}

	var createCdsList []api.CreateCdsModel
	for _, d := range m.Disks.DataDisk {
		createCdsList = append(createCdsList, api.CreateCdsModel{
			StorageType: api.StorageType(d.Category),
			CdsSizeInGB: d.Size,
		})
	}

	var tags []model.TagModel
	for _, item := range m.Tags {
		tags = append(tags, model.TagModel{
			TagKey:   item.Key,
			TagValue: item.Value,
		})
	}

	periodUnit := "Month"
	period := m.Charge.Period
	if m.Charge.ChargeType == cloud.PrePaid && m.Charge.PeriodUnit == cloud.Year {
		period *= 12
	}
	request := &api.CreateInstanceBySpecArgs{
		ImageId: m.ImageId,
		Billing: api.Billing{
			PaymentTiming: api.PaymentTimingType(_inEcsChargeType[m.Charge.ChargeType]), // https://cloud.baidu.com/doc/BCC/s/6jwvyo0q2#billing
			Reservation: &api.Reservation{
				ReservationLength:   period,
				ReservationTimeUnit: periodUnit,
			},
		},
		Spec:                  m.InstanceType, // https://cloud.baidu.com/doc/BCC/s/6jwvyo0q2#instancetype
		RootDiskSizeInGb:      m.Disks.SystemDisk.Size,
		RootDiskStorageType:   api.StorageType(m.Disks.SystemDisk.Category), // https://cloud.baidu.com/doc/BCC/s/6jwvyo0q2#storagetype
		CreateCdsList:         createCdsList,
		NetWorkCapacityInMbps: m.Network.InternetMaxBandwidthOut,
		PurchaseCount:         num,
		AdminPass:             m.Password,
		ZoneName:              m.Zone,
		SubnetId:              m.Network.SubnetId,
		SecurityGroupId:       strings.Split(m.Network.SecurityGroup, ",")[0],
		Tags:                  tags,
		InternetChargeType:    m.Network.InternetChargeType, // https://cloud.baidu.com/doc/BCC/s/6jwvyo0q2#internetchargetype
		InternalIps:           nil,
		DeployIdList:          nil,
		DetetionProtection:    0,
	}
	r, err := p.bccClient.CreateInstanceBySpec(request)
	if err != nil {
		return nil, err
	} else {
		return r.InstanceIds, nil
	}
}

func (p *BaiduCloud) ProviderType() string {
	return cloud.BaiduCloud
}

// GetInstances 缺失InternetChargeType，InternetIpType
func (p *BaiduCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	for _, item := range ids {
		r, err := p.bccClient.GetInstanceDetail(item)
		if err != nil {
			return nil, err
		}
		var SecurityGroup, _ = p.bccClient.ListSecurityGroup(&api.ListSecurityGroupArgs{
			InstanceId: item,
		})
		var SecurityGroups []string
		for _, v := range SecurityGroup.SecurityGroups {
			SecurityGroups = append(SecurityGroups, v.Id)
		}
		var expireAt *time.Time
		if r.Instance.ExpireTime != "" {
			expireTime, _ := time.Parse("2006-01-02T15:04:05Z", r.Instance.ExpireTime)
			expireAt = &expireTime
		}

		instances = append(instances, cloud.Instance{
			Id:       r.Instance.InstanceId,
			CostWay:  _ecsChargeType[r.Instance.PaymentTiming],
			Provider: cloud.BaiduCloud,
			IpInner:  r.Instance.InternalIP,
			IpOuter:  r.Instance.PublicIP,
			Network: &cloud.Network{
				VpcId:                   r.Instance.VpcId,
				SubnetId:                r.Instance.SubnetId,
				SecurityGroup:           strings.Join(SecurityGroups, ","),
				InternetChargeType:      "",
				InternetMaxBandwidthOut: r.Instance.NetworkCapacityInMbps,
				InternetIpType:          "",
			},
			ImageId:  r.Instance.ImageId,
			Status:   _ecsStatus[string(r.Instance.Status)],
			ExpireAt: expireAt,
		})
	}

	return instances, nil
}
func (p *BaiduCloud) GetInstancesByTags(region string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	var AllinstanceIds [][]string
	for _, v := range tags {
		var instanceOfTag []string
		result, err := p.bccClient.ListServersByMarkerV3(&api.ListServerRequestV3Args{
			Tag: model.TagModel{TagKey: v.Key, TagValue: v.Value},
		})
		if err != nil {
			return nil, err
		}
		for _, v := range result.Instances {
			instanceOfTag = append(instanceOfTag, v.InstanceId)
		}
		AllinstanceIds = append(AllinstanceIds, instanceOfTag)
	}
	interInstanceId := utils.Intersect(AllinstanceIds) // 多个instanceId切片求交集
	instances, err = p.ListbyId(interInstanceId)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (p *BaiduCloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	return p.GetInstancesByTags(regionId, []cloud.Tag{{
		Key:   cloud.ClusterName,
		Value: clusterName,
	}})
}

func (p *BaiduCloud) BatchDelete(ids []string, regionId string) error {
	for _, id := range ids {
		err := p.bccClient.DeleteInstance(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *BaiduCloud) StartInstances(ids []string) error {
	for _, id := range ids {
		err := p.bccClient.StartInstance(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *BaiduCloud) StopInstances(ids []string) error {
	for _, id := range ids {
		err := p.bccClient.StopInstance(id, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *BaiduCloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	request := &vpc.CreateVPCArgs{
		Name:        req.VpcName,
		Cidr:        req.CidrBlock,
		ClientToken: "",
		Description: "",
		Tags:        nil,
	}

	response, err := p.vpcClient.CreateVPC(request)
	if err != nil {
		return cloud.CreateVpcResponse{}, err
	}

	return cloud.CreateVpcResponse{
		VpcId:     response.VPCID,
		RequestId: "",
	}, nil
}

// GetVPC 缺少createAt， status也没有返回值，设置为默认
func (p *BaiduCloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	response, err := p.vpcClient.GetVPCDetail(req.VpcId)
	if err != nil {
		return cloud.GetVpcResponse{}, err
	}
	return cloud.GetVpcResponse{
		Vpc: cloud.VPC{
			VpcId:     response.VPC.VPCId,
			VpcName:   response.VPC.Name,
			CidrBlock: response.VPC.Cidr,
			RegionId:  req.RegionId,
			Status:    cloud.VPCStatusAvailable,
			CreateAt:  "",
		},
	}, nil
}

// CreateSwitch GatewayIp 没有用到
func (p *BaiduCloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {
	r, err := p.vpcClient.CreateSubnet(&vpc.CreateSubnetArgs{
		ClientToken:      "",
		Name:             req.VSwitchName,
		ZoneName:         req.ZoneId,
		Cidr:             req.CidrBlock,
		VpcId:            req.VpcId,
		VpcSecondaryCidr: "",
		SubnetType:       "BCC", // BCC BCC_NAT BBC三种
		Description:      "",
		Tags:             nil,
	})

	if err != nil {
		return cloud.CreateSwitchResponse{}, err
	} else {
		return cloud.CreateSwitchResponse{
			RequestId: "",
			SwitchId:  r.SubnetId,
		}, nil
	}
}

// GetSwitch 缺失GatewayIp  Vsstatus设为默认
func (p *BaiduCloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {
	r, err := p.vpcClient.GetSubnetDetail(req.SwitchId)
	if err != nil {
		return cloud.GetSwitchResponse{}, err
	} else {
		return cloud.GetSwitchResponse{
			Switch: cloud.Switch{
				VpcId:                   r.Subnet.VPCId,
				SwitchId:                r.Subnet.SubnetId,
				Name:                    r.Subnet.Name,
				IsDefault:               0,
				AvailableIpAddressCount: r.Subnet.AvailableIp,
				VStatus:                 cloud.SubnetAvailable,
				CreateAt:                r.Subnet.CreatedTime,
				ZoneId:                  r.Subnet.ZoneName,
				CidrBlock:               r.Subnet.Cidr,
				GatewayIp:               "",
			},
		}, nil
	}
}

// rules 不用，securitGroupType没有用到
func (p *BaiduCloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {
	request := &api.CreateSecurityGroupArgs{
		Name:  req.SecurityGroupName,
		Desc:  "",
		VpcId: req.VpcId,
	}

	r, err := p.bccClient.CreateSecurityGroup(request)
	if err != nil {
		return cloud.CreateSecurityGroupResponse{}, err
	} else {
		return cloud.CreateSecurityGroupResponse{
			SecurityGroupId: r.SecurityGroupId,
			RequestId:       "",
		}, nil
	}
}

// AddIngressSecurityGroupRule vpcid不用
func (p *BaiduCloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	request := &api.AuthorizeSecurityGroupArgs{
		Rule: &api.SecurityGroupRuleModel{
			SourceIp:        req.CidrIp,
			DestIp:          "",
			Protocol:        req.IpProtocol,
			SourceGroupId:   "",
			Ethertype:       "",
			PortRange:       fmt.Sprintf("%s-%s", strconv.Itoa(req.PortFrom), strconv.Itoa(req.PortTo)),
			DestGroupId:     "",
			SecurityGroupId: req.SecurityGroupId,
			Remark:          "",
			Direction:       "ingress", // required
		},
	}
	return p.bccClient.AuthorizeSecurityGroupRule(req.SecurityGroupId, request)
}

// AddEgressSecurityGroupRule vpcid不用
func (p *BaiduCloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	request := &api.AuthorizeSecurityGroupArgs{
		Rule: &api.SecurityGroupRuleModel{
			SourceIp:        "",
			DestIp:          req.CidrIp,
			Protocol:        req.IpProtocol,
			SourceGroupId:   "",
			Ethertype:       "",
			PortRange:       fmt.Sprintf("%s-%s", strconv.Itoa(req.PortFrom), strconv.Itoa(req.PortTo)),
			DestGroupId:     "",
			SecurityGroupId: req.SecurityGroupId,
			Remark:          "",
			Direction:       "egress", // required
		},
	}

	return p.bccClient.AuthorizeSecurityGroupRule(req.SecurityGroupId, request)
}

// maxkeys每页包含的最大数量，最大数量通常不超过1000，缺省值为1000。 缺少creatAt和RegionId
func (p *BaiduCloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {
	r, err := p.bccClient.ListSecurityGroup(&api.ListSecurityGroupArgs{
		Marker:     "",
		MaxKeys:    1000,
		InstanceId: "",
		VpcId:      req.VpcId,
	})
	if err != nil {
		return cloud.DescribeSecurityGroupsResponse{}, err
	} else {
		var groups []cloud.SecurityGroup
		for _, item := range r.SecurityGroups {
			groups = append(groups, cloud.SecurityGroup{
				SecurityGroupId:   item.Id,
				SecurityGroupType: "normal",
				SecurityGroupName: item.Name,
				CreateAt:          "",
				VpcId:             item.VpcId,
				RegionId:          req.RegionId,
			})
		}
		return cloud.DescribeSecurityGroupsResponse{
			Groups: groups,
		}, nil
	}
}

func (p *BaiduCloud) GetRegions() (cloud.GetRegionsResponse, error) {
	regions := cloud.GetRegionsResponse{Regions: []cloud.Region{
		{
			RegionId:  "BJ",
			LocalName: "华北-北京",
		},
		{
			RegionId:  "GZ",
			LocalName: "华南-广州",
		},
		{
			RegionId:  "SU",
			LocalName: "华东-苏州",
		},
		{
			RegionId:  "HKG",
			LocalName: "中国香港",
		},
		{
			RegionId:  "FWH",
			LocalName: "金融华中-武汉",
		},
		{
			RegionId:  "BD",
			LocalName: "华北-保定",
		},
	}}

	return regions, nil
}

func (p *BaiduCloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	r, err := p.bccClient.ListZone()
	if err != nil {
		return cloud.GetZonesResponse{}, err
	} else {
		var zones []cloud.Zone
		for _, item := range r.Zones {
			zones = append(zones, cloud.Zone{
				ZoneId:    item.ZoneName,
				LocalName: item.ZoneName,
			})
		}
		return cloud.GetZonesResponse{
			Zones: zones,
		}, nil
	}
}

// 缺失family,status设置为默认
func (p *BaiduCloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	zoneIds := make([]string, 0, 8)
	if req.ZoneId == "" {
		zones, err := p.GetZones(cloud.GetZonesRequest{})
		if err != nil {
			return cloud.DescribeAvailableResourceResponse{}, err
		}
		for _, zone := range zones.Zones {
			zoneIds = append(zoneIds, zone.ZoneId)
		}
	} else {
		zoneIds = append(zoneIds, req.ZoneId)
	}
	instanceTypes := make(map[string][]cloud.InstanceType)
	for _, zoneId := range zoneIds {
		r, err := p.bccClient.ListFlavorSpec(&api.ListFlavorSpecArgs{ZoneName: zoneId})
		if err != nil {
			return cloud.DescribeAvailableResourceResponse{}, err
		} else {

			for _, item := range r.ZoneResources {
				for _, flavor := range item.BccResources.FlavorGroups {
					for _, bbcFlavor := range flavor.Flavors {
						chargeType := _insTypeChargeType[bbcFlavor.ProductType]
						if chargeType == "" {
							continue
						}
						instanceTypes[zoneId] = append(instanceTypes[zoneId], cloud.InstanceType{
							ChargeType:  chargeType,
							IsGpu:       false,
							Core:        bbcFlavor.CpuCount,
							Memory:      bbcFlavor.MemoryCapacityInGB,
							Family:      bbcFlavor.SpecId,
							InsTypeName: bbcFlavor.Spec,
							Status:      cloud.InsTypeAvailable,
						})
					}
				}
			}
		}
	}
	return cloud.DescribeAvailableResourceResponse{
		InstanceTypes: instanceTypes,
	}, nil
}

// 缺失family,status设置为默认
func (p *BaiduCloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	return cloud.DescribeInstanceTypesResponse{}, nil
}

// maxkeys每页包含的最大数量，最大数量通常不超过1000
// TODO 通过instancetype找可用镜像  查了文档，没有这个接口
func (p *BaiduCloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	request := &api.ListImageArgs{
		Marker:    "",
		MaxKeys:   1000,
		ImageType: _imageType[req.ImageType],
		ImageName: "",
	}
	images := make([]cloud.Image, 0, 1000)
	r, err := p.bccClient.ListImage(request)
	if err != nil {
		return cloud.DescribeImagesResponse{}, err
	} else {
		for _, item := range r.Images {
			images = append(images, cloud.Image{
				Platform:  item.OsName,
				OsType:    item.OsType,
				OsName:    item.Name,
				Size:      0,
				ImageId:   item.Id,
				ImageName: item.Name,
			})
		}
		return cloud.DescribeImagesResponse{
			Images: images,
		}, nil
	}
}

// maxkeys每页包含的最大数量，最大数量通常不超过1000
// status设置为默认，createAt、regionId缺失
func (p *BaiduCloud) DescribeVpcs(req cloud.DescribeVpcsRequest) (cloud.DescribeVpcsResponse, error) {
	request := &vpc.ListVPCArgs{
		Marker:    "",
		MaxKeys:   1000,
		IsDefault: "",
	}
	vpcs := make([]cloud.VPC, 0, 1000)
	r, err := p.vpcClient.ListVPC(request)
	if err != nil {
		return cloud.DescribeVpcsResponse{}, err
	} else {
		for _, item := range r.VPCs {
			vpcs = append(vpcs, cloud.VPC{
				VpcId:     item.VPCID,
				VpcName:   item.Name,
				CidrBlock: item.Cidr,
				RegionId:  req.RegionId,
				Status:    cloud.VPCStatusAvailable,
				CreateAt:  "",
			})
		}
		return cloud.DescribeVpcsResponse{
			Vpcs: vpcs,
		}, nil
	}
}

// maxkeys每页包含的最大数量，最大数量通常不超过1000
// VsStatussh设置为默认，gatewayIpqu缺失
func (p *BaiduCloud) DescribeSwitches(req cloud.DescribeSwitchesRequest) (cloud.DescribeSwitchesResponse, error) {
	request := &vpc.ListSubnetArgs{
		Marker:     "",
		MaxKeys:    1000,
		VpcId:      req.VpcId,
		ZoneName:   "",
		SubnetType: "",
	}
	switchs := make([]cloud.Switch, 0, 1000)
	r, err := p.vpcClient.ListSubnets(request)
	if err != nil {
		return cloud.DescribeSwitchesResponse{}, err
	} else {
		var wg sync.WaitGroup
		var mutex sync.Mutex
		wg.Add(len(r.Subnets))
		for _, item := range r.Subnets {
			go func(switchId string) {
				defer func() {
					wg.Done()
					if e := recover(); e != nil {
						logs.Logger.Errorf("ShowServer %s panic, %v", switchId, e)
					}
				}()
				if switchId == "" {
					return
				}
				response, err := p.GetSwitch(cloud.GetSwitchRequest{
					SwitchId: switchId,
				})
				if err != nil {
					logs.Logger.Errorf("ShowServer failed, %s, %s", switchId, err.Error())
					return
				}
				mutex.Lock()
				switchs = append(switchs, response.Switch)
				mutex.Unlock()
			}(item.SubnetId)
		}
		wg.Wait()
		return cloud.DescribeSwitchesResponse{Switches: switchs}, nil
	}
}

// 缺失createAt和prefixListId
func (p *BaiduCloud) DescribeGroupRules(req cloud.DescribeGroupRulesRequest) (cloud.DescribeGroupRulesResponse, error) {
	request := &api.ListSecurityGroupArgs{
		Marker:     "",
		MaxKeys:    128,
		InstanceId: "",
		VpcId:      "",
	}
	rules := make([]cloud.SecurityGroupRule, 0, 128)
	r, err := p.bccClient.ListSecurityGroup(request)
	if err != nil {
		return cloud.DescribeGroupRulesResponse{}, nil
	} else {
		for _, item := range r.SecurityGroups {
			if item.Id == req.SecurityGroupId {
				for _, rule := range item.Rules {
					portResult := strings.Split(rule.PortRange, "-")
					portFrom, _ := strconv.Atoi(portResult[0])
					if err != nil {
						return cloud.DescribeGroupRulesResponse{}, err
					}
					portTo := portFrom
					if len(portResult) == 2 {
						portTo, _ = strconv.Atoi(portResult[1])
						if err != nil {
							return cloud.DescribeGroupRulesResponse{}, err
						}
					}
					var cidr string
					if rule.Direction == "egress" {
						cidr = rule.DestIp
					} else {
						cidr = rule.SourceIp
					}

					rules = append(rules, cloud.SecurityGroupRule{
						VpcId:           item.VpcId,
						SecurityGroupId: req.SecurityGroupId,
						PortFrom:        portFrom,
						PortTo:          portTo,
						Protocol:        rule.Protocol,
						Direction:       rule.Direction,
						GroupId:         item.Id,
						CidrIp:          cidr,
						PrefixListId:    "",
						CreateAt:        "",
					})
				}

			}
		}
	}
	return cloud.DescribeGroupRulesResponse{
		Rules: rules,
	}, nil
}

func (p *BaiduCloud) ListbyId(instanceIds []string) (instances []cloud.Instance, err error) {
	arg := &api.ListInstanceByInstanceIdArgs{
		InstanceIds: instanceIds,
	}
	res, err := p.bccClient.ListInstanceByInstanceIds(arg)
	if err != nil {
		return nil, err
	}
	var expireAt *time.Time
	for _, v := range res.Instances {
		if v.ExpireTime != "" {
			expireTime, _ := time.Parse("2006-01-02T15:04:05Z", v.ExpireTime)
			expireAt = &expireTime
		}
		instance := cloud.Instance{
			Id:       v.InstanceId,
			CostWay:  _ecsChargeType[v.PaymentTiming],
			Provider: cloud.BaiduCloud,
			IpInner:  v.InternalIP,
			IpOuter:  v.PublicIP,
			Network: &cloud.Network{
				VpcId:         v.VpcId,
				SubnetId:      v.SubnetId,
				SecurityGroup: "",
			},
			ImageId:  v.ImageId,
			Status:   _ecsStatus[string(v.Status)],
			ExpireAt: expireAt,
		}
		instances = append(instances, instance)
	}
	return instances, nil
}
func (p *BaiduCloud) CreateKeyPair(req cloud.CreateKeyPairRequest) (cloud.CreateKeyPairResponse, error) {
	input := &api.CreateKeypairArgs{
		Name: req.KeyPairName,
	}
	output, err := p.bccClient.CreateKeypair(input)
	if err != nil {
		logs.Logger.Errorf("CreateKeyPair BaiduCloud failed.err:[%v] req:[%v]", err, req)
		return cloud.CreateKeyPairResponse{}, err
	}
	return cloud.CreateKeyPairResponse{
		KeyPairId:   output.Keypair.KeypairId,
		KeyPairName: output.Keypair.Name,
		PrivateKey:  output.Keypair.PrivateKey,
		PublicKey:   output.Keypair.PublicKey,
	}, nil

}
func (p *BaiduCloud) ImportKeyPair(req cloud.ImportKeyPairRequest) (cloud.ImportKeyPairResponse, error) {
	input := &api.ImportKeypairArgs{
		Name:      req.KeyPairName,
		PublicKey: req.PublicKey,
	}
	output, err := p.bccClient.ImportKeypair(input)
	if err != nil {
		logs.Logger.Errorf("ImportKeyPair BaiduCloud failed.err:[%v] req:[%v]", err, req)
		return cloud.ImportKeyPairResponse{}, err
	}
	return cloud.ImportKeyPairResponse{
		KeyPairId:   output.Keypair.KeypairId,
		KeyPairName: output.Keypair.Name,
	}, nil
}
func (p *BaiduCloud) DescribeKeyPairs(req cloud.DescribeKeyPairsRequest) (cloud.DescribeKeyPairsResponse, error) {
	input := &api.ListKeypairArgs{
		Marker:  req.OlderMarker,
		MaxKeys: req.PageSize,
	}
	output, err := p.bccClient.ListKeypairs(input)
	if err != nil {
		logs.Logger.Errorf("DescribeKeyPairs BaiduCloud failed.err:[%v] req:[%v]", err, req)
		return cloud.DescribeKeyPairsResponse{}, err
	}
	totalCount := len(output.Keypairs)
	if totalCount == 0 {
		return cloud.DescribeKeyPairsResponse{}, nil
	}
	keyPairs := make([]cloud.KeyPair, 0, totalCount)
	for _, pair := range output.Keypairs {
		keyPairs = append(keyPairs, cloud.KeyPair{KeyPairId: pair.KeypairId, KeyPairName: pair.Name})
	}
	if output.IsTruncated {
		return cloud.DescribeKeyPairsResponse{
			TotalCount: totalCount,
			KeyPairs:   keyPairs,
			NewMarker:  output.NextMarker,
		}, nil
	} else {
		return cloud.DescribeKeyPairsResponse{
			TotalCount: totalCount,
			KeyPairs:   keyPairs,
		}, nil
	}
}
