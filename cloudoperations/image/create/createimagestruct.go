package imagecreate

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// CreateImageInput takes the required parameters for creating image.
// These parameters vary based on cloud choosed.
type CreateImageInput struct {
	InstanceIds []string `json:"instanceids"`
	Cloud       cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/create
