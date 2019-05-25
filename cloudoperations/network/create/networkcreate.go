package networkcreate

import (
	"fmt"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	network "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"

	//log "neuron/logger"
	"strings"
)

// CreateNetworkResponse is a struct that will return the filtered/unfiltered responses of variuos clouds.
type CreateNetworkResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse network.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// CreateNetwork is responsible for creating network and send back the response to the called source.
// appropriate user and his cloud profile details which was passed while calling it.
func (net *NetworkCreateInput) CreateNetwork() (CreateNetworkResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateNetwork")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: net.Cloud.Profile,
				Cloud:   net.Cloud.Name,
			},
		)

		if err != nil {
			return CreateNetworkResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: net.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// Fetching all the networks across cloud aws
		networkin := new(network.NetworkCreateInput)
		networkin.Name = net.Name
		networkin.VpcCidr = net.VpcCidr
		networkin.SubCidrs = net.SubCidr
		networkin.Type = net.Type
		networkin.Ports = net.Ports
		networkin.GetRaw = net.Cloud.GetRaw
		response, netErr := networkin.CreateNetwork(authinpt)
		if netErr != nil {
			return CreateNetworkResponse{}, netErr
		}
		return CreateNetworkResponse{AwsResponse: response}, nil

	case "azure":
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return CreateNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateNetwork")
	}
}

// New returns the new NetworkCreateInput instance with empty values
func New() *NetworkCreateInput {
	net := &NetworkCreateInput{}
	return net
}
