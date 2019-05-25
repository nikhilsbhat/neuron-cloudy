package deleteloadbalancer

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	loadbalance "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/loadbalancer"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// LoadBalancerDeleteResponse will return the filtered/unfiltered responses of variuos clouds.
type LoadBalancerDeleteResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []loadbalance.LoadBalanceDeleteResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// DeleteLoadBalancer will help in deleting the loadbalancer created by CreateLoadBalancer
// Appropriate user and his cloud profile details has to be passed while calling it.
func (d *LbDeleteInput) DeleteLoadBalancer() (LoadBalancerDeleteResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(d.Cloud.Name)); status != true {
		return LoadBalancerDeleteResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(d.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: d.Cloud.Profile,
				Cloud:   d.Cloud.Name,
			},
		)
		if err != nil {
			return LoadBalancerDeleteResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: d.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = d.Cloud.Region
		authinpt.Session = sess
		switch strings.ToLower(d.Type) {
		case "classic":
			authinpt.Resource = "elb"
		case "application":
			authinpt.Resource = "elb2"
		}

		lbin := new(loadbalance.DeleteLoadbalancerInput)
		lbin.LbNames = d.LbNames
		lbin.LbArns = d.LbArns
		lbin.Type = d.Type
		lbin.GetRaw = d.Cloud.GetRaw
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
