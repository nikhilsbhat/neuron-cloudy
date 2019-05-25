// Package servercreate makes the tool cloud agnostic with respect to creation of server.
// The decision will be made here to route the request to respective package based on input.
package servercreate

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// ServerCreateInput actually implements CreateServer and holds the data in creating server.
type ServerCreateInput struct {
	// Name that has to be assigend to the instance which will be created.
	InstanceName string `json:"instancename"`

	// The number of instance that has to be created.
	Count int64 `json:"count"`

	// Id of the base image that has to be used for creating the vm's.
	ImageId string `json:"imageid"`

	// Id of the subnetwork in which the instance/vm has to be created.
	SubnetId string `json:"subnetid"`

	// Name of the ssh keypair that has to used for creation of vm's, and this key will be used to log into the servers.
	KeyName string `json:"keyname"`

	// Flavor defines the hardware configurations of the vm that has to be created [ex: t2.micro(aws),n1-standard-2(gcp) etc.]
	Flavor string `json:"flavor"`

	// The raw codes that has to be executed immediatly after server boots up goes here.
	UserData string `json:"userdata"`

	// This defines whether a public IP has to be assigned to VM or not
	AssignPubIp bool `json:"assignpubip"`

	// All cloud info goes here
	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/create
