package imagedelete

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// DeleteImageInput takes the required parameters for deleting images.
// These parameters vary based on cloud choosed.
type DeleteImageInput struct {
	ImageIds []string `json:"imageids"`
	Cloud    cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/delete
