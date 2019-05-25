package networkget

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	network "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetSubnetsResponse will return the filtered/unfiltered responses of variuos clouds.
type GetSubnetsResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse network.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

// GetSubnets will fetch the details of subnets specified else it pull the data out for all subnets in that particulat region
// and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (sub GetNetworksInput) GetSubnets() (GetSubnetsResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(sub.Cloud.Name)); status != true {
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetSubnets")
	}

	switch strings.ToLower(sub.Cloud.Name) {
	case "aws":

		creds, err := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: sub.Cloud.Profile,
				Cloud:   sub.Cloud.Name,
			},
		)

		if err != nil {
			return GetSubnetsResponse{}, err
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: sub.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authInpt := auth.EstablishConnectionInput{Region: sub.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call getsubnets and get the things done
		networkin := new(network.GetNetworksInput)
		networkin.GetRaw = sub.Cloud.GetRaw
		if sub.SubnetIds != nil {
			networkin.SubnetIds = sub.SubnetIds
			response, getSubErr := networkin.GetSubnets(authInpt)
			if getSubErr != nil {
				return GetSubnetsResponse{}, getSubErr
			}
			return GetSubnetsResponse{AwsResponse: response}, nil
		} else if sub.VpcIds != nil {
			networkin.VpcIds = sub.VpcIds
			response, getSubErr := networkin.GetSubnetsFromVpc(authInpt)
			if getSubErr != nil {
				return GetSubnetsResponse{}, getSubErr
			}
			return GetSubnetsResponse{AwsResponse: response}, nil
		} else {
			return GetSubnetsResponse{}, fmt.Errorf("You have not passed valid input to get details of server, the input struct looks like empty")
		}

	case "azure":
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return GetSubnetsResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetSubnets")
	}
}
