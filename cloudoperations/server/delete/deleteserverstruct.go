package deleteserver

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// DeleteServersInput will implement DeleteServer nad holds data for deletion of server/vm
type DeleteServersInput struct {
	// Ids of instance/vm's that has to be deleted
	InstanceIds []string `json:"instanceids"`

	// Id of the network from which the server has to be deleted, be cautious while using this
	// because it actually terminates all the instances in the network if subnetwork is not metioned.
	VpcId string `json:"vpcid"`

	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/delete
