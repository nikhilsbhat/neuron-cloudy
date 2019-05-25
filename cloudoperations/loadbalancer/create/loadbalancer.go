package createloadbalancer

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	loadbalance "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/loadbalancer"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// LoadBalanceResponse will return the filtered/unfiltered responses of variuos clouds.
type LoadBalanceResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse loadbalance.LoadBalanceResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// CreateLoadBalancer will create the loadbalancer based on the input in the struct LbCreateInput.
// Appropriate user and his cloud profile details has to be passed while calling it.
func (lb *LbCreateInput) CreateLoadBalancer() (LoadBalanceResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(lb.Cloud.Name)); status != true {
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(lb.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: lb.Cloud.Profile,
				Cloud:   lb.Cloud.Name,
			},
		)
		if err != nil {
			return LoadBalanceResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: lb.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = lb.Cloud.Region
		authinpt.Session = sess
		switch strings.ToLower(lb.Type) {
		case "classic":
			authinpt.Resource = "elb"
		case "application":
			authinpt.Resource = "elb2"
		}

		lbin := new(loadbalance.LoadBalanceCreateInput)
		lbin.GetRaw = lb.Cloud.GetRaw
		lbin.Name = lb.Name
		lbin.VpcId = lb.VpcId
		lbin.SubnetIds = lb.SubnetIds
		lbin.AvailabilityZones = lb.AvailabilityZones
		lbin.SecurityGroupIds = lb.SecurityGroupIds
		lbin.Scheme = lb.Scheme
		lbin.Type = lb.Type
		lbin.LbPort = lb.LbPort
		lbin.InstPort = lb.InstPort
		lbin.Lbproto = lb.Lbproto
		lbin.Instproto = lb.Instproto
		lbin.HttpCode = lb.HttpCode
		lbin.HealthPath = lb.HealthPath
		lbin.SslCert = lb.SslCert
		lbin.SslPolicy = lb.SslPolicy
		lbin.IpAddressType = lb.IpAddressType
		response, lberr := lbin.CreateLoadBalancer(*authinpt)
		if lberr != nil {
			return LoadBalanceResponse{}, lberr
		}
		return LoadBalanceResponse{AwsResponse: response}, nil

	case "azure":
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return LoadBalanceResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateLoadBalancer")
	}
}

// New returns the new instance of LbCreateInput with empty values.
func New() *LbCreateInput {
	net := &LbCreateInput{}
	return net
}
