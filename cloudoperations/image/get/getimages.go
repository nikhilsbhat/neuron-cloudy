package imagesget

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	image "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/image"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetImagesResponse contains the details of the images collected by GetImage.
// This also can contain the response from various cloud, but will deliver what was passed to it.
type GetImagesResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []image.ImageResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// GetImage will collect all the required information of the images specified to it and send back the response.
func (img *GetImagesInput) GetImage() (GetImagesResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(img.Cloud.Name)); status != true {
		return GetImagesResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetImages")
	}

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Cloud.Profile, Cloud: img.Cloud.Name})
		if crderr != nil {
			return GetImagesResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: img.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		getimage := new(image.GetImageInput)
		getimage.ImageIds = img.ImageIds
		getimage.GetRaw = img.Cloud.GetRaw
		result, err := getimage.GetImage(authinpt)
		if err != nil {
			return GetImagesResponse{}, err
		}
		return GetImagesResponse{AwsResponse: result}, nil

	case "azure":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetImagesResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetImage")
	}
}

// GetAllImage will fetch the details of all the images in the specified acoount ot region.
func (img *GetImagesInput) GetAllImage() (GetImagesResponse, error) {

	switch strings.ToLower(img.Cloud.Name) {
	case "aws":

		creds, crderr := common.GetCredentials(&common.GetCredentialsInput{Profile: img.Cloud.Profile, Cloud: img.Cloud.Name})
		if crderr != nil {
			return GetImagesResponse{}, crderr
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: img.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: img.Cloud.Region, Resource: "ec2", Session: sess}

		getimages := new(image.GetImageInput)
		getimages.GetRaw = img.Cloud.GetRaw
		result, err := getimages.GetAllImage(authinpt)
		if err != nil {
			return GetImagesResponse{}, err
		}
		return GetImagesResponse{AwsResponse: result}, nil

	case "azure":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetImagesResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetImagesResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetAllImage")
	}
}

// New returns the new instance of GetImagesInput with empty values.
func New() *GetImagesInput {
	net := &GetImagesInput{}
	return net
}
