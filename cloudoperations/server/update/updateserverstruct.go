package updateservers

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// UpdateServersInput actually implements UpdateServers and it holds the data for updating server.
type UpdateServersInput struct {
	// Ids of the instances/vms which has to be updated
	InstanceIds []string `json:"instanceids"`

	// Action item that has to be performed on the VM
	Action string `json:"action"`

	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/update
