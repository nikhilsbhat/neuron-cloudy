package networkupdate

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awsnetwork "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
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

		// Gets the established session so that it can carry out the process in cloud.
		sess := (net.Cloud.Client).(*session.Session)

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call UpdateNetwork of interface and get the things done
		serverin := awsnetwork.UpdateNetworkInput{}
		serverin.Resource = net.Catageory.Resource
		serverin.Action = net.Catageory.Action
		serverin.GetRaw = net.Cloud.GetRaw
		serverin.Network.Name = net.Catageory.Name
		serverin.Network.VpcCidr = net.Catageory.VpcCidr
		serverin.Network.VpcId = net.Catageory.VpcId
		serverin.Network.SubCidrs = net.Catageory.SubCidrs
		serverin.Network.Type = net.Catageory.Type
		serverin.Network.Ports = net.Catageory.Ports
		serverin.Network.Zone = net.Catageory.Zone

		response, err := serverin.UpdateNetwork(authinpt)
		if err != nil {
			return UpdateNetworkResponse{}, err
		}
		return UpdateNetworkResponse{AwsResponse: response}, nil

	case "azure":
		return UpdateNetworkResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return UpdateNetworkResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return UpdateNetworkResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return UpdateNetworkResponse{}, fmt.Errorf(common.DefaultCloudResponse + "NetworkUpdate")
	}
}

// New returns the new NetworkUpdateInput instance with empty values
func New() *NetworkUpdateInput {
	net := &NetworkUpdateInput{}
	return net
}
