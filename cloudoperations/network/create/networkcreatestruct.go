// Package networkcreate makes the tool cloud agnostic with respect to creation of network.
// The decision will be made here to route the request to respective package based on input.
package networkcreate

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// NetworkCreateInput implements method CreateNetwork and holds parameter for creating network.
type NetworkCreateInput struct {
	// Name for the Network that has to be created.
	Name string `json:"name"`
	// VpcCidr is the CIDR block which will be used to create VPC and this
	// contains info that how many IP should be present in the network
	// so decide that in prior before calling this.
	VpcCidr string `json:"vpccidr"`
	// SubCidr refers to the list of CIDR for the subnet that has to be created in the VPC.
	// Pass an array of CIDR's and neuron will take care of creating
	// appropriate number of subnets and attaching to created VPC.
	SubCidr []string `json:"subcidr"`
	// Type of the network that has to be created, public or private.
	// Accordingly IGW will be created and attached.
	Type string `json:"type"`
	// Ports that has to be opened for the network,
	// if not passed, by default 22 will be made open so that
	// one can access machines that will be created inside the created network.
	Ports []string `json:"ports"`
	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for network/create
