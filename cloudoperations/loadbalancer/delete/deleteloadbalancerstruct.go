package deleteloadbalancer

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// LbDeleteInput takes the inputs required for deletion of appropriate loadbalancer. This supports various types of loadbalancer.
// These parameters vary based on cloud choosed.
type LbDeleteInput struct {
	LbNames []string `json:"lbnames"`
	LbArns  []string `json:"lbarns"`
	Type    string   `json:"type"`
	Cloud   cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for loadbalance/delete
