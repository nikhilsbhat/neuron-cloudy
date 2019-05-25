package updateservers

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	server "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/server"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// UpdateServersResponse will return the filtered/unfiltered responses of variuos clouds.
type UpdateServersResponse struct {

	// Contains filtered/unfiltered response of AWS.
	AwsResponse []server.ServerResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// UpdateServers will update servers (start/stop, change ebs etc)
//  with the instructions passed to him and give back the response who called this.
// Below method will take care of fetching details of
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *UpdateServersInput) UpdateServers() (UpdateServersResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return UpdateServersResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":

		creds, crederr := common.GetCredentials(
			&common.GetCredentialsInput{
				Profile: serv.Cloud.Profile,
				Cloud:   serv.Cloud.Name,
			},
		)

		if crederr != nil {
			return UpdateServersResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: serv.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call UpdateServer of interface and get the things done
		serverin := server.UpdateServerInput{InstanceIds: serv.InstanceIds, Action: serv.Action, GetRaw: serv.Cloud.GetRaw}
		response, err := serverin.UpdateServer(authinpt)
		if err != nil {
			return UpdateServersResponse{}, err
		}
		return UpdateServersResponse{AwsResponse: response}, nil

	case "azure":
		return UpdateServersResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return UpdateServersResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return UpdateServersResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:
		return UpdateServersResponse{}, fmt.Errorf(common.DefaultCloudResponse + "UpdateServers")
	}
}

// New returns the new UpdateServersInput instance with empty values
func New() *UpdateServersInput {
	net := &UpdateServersInput{}
	return net
}
