package miscoperations

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awscommon "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/common"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetRegionsResponse return the filtered/unfiltered responses of variuos clouds.
type GetRegionsResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse awscommon.CommonResponse `json:"Regions,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// GetRegions will fetch the information about the regions specified, else the details of entire region across the region.
func (reg *GetRegionInput) GetRegions() (GetRegionsResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(reg.Cloud.Name)); status != true {
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(reg.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: reg.Cloud.Profile,
				Cloud:   reg.Cloud.Name,
			},
		)
		if err != nil {
			return GetRegionsResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: reg.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: reg.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call create_vpc and get the things done
		regionin := awscommon.CommonInput{}
		regionin.GetRaw = reg.Cloud.GetRaw
		response, regErr := regionin.GetRegions(authinpt)
		if regErr != nil {
			return GetRegionsResponse{}, regErr
		}
		return GetRegionsResponse{AwsResponse: response}, nil

	case "azure":
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetRegionsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetRegions")
	}
}

// New will return the new instance of GetRegionInput with empty values.
func New() *GetRegionInput {
	net := &GetRegionInput{}
	return net
}
