package networkupdate

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awsnetwork "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// UpdateNetworkResponse will return the filtered/unfiltered responses of variuos clouds.
type UpdateNetworkResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse awsnetwork.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// UpdateNetwork will update network and its components
// appropriate user and his cloud profile details which was passed while calling it.
func (net *NetworkUpdateInput) UpdateNetwork() (UpdateNetworkResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return UpdateNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "UpdateNetwork")
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
			return UpdateNetworkResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: net.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call UpdateNetwork of interface and get the things done
		serverin := awsnetwork.UpdateNetworkInput{
			Resource: net.Catageory.Resource,
			Action:   net.Catageory.Action,
			GetRaw:   net.Cloud.GetRaw,
			Network: awsnetwork.NetworkCreateInput{
				Name:     net.Catageory.Name,
				VpcCidr:  net.Catageory.VpcCidr,
				VpcId:    net.Catageory.VpcId,
				SubCidrs: net.Catageory.SubCidrs,
				Type:     net.Catageory.Type,
				Ports:    net.Catageory.Ports,
				Zone:     net.Catageory.Zone,
			},
		}
		response, err := serverin.UpdateNetwork(authinpt)
		if err != nil {
			return UpdateNetworkResponse{}, err
		}
		return UpdateNetworkResponse{AwsResponse: response}, nil

	case "azure":
		return UpdateNetworkResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return UpdateNetworkResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return UpdateNetworkResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:
		return UpdateNetworkResponse{DefaultResponse: common.DefaultCloudResponse + "NetworkUpdate"}, nil
	}
}

// New returns the new NetworkUpdateInput instance with empty values
func New() *NetworkUpdateInput {
	net := &NetworkUpdateInput{}
	return net
}
