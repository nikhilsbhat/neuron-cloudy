// Package loadbalancer will help one in creating/deleting/updating/fetching loabalancers in aws.
// But this is tailor made for this application if one needs
// customized result, he/she has to write one similar to this for them by calling the (master) interface.
// This package is capable of returining both custom response
// and raw from cloud depending on what you pass.
package loadbalancer

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	network "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
)

// LoadBalanceCreateInput implements CreateLoadBalancer to create loadbalancer.
// It has various parameter which helps in taking decision for creating loabalancer and other aspects of it.
type LoadBalanceCreateInput struct {
	//optional parameter; If you provide the name to the loadbalancer well and good, else we will name it with a default one.
	Name string `json:"name"`

	//optional parameter; The Id of vpc in which the loadbalancer has to be created. Use this only if you don't want to pass subnets directly.
	//once this option is used we automatically fetch the random subnets from this network.
	VpcId string `json:"vpcid"`

	//optional parameter; The Ids of subnet in which the loadbalancer has to be created.
	SubnetIds []string `json:"subnetids"`

	//optional parameter; The names of availability zones to which loadbalancers has to be tagged. Either this or subnets has to be passed, passing both does't work
	AvailabilityZones []string `json:"availabilityzones"`

	//optional parameter; The Ids of secutiry group to be attached to loadbalancer.
	//If not mentioned, default security group of VPC will be attached.
	SecurityGroupIds []string `json:"securitygroupids"`

	//optional parameter; This field is to select the catageory of loadbalancer ex: internal, internet-facing. If not mentioned internet-facing will be created by default.
	Scheme string `json:"scheme"`
	//mandatory parameter; The type of loadbalancer required ex: classic, application etc.
	Type string `json:"type"`

	//required only if the LB protocol is HTTPS else can be initiazed with dummy value
	SslCert   string `json:"sslCert"`
	SslPolicy string `json:"sslPolicy"`

	//mandatory parameter; The port of the loabalancer. ex: 8080, 80 etc.
	LbPort   int64 `json:"lbport"`
	InstPort int64 `json:"instport"`

	//mandatory parameter; The protocol of loadbalancer. ex: HTTPS, HTTP.
	Lbproto   string `json:"lbproto"`
	Instproto string `json:"instproto"`

	//optional parameter; The http code. ex: 200, 404 etc.
	HttpCode   string `json:"httpcode"`
	HealthPath string `json:"healthpath"`

	//optional parameter; Select Ip address type ex: ipv4, ipv6. If nothing is passed ipv4 is considered by default.
	IpAddressType string `json:"ipaddresstype"`

	//optional parameter; Only when you need unfiltered result from cloud, enable this field by setting it to true. By default it is set to false.
	GetRaw bool `json:"getraw"`
}

// LoadBalanceResponse is the output format of CreateLoadBalancer, this holds both filetered and unfiletred response from cloud.
// But one has to enable flag 'GetRaw' in LoadBalanceCreateInput to get unfiletred output.
type LoadBalanceResponse struct {
	Name                   string                           `json:"name,omitempty"`
	Type                   string                           `json:"type,omitempty"`
	LbDns                  string                           `json:"lbdns,omitempty"`
	LbArn                  string                           `json:"lbarn,omitempty"`
	LbArns                 []string                         `json:"lbarns,omitempty"`
	TargetArn              interface{}                      `json:"targetarn,omitempty"`
	ListnerArn             interface{}                      `json:"listnerarn,omitempty"`
	Createdon              string                           `json:"createdon,omitempty"`
	VpcId                  string                           `json:"vpcid,omitempty"`
	Scheme                 string                           `json:"scheme,omitempty"`
	DefaultResponse        interface{}                      `json:"defaultresponse,omitempty"`
	LbDeleteStatus         string                           `json:"lbdeletestatus,omitempty"`
	ClassicLb              []LoadBalanceResponse            `json:"classiclb,omitempty"`
	ApplicationLb          []LoadBalanceResponse            `json:"applicationlb,omitempty"`
	CreateClassicLbRaw     *elb.CreateLoadBalancerOutput    `json:"createclassiclbraw,omitempty"`
	GetClassicLbsRaw       *elb.DescribeLoadBalancersOutput `json:"getclassiclbsraw,omitempty"`
	GetClassicLbRaw        *elb.LoadBalancerDescription     `json:"getclassiclbraw,omitempty"`
	CreateApplicationLbRaw ApplicationLbRaw                 `json:"createapplicationlbraw,omitempty"`
	GetApplicationLbRaw    ApplicationLbRaw                 `json:"getapplicationlbraw,omitempty"`
}

// ApplicationLbRaw is a not used alone and has to be used along with LoadBalanceResponse.
// This means, no function will which is accessible to user will send this in response but is passed via LoadBalanceResponse.
type ApplicationLbRaw struct {
	CreateApplicationLbRaw *elbv2.CreateLoadBalancerOutput    `json:"createapplicationlbraw,omitempty"`
	GetApplicationLbsRaw   *elbv2.DescribeLoadBalancersOutput `json:"getapplicationlbsraw,omitempty"`
	GetApplicationLbRaw    *elbv2.LoadBalancer                `json:"getapplicationlbraw,omitempty"`
	CreateTargetGroupRaw   *elbv2.CreateTargetGroupOutput     `json:"createtargetgroupraw,omitempty"`
	GetTargetGroupRaw      *elbv2.DescribeTargetGroupsOutput  `json:"gettargetgroupraw,omitempty"`
	CreateListnersRaw      *elbv2.CreateListenerOutput        `json:"createlistnersraw,omitempty"`
	GetListnersRaw         *elbv2.DescribeListenersOutput     `json:"getlistnersraw,omitempty"`
}

// CreateLoadBalancer is spoc for creating loadbalancer in aws and it can create both classica and applciation loabalancers.
// This return both custom and raw response, pass appropriate value to get the result.
func (load *LoadBalanceCreateInput) CreateLoadBalancer(con aws.EstablishConnectionInput) (LoadBalanceResponse, error) {

	// creating LB according to the input ex: application/classic
	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return LoadBalanceResponse{}, sesserr
	}

	lbin := new(aws.LoadBalanceCreateInput)
	//giving name to the loadbalancer which wil be created
	lbin.Name = load.Name
	// collecting subnet details
	if load.SubnetIds != nil {
		lbin.Subnets = load.SubnetIds
	} else {
		subnetsIn := network.GetNetworksInput{VpcIds: []string{load.VpcId}}
		subnetsResult, suberr := subnetsIn.GetSubnetsFromVpc(con)
		if suberr != nil {
			return LoadBalanceResponse{}, suberr
		}
		for _, subnet := range subnetsResult.Subnets {
			lbin.Subnets = append(lbin.Subnets, subnet.Id)
		}
	}

	// fetching security group so that I can attach it to the loabalancer which I am about to create
	if load.SecurityGroupIds != nil {
		lbin.SecurityGroups = load.SecurityGroupIds
	} else {
		secInput := network.NetworkComponentInput{VpcIds: []string{load.VpcId}}
		secResult, err := secInput.GetSecFromVpc(con)
		if err != nil {
			return LoadBalanceResponse{}, err
		}
		lbin.SecurityGroups = secResult.SecGroupIds
	}
	// creating load balancer

	// selecting scheme
	if load.Scheme == "external" {
		lbin.Scheme = "internet-facing"
	} else if load.Scheme == "internal" {
		lbin.Scheme = "internal"
	} else {
		lbin.Scheme = "internet-facing"
	}

	//setting availability zones
	if load.AvailabilityZones != nil {
		lbin.AvailabilityZones = load.AvailabilityZones
	}

	switch strings.ToLower(load.Type) {
	case "classic":

		lbin.InstPort = load.InstPort
		lbin.Instproto = load.Instproto
		lbin.LbPort = load.LbPort
		lbin.Lbproto = load.Lbproto
		lbin.SslCert = load.SslCert
		lbCreateResponse, err := elb.CreateClassicLb(*lbin)

		if err != nil {
			return LoadBalanceResponse{}, err
		}

		response := new(LoadBalanceResponse)
		if load.GetRaw == true {
			response.CreateClassicLbRaw = lbCreateResponse
			return *response, nil
		}

		response.Name = load.Name
		response.Type = load.Type
		response.LbDns = *lbCreateResponse.DNSName
		return *response, nil

	case "application":

		if load.IpAddressType == "" {
			lbin.IpAddressType = "ipv4"
		} else {
			lbin.IpAddressType = load.IpAddressType
		}
		// creating load balancer logic
		lbCreateResponse, lberr := elb.CreateApplicationLb(*lbin)
		if lberr != nil {
			return LoadBalanceResponse{}, lberr
		}

		lbin.Name = load.Name + "-target"
		lbin.LbPort = load.LbPort
		lbin.Lbproto = load.Lbproto
		lbin.VpcId = load.VpcId
		lbin.Instproto = load.Instproto
		lbin.InstPort = load.InstPort
		lbin.HealthPath = load.HealthPath
		lbin.HttpCode = load.HttpCode
		// creating target group
		targetGroupResponse, tarerr := elb.CreateTargetGroups(lbin)
		if tarerr != nil {
			return LoadBalanceResponse{}, tarerr
		}

		lbin.SslCert = load.SslCert
		lbin.TargetArn = *targetGroupResponse.TargetGroups[0].TargetGroupArn
		lbin.LbArn = *lbCreateResponse.LoadBalancers[0].LoadBalancerArn
		lbin.LbPort = load.LbPort
		lbin.Lbproto = load.Lbproto
		lbin.SslPolicy = load.SslPolicy
		listnerCreateResponse, liserr := elb.CreateApplicationListners(lbin)
		if liserr != nil {
			return LoadBalanceResponse{}, liserr
		}

		response := new(LoadBalanceResponse)

		if load.GetRaw == true {
			response.CreateApplicationLbRaw.CreateApplicationLbRaw = lbCreateResponse
			response.CreateApplicationLbRaw.CreateTargetGroupRaw = targetGroupResponse
			response.CreateApplicationLbRaw.CreateListnersRaw = listnerCreateResponse
			return *response, nil
		}

		response.Name = load.Name
		response.Type = load.Type
		response.LbDns = *lbCreateResponse.LoadBalancers[0].DNSName
		response.LbArn = *lbCreateResponse.LoadBalancers[0].LoadBalancerArn
		response.TargetArn = *targetGroupResponse.TargetGroups[0].TargetGroupArn
		response.ListnerArn = *listnerCreateResponse.Listeners[0].ListenerArn
		return *response, nil

	default:
		return LoadBalanceResponse{}, fmt.Errorf("You provided unknown loadbalancer type, enter a valid LB type")
	}
}
