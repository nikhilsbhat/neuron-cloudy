package getloadbalancer

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetLoadbalancerInput takes all the required inputs for fetching the details of the loadbalancers.
type GetLoadbalancerInput struct {
	// LbNames refers to the list of loadbalancer names of which the information has to be retrieved.
	LbNames []string `json:"lbnames"`
	// LbArns refers to the list of ARN values of the loadbalancer has to be retrieved.
	LbArns []string `json:"lbarns"`
	// Type refers to the type of loadbalancer of which the information has to be retrieved.
	Type  string `json:"type"`
	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for loadbalance/get
