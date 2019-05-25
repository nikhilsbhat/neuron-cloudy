package getloadbalancer

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetLoadbalancerInput takes all the required inputs for fetching the details of the loadbalancers.
type GetLoadbalancerInput struct {
	LbNames []string `json:"lbnames"`
	LbArns  []string `json:"lbarns"`
	Type    string   `json:"type"`
	Cloud   cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for loadbalance/get
