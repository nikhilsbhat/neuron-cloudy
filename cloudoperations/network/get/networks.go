// Package networkget makes the tool cloud agnostic for fetching network details.
// The decision will be made here to route the request to respective package based on input.
package networkget

import (
	awsnetwork "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
	gcp "github.com/nikhilsbhat/neuron-cloudy/cloud/gcp/operations"
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetNetworksInput implements method GetNetworks, GetSubnets.
type GetNetworksInput struct {
	// SubnetIds are IDs or names of the SUBNET's of which the informaion has to be fetched.
	SubnetIds []string `json:"subnetids"`
	// ProjectID refers to the ID of the project in which the selected resource exists.
	ProjectID string
	// NetworkID refers to the name/id of the network of which the information has to be retrived.
	NetworkID []string `json:"networkid"`
	Cloud     cmn.Cloud
}

// GetNetworksResponse will return the filtered/unfiltered responses of variuos clouds.
type GetNetworksResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []awsnetwork.NetworkResponse `json:"AwsResponse,omitempty"`
	// Contains filtered/unfiltered response of Azure.
	AzureResponse string                `json:"AzureResponse,omitempty"`
	GCPResponse   []gcp.NetworkResponse `json:"GcpResponse,omitempty"`
	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

//Nothing much from this file. This file contains only the structs for network/get.

// GetSubnetsResponse will return the filtered/unfiltered responses of variuos clouds.
type GetSubnetsResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse awsnetwork.NetworkResponse `json:"AwsResponse,omitempty"`
	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`
	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}
