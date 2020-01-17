package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// LoadBalanceCreateInput implements CreateLoadBalancer to create loadbalancer.
// It has various parameter which helps in taking decision for creating loabalancer and other aspects of it.
type LoadBalanceCreateInput struct {
	// Name refers to the name of the loadbalancer to be created or of which the information to be retrieved.
	// optional parameter; If you provide the name to the loadbalancer well and good, else we will name it with a default one.
	Name string
	// VpcId is the ID of the network of which the loadbalancer is part of.
	// optional parameter; The Id of vpc in which the loadbalancer has to be created. Use this only if you don't want to pass subnets directly.
	// once this option is used we automatically fetch the random subnets from this network.
	VpcId string
	// SubnetIds is an array of subnetworks which wpuld be part of the network retrieved.
	// optional parameter;
	SubnetIds []string
	// AvailabilityZones is an array of availability zones to which loadbalancers has to be tagged.
	// optional parameter;
	AvailabilityZones []string
	// SecurityGroupIds are the Ids of secutiry group to be attached to loadbalancer.
	// optional parameter; If not mentioned, default security group of VPC will be attached.
	SecurityGroupIds []string
	// Scheme is to select the catageory of loadbalancer ex: internal, internet-facing. If not mentioned internet-facing will be created by default.
	// optional parameter;
	Scheme string
	// Type are the type of loadbalancer required ex: classic, application etc.
	// mandatory parameter;
	Type string
	// SslCert takes the link to the certificate which will be used to loadbalancer.
	// required only if the LB protocol is HTTPS else can be initiazed with dummy value
	SslCert string
	// SslPolicy defines the policy for the ssl.
	SslPolicy string
	// LbPort is the port of loadbalacner to be opened with.
	// mandatory parameter; The port of the loabalancer. ex: 8080, 80 etc.
	LbPort int64
	// InstPort is the port on which loadbalacner has to talk to instance.
	InstPort int64
	// Lbproto defines the protocol of loadbalancer.
	// mandatory parameter; The protocol of loadbalancer. ex: HTTPS, HTTP.
	Lbproto string
	// Instproto defines the protocol of instance.
	Instproto string
	// HttpCode has the http code. ex: 200, 404 etc.
	// optional parameter;
	HttpCode string
	// HealthPath is the path for loadbalacer to check the healt of the backend systems.
	HealthPath string
	// IpAddressType refers to Ip address type ex: ipv4, ipv6. If nothing is passed ipv4 is considered by default.
	// optional parameter;
	IpAddressType string
	// GetRaw returns unfiltered response from the cloud if it is set to true.
	// optional parameter;
	GetRaw bool
}

// LoadBalanceResponse is the output format of CreateLoadBalancer, this holds both filetered and unfiletred response from cloud.
// But one has to enable flag 'GetRaw' in LoadBalanceCreateInput to get unfiletred output.
type LoadBalanceResponse struct {
	// Name refers to the name of the loadbalancer of which the information has to be retrieved.
	Name string `json:"name,omitempty"`
	// Type are the type of loadbalancer required ex: classic, application etc.
	Type string `json:"type,omitempty"`
	// LbDns is the DNS associated to the loadbalancer created/retrieved.
	LbDns string `json:"lbdns,omitempty"`
	// LbArn is the ARN value associated to the loadbalancer created/retrieved.
	LbArn string `json:"lbarn,omitempty"`
	// LbArns are the array of ARN value associated to the loadbalancer created/retrieved.
	LbArns []string `json:"lbarns,omitempty"`
	// TargetArn refers to the ARN value associated to the target group part of loadbalancer.
	TargetArn interface{} `json:"targetarn,omitempty"`
	// ListnerArn refers to the ARN value associated to the listner group part of loadbalancer.
	ListnerArn interface{} `json:"listnerarn,omitempty"`
	// Createdon refers to the date on which the loadbalancer is created.
	Createdon string `json:"createdon,omitempty"`
	// VpcId is the ID of the network of which the loadbalancer is part of.
	VpcId string `json:"vpcid,omitempty"`
	// Scheme is to select the catageory of loadbalancer ex: internal, internet-facing. If not mentioned internet-facing will be created by default.
	Scheme          string      `json:"scheme,omitempty"`
	DefaultResponse interface{} `json:"defaultresponse,omitempty"`
	LbDeleteStatus  string      `json:"lbdeletestatus,omitempty"`
	// ClassicLb has the responses of classic loadbalancer.
	ClassicLb []LoadBalanceResponse `json:"classiclb,omitempty"`
	// ApplicationLb has the responses of application loadbalancer.
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
		subnetsIn := GetNetworksInput{VpcIds: []string{load.VpcId}}
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
		secInput := NetworkComponentInput{VpcIds: []string{load.VpcId}}
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

		if len(load.IpAddressType) == 0 {
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
