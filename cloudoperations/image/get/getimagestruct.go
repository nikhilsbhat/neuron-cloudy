package imagesget

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetImagesInput fetches the information about the images selected and is achieved by implementing the methoda for the same.
type GetImagesInput struct {
	ImageIds []string
	Cloud    cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for image/get
