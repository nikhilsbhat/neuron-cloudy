package deleteloadbalancer

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awslb "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// LoadBalancerDeleteResponse will return the filtered/unfiltered responses of variuos clouds.
type LoadBalancerDeleteResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []awslb.LoadBalanceDeleteResponse `json:"AwsResponse,omitempty"`
	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`
	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// DeleteLoadBalancer will help in deleting the loadbalancer created by CreateLoadBalancer
// Appropriate user and his cloud profile details has to be passed while calling it.
func (lb *LbDeleteInput) DeleteLoadBalancer() (LoadBalancerDeleteResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(lb.Cloud.Name)); status != true {
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(lb.Cloud.Name) {
	case "aws":

		// Gets the established session so that it can carry out the process in cloud
		sess := (lb.Cloud.Client).(*session.Session)

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

		lbin := new(awslb.DeleteLoadbalancerInput)
		lbin.LbNames = lb.LbNames
		lbin.LbArns = lb.LbArns
		lbin.Type = lb.Type
		lbin.GetRaw = lb.Cloud.GetRaw
		response, lberr := lbin.DeleteLoadbalancer(*authinpt)
		if lberr != nil {
			return LoadBalancerDeleteResponse{}, lberr
		}
		return LoadBalancerDeleteResponse{AwsResponse: response}, nil

	case "azure":
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteLoadBalancer")
	}
}

// New returns the new instance of LbDeleteInput with the empty default values.
func New() *LbDeleteInput {
	net := &LbDeleteInput{}
	return net
}
