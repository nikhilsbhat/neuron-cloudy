package networkdelete

import (
	"fmt"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	network "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"

	"strings"
)

// DeleteNetworkResponse returns the filtered/unfiltered responses of variuos clouds.
type DeleteNetworkResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse network.DeleteNetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// DeleteNetwork will help in deleting network and its components.
// Appropriate user and his cloud profile details which was passed while calling it.
func (net *DeleteNetworkInput) DeleteNetwork() (DeleteNetworkResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "UpdateNetwork")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(&common.GetCredentialsInput{Profile: net.Cloud.Profile, Cloud: net.Cloud.Name})
		if err != nil {
			return DeleteNetworkResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: net.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// deleteting network from aws
		networkin := new(network.DeleteNetworkInput)
		networkin.VpcIds = net.VpcIds
		networkin.GetRaw = net.Cloud.GetRaw
		response, netErr := networkin.DeleteNetwork(authinpt)
		if netErr != nil {
			return DeleteNetworkResponse{}, netErr
		}
		return DeleteNetworkResponse{AwsResponse: response}, nil

	case "azure":
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return DeleteNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteNetwork")
	}
}

// New returns the new instance of DeleteNetworkInput with empty values
func New() *DeleteNetworkInput {
	net := &DeleteNetworkInput{}
	return net
}
