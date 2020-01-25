// Package miscoperations makes the tool cloud agnostic in creation/deletion/updation/retrieval of other cloud resources.
// The decision will be made here to route the request to respective package based on input.
package miscoperations

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetRegionInput holds the information required to fetch the details of region, and implements the methods for it
type GetRegionInput struct {
	Cloud cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for server/get
