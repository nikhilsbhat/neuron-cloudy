package getservers

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetServersInput actually implements GetAllServers and GetServersDetails it holds the data for fetching server.
type GetServersInput struct {
	// Ids of the instance of which the data has to be fetched
	InstanceIds []string `json:"instanceids"`

	// Ids of the network from which the details of VM's has to be fetched. (make sure you pass either of them and not both)
	VpcIds []string `json:"vpcids"`

	// Ids of subnetwork from which the details of VM's has to be fetched.
	SubnetIds []string `json:"subnetids"`

	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/get
