// Package networkupdate which makes the tool cloud agnostic for updating network components.
// The decision will be made here to route the request to respective package based on input.
package networkupdate

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// NetworkUpdateInput implements method GetNetworks, GetSubnets.
type NetworkUpdateInput struct {

	// The type of resources and the action to be performed in it
	// goes here the detailed inputs is in below struct
	Catageory

	Cloud cmn.Cloud
}

// Catageory holds the details of the network and its components which has to be updated.
type Catageory struct {
	// The type of resource that has to be updated.
	Resource string `json:"resource"`

	// Select the action that has to be performed on the resource
	// passed in above option.
	Action string `json:"action"`

	// Pass the name here for the resource that has to be created.
	Name string `json:"name"`

	// The CIDR block which will be used to create VPC and this
	// contains info that how many IP should be present in the network
	// so decide that in prior before calling this.
	VpcCidr string `json:"vpccidr"`

	// The CIDR for the subnet that has to be created in the VPC.
	// Pass an array of CIDR's and neuron will take care of creating
	// appropriate number of subnets and attaching to created VPC
	SubCidrs []string `json:"subcidrs"`

	// The type of the network that has to be created, public or private.
	// Accordingly IGW will be created and attached.
	Type string `json:"type"`

	// The ports that has to be opened for the network,
	// if not passed, by default 22 will be made open so that
	// one can access machines that will be created inside the created network.
	Ports []string `json:"ports"`

	// Pass the Id of the vpc here if you select to update a resource inside it.
	VpcId string `json:"vpcid"`

	// Pass the zone here if you need to create subnet in the required zone.
	Zone string `json:"zone"`
}

//Nothing much from this file. This file contains only the structs for network/update
