package getloadbalancer

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	loadbalance "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/loadbalancer"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetLoadbalancerResponse will return the filtered/unfiltered responses of variuos clouds.
type GetLoadbalancerResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []loadbalance.LoadBalanceResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// GetLoadbalancers fetches the information of the appropriate loadbalancers.
// Appropriate user and his cloud profile details which was passed while calling it.
func (lb *GetLoadbalancerInput) GetLoadbalancers() (GetLoadbalancerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(lb.Cloud.Name)); status != true {
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(lb.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: lb.Cloud.Profile, Cloud: lb.Cloud.Name})
		if err != nil {
			return GetLoadbalancerResponse{}, err
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

		lbin := new(loadbalance.GetLoadbalancerInput)
		lbin.GetRaw = lb.Cloud.GetRaw
		lbin.LbNames = lb.LbNames
		lbin.LbArns = lb.LbArns
		lbin.Type = lb.Type
		response, lberr := lbin.Getloadbalancers(*authinpt)
		if lberr != nil {
			return GetLoadbalancerResponse{}, lberr
		}
		return GetLoadbalancerResponse{AwsResponse: response}, nil

	case "azure":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetLoadbalancers")
	}
}

// GetAllLoadbalancer fetches the information of all the loadbalancers in the specified region or the the one passed.
func (lb *GetLoadbalancerInput) GetAllLoadbalancer() (GetLoadbalancerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(lb.Cloud.Name)); status != true {
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(lb.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: lb.Cloud.Profile, Cloud: lb.Cloud.Name})
		if err != nil {
			return GetLoadbalancerResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: lb.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := new(auth.EstablishConnectionInput)
		authinpt.Region = lb.Cloud.Region
		authinpt.Session = sess
		authinpt.Resource = "elb12"
		lbin := new(loadbalance.GetLoadbalancerInput)
		lbin.GetRaw = lb.Cloud.GetRaw

		switch strings.ToLower(lb.Type) {
		case "classic":
			response, lberr := lbin.GetAllClassicLb(*authinpt)
			if lberr != nil {
				return GetLoadbalancerResponse{}, lberr
			}
			return GetLoadbalancerResponse{AwsResponse: response}, nil
		case "application":
			response, lberr := lbin.GetAllApplicationLb(*authinpt)
			if lberr != nil {
				return GetLoadbalancerResponse{}, lberr
			}
			return GetLoadbalancerResponse{AwsResponse: response}, nil
		case "":
			response, lberr := lbin.GetAllLoadbalancer(*authinpt)
			if lberr != nil {
				return GetLoadbalancerResponse{}, lberr
			}
			return GetLoadbalancerResponse{AwsResponse: response}, nil
		default:
			return GetLoadbalancerResponse{}, fmt.Errorf("The loadbalancer type you entered is unknown to me")
		}

	case "azure":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetLoadbalancerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetLoadbalancers")
	}
}

// New return the new instance of GetLoadbalancerInput with an empty values.
func New() *GetLoadbalancerInput {
	net := &GetLoadbalancerInput{}
	return net
}
