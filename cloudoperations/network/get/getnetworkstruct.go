// Package networkget makes the tool cloud agnostic for fetching network details.
// The decision will be made here to route the request to respective package based on input.
package networkget

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetNetworksInput implements method GetNetworks, GetSubnets.
type GetNetworksInput struct {
	// Ids or names of VPC's of which the information has to be fetched.
	VpcIds []string `json:"vpcids"`

	// Ids or names of the SUBNET's of which the informaion has to be fetched.
	SubnetIds []string `json:"subnetids"`

	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for network/get.
