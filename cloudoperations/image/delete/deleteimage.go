package imagedelete

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	image "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/image"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// DeleteImageResponse contains the details of the images deleted by DeleteImage.
// This also can contain the response from various cloud, but will deliver what was passed to it.
type DeleteImageResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []image.ImageResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// DeleteImage deletes the images based on the inputu passed via DeleteImageInput struct.
func (img *DeleteImageInput) DeleteImage() (DeleteImageResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(img.Cloud.Name)); status != true {
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteImage")
	}

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		// gets the established session so that we can carry out the process in cloud.
		sess := (img.Cloud.Client).(*session.Session)

		// authorizing further request
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		delimages := new(image.DeleteImageInput)
		delimages.ImageIds = img.ImageIds
		result, err := delimages.DeleteImage(authinpt)
		if err != nil {
			return DeleteImageResponse{}, err
		}
		response := make([]image.ImageResponse, 0)
		response = append(response, result)
		return DeleteImageResponse{AwsResponse: response}, nil

	case "azure":
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return DeleteImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetImage")
	}
}
