// Package imagecreate makes the tool cloud agnostic for the image creation process.
// The decision will be made here to route the request to respective package based on input.
package imagecreate

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// CreateImageInput takes the required parameters for creating image.
// These parameters vary based on cloud choosed.
type CreateImageInput struct {
	// InstanceIds are the list of instance IDs of which an image would be captured.
	InstanceIds []string `json:"instanceids"`
	Cloud       cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/create
