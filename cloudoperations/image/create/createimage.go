package imagecreate

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awsimage "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// CreateImageResponse contains the details of the images captured by CreateImage.
// This also can contain the response from various cloud, but will deliver what was passed to it.
type CreateImageResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []awsimage.ImageResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// CreateImage will capture image of the server specified, this gives back the response who called.
func (img *CreateImageInput) CreateImage() (CreateImageResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(img.Cloud.Name)); status != true {
		return CreateImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateImage")
	}

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		// gets the established session so that we can carry out the process in cloud.
		sess := (img.Cloud.Client).(*session.Session)

		// authorizing further request
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		responseImage := make([]awsimage.ImageResponse, 0)

		for _, id := range img.InstanceIds {
			imgcreate := new(awsimage.ImageCreateInput)
			imgcreate.InstanceId = id
			imgcreate.GetRaw = img.Cloud.GetRaw
			response, imgerr := imgcreate.CreateImage(authinpt)
			if imgerr != nil {
				return CreateImageResponse{}, imgerr
			}
			responseImage = append(responseImage, response)
		}
		return CreateImageResponse{AwsResponse: responseImage}, nil

	case "azure":
		return CreateImageResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return CreateImageResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return CreateImageResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return CreateImageResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateImage")
	}
}

// New returns the new instance of CreateImageInput with empty values.
func New() *CreateImageInput {
	net := &CreateImageInput{}
	return net
}
