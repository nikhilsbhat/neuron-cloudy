// Package imagedelete makes the tool cloud agnostic for the image deletion process.
// The decision will be made here to route the request to respective package based on input.
package imagedelete

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// DeleteImageInput takes the required parameters for deleting images.
// These parameters vary based on cloud choosed.
type DeleteImageInput struct {
	// ImageIds are the list of image IDs which has to be deleted.
	ImageIds []string `json:"imageids"`
	Cloud    cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/delete
