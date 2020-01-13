package networkget

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awscommon "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/common"
	awsnetwork "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
	gcp "github.com/nikhilsbhat/neuron-cloudy/cloud/gcp/operations"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetNetworks is responsible for fetching the details of a particular network passed
// or all the details of the networks present in the region.
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
		networkin := awsnetwork.GetNetworksInput{}
		networkin.VpcIds = net.NetworkID
		networkin.GetRaw = net.Cloud.GetRaw
		response, netErr := networkin.GetNetwork(authinpt)
		if netErr != nil {
			return GetNetworksResponse{}, netErr
		}
		return GetNetworksResponse{AwsResponse: response}, nil

	case "azure":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":

		fmt.Fprintf(os.Stdout, "%v\n", "We are in alpha for Google Cloud support, watchout for the output")
		getCluster := new(gcp.GetNetworkInput)
		getCluster.ProjectID = net.ProjectID
		getCluster.NetworkID = net.NetworkID[0]
		resp, err := getCluster.GetNetwork(net.Cloud.Client)
		if err != nil {
			return GetNetworksResponse{}, err
		}
		return GetNetworksResponse{GCPResponse: resp}, nil

	case "openstack":
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetNetworksResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}
}

// GetAllNetworks will fetch the details of all networks across all regions from the cloud specified.
func (net GetNetworksInput) GetAllNetworks() ([]GetNetworksResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(net.Cloud.Name)); status != true {
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllNetworks")
	}

	switch strings.ToLower(net.Cloud.Name) {
	case "aws":

		// Gets the established session so that it can carry out the process in cloud.
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

		networkResponse := make([]GetNetworksResponse, 0)
		for _, region := range regions.Regions {
			//authorizing to request further
			authinpt := auth.EstablishConnectionInput{Region: region, Resource: "ec2", Session: sess}

			networkin := awsnetwork.GetNetworksInput{GetRaw: net.Cloud.GetRaw}
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

		fmt.Fprintf(os.Stdout, "%v\n", "We are in alpha for Google Cloud support, watchout for the output")
		getCluster := new(gcp.GetNetworkInput)
		getCluster.ProjectID = net.ProjectID
		resp, err := getCluster.GetNetworks(net.Cloud.Client)
		if err != nil {
			return nil, err
		}
		networkResponse := make([]GetNetworksResponse, 0)
		return append(networkResponse, GetNetworksResponse{GCPResponse: resp}), nil

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
