package deleteloadbalancer

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// LbDeleteInput takes the inputs required for deletion of appropriate loadbalancer. This supports various types of loadbalancer.
// These parameters vary based on cloud choosed.
type LbDeleteInput struct {
	// LbNames refers to the list of loadbalancer that has to be deleted.
	LbNames []string `json:"lbnames"`
	// LbArns refers to the list of ARN values of the loadbalancer which has to be deleted.
	LbArns []string `json:"lbarns"`
	// Type refers to the type of loadbalancer of which has to be deleted.
	Type  string `json:"type"`
	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for loadbalance/delete
