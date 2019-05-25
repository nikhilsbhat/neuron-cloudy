package deleteserver

import (
	"fmt"
	"strings"

	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	server "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/server"
	awssess "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/sessions"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// DeleteServerResponse will return the filtered/unfiltered responses of variuos clouds.
type DeleteServerResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []server.ServerResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// DeleteServer will delete servers as per the parameter passed to it
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *DeleteServersInput) DeleteServer() (DeleteServerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
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
			return DeleteServerResponse{}, crederr
		}
		// I will establish session so that we can carry out the process in cloud
		sessionInput := awssess.CreateSessionInput{Region: serv.Cloud.Region, KeyId: creds.KeyId, AcessKey: creds.SecretAccess}
		sess := sessionInput.CreateAwsSession()

		//authorize
		authInpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call DeleteServer of interface and get the things done
		if serv.InstanceIds != nil {
			serverin := server.DeleteServerInput{}
			serverin.InstanceIds = serv.InstanceIds
			serverin.GetRaw = serv.Cloud.GetRaw
			serverResponse, serverr := serverin.DeleteServer(authInpt)
			if serverr != nil {
				return DeleteServerResponse{}, serverr
			}
			return DeleteServerResponse{AwsResponse: serverResponse}, nil
		} else if serv.VpcId != "" {
			serverin := server.DeleteServerInput{}
			serverin.VpcId = serv.VpcId
			serverin.GetRaw = serv.Cloud.GetRaw
			serverResponse, serverr := serverin.DeleteServerFromVpc(authInpt)
			if serverr != nil {
				return DeleteServerResponse{}, serverr
			}
			return DeleteServerResponse{AwsResponse: serverResponse}, nil
		} else {
			return DeleteServerResponse{}, fmt.Errorf("You have not passed valid input to get details of server, the input looks like empty")
		}

	case "azure":
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return DeleteServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "DeleteServer")
	}
}

// New returns the new DeleteServersInput instance with empty values
func New() *DeleteServersInput {
	net := &DeleteServersInput{}
	return net
}
