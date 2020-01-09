package networkget

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awscommon "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/common"
	network "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetNetworksResponse will return the filtered/unfiltered responses of variuos clouds.
type GetNetworksResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []network.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// GetNetworks is responsible for fetching the details of a particular network passed
// or all the details of the networks present in the region.
// appropriate user and his cloud profile details which was passed while calling it.
func (net *GetNetworksInput) GetNetworks() (GetNetworksResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		// Gets the established session so that it can carry out the process in cloud.
		sess := (net.Cloud.Client).(*session.Session)

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// Fetching all the networks across cloud aws
		networkin := network.GetNetworksInput{}
		networkin.VpcIds = net.VpcIds
		networkin.GetRaw = net.Cloud.GetRaw
		response, netErr := networkin.GetNetwork(authinpt)
		if netErr != nil {
			return GetNetworksResponse{}, netErr
		}
		return GetNetworksResponse{AwsResponse: response}, nil

	case "azure":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}
}

// GetAllNetworks will fetch the details of all networks across all regions from the cloud specified.
// appropriate user and his cloud profile details which was passed while calling it.
func (net GetNetworksInput) GetAllNetworks() ([]GetNetworksResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllNetworks")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		// Gets the establish session so that it can carry out the process in cloud.
		sess := (net.Cloud.Client).(*session.Session)

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: net.Cloud.Region, Resource: "ec2", Session: sess}

		// calls GetAllNetworks of interface and get the things done
		// Fetching all the regions from the cloud aws
		regionin := awscommon.CommonInput{}
		regions, regerr := regionin.GetRegions(authinpt)
		if regerr != nil {
			return nil, regerr
		}
		// Fetching all the networks across all the regions of cloud aws
		/*reg := make(chan []DengineAwsInterface.NetworkResponse, len(get_region_response.Regions))

		getnetworkdetails_input := GetAllNetworksInput{net.Cloud, net.Region}
		getnetworkdetails_input.getnetworkdetails(get_region_response.Regions, reg)
		for region_detail := range reg {
				get_all_network_response = append(get_all_network_response, GetAllNetworksResponse{AwsResponse: region_detail})
		}*/
		networkResponse := make([]GetNetworksResponse, 0)
		for _, region := range regions.Regions {
			//authorizing to request further
			authinpt := auth.EstablishConnectionInput{Region: region, Resource: "ec2", Session: sess}

			networkin := network.GetNetworksInput{GetRaw: net.Cloud.GetRaw}
			response, netErr := networkin.GetAllNetworks(authinpt)
			if netErr != nil {
				return nil, netErr
			}
			networkResponse = append(networkResponse, GetNetworksResponse{AwsResponse: response})
		}
		return networkResponse, nil

	case "azure":
		return nil, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return nil, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return nil, fmt.Errorf(common.DefaultOpResponse)
	default:
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllNetworks")
	}
}

// New returns the new GetNetworksInput instance with empty values
func New() *GetNetworksInput {
	net := &GetNetworksInput{}
	return net
}
