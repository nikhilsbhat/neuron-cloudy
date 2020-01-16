package servercreate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awsserver "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// ServerCreateResponse will return the filtered/unfiltered responses of variuos clouds.
type ServerCreateResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []awsserver.ServerResponse `json:"AwsResponse,omitempty"`
	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`
	// Default response if no inputs or matching the values required.
	DefaultResponse interface{} `json:"DefaultResponse,omitempty"`
}

// CreateServer will create the server with the requirement passed to him
// appropriate user and his cloud profile details which was passed while calling it.
func (serv ServerCreateInput) CreateServer() (ServerCreateResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}
	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":

		sess := (serv.Cloud.Client).(*session.Session)

		//authorizing to request further
		authInpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}

		// I will call CreateServer of interface and get the things done

		serverin := awsserver.CreateServerInput{}
		serverin.InstanceName = serv.InstanceName
		serverin.ImageId = serv.ImageId
		serverin.InstanceType = serv.Flavor
		serverin.KeyName = serv.KeyName
		serverin.MaxCount = serv.Count
		serverin.SubnetId = serv.SubnetId
		serverin.UserData = serv.UserData
		serverin.AssignPubIp = serv.AssignPubIp
		serverin.GetRaw = serv.Cloud.GetRaw
		response, err := serverin.CreateServer(authInpt)
		if err != nil {
			return ServerCreateResponse{}, err
		}
		return ServerCreateResponse{AwsResponse: response}, nil

	case "azure":
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return ServerCreateResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:

		return ServerCreateResponse{}, fmt.Errorf(common.DefaultCloudResponse + "CreateServer")
	}
}

// CreateServerMock will help the user to know what all parameter CreateServer takes as part of ServerCreateInput
func CreateServerMock() (ServerCreateResponse, error) {

	input := new(ServerCreateInput)
	defaultval := reflect.Indirect(reflect.ValueOf(input))

	defaults := make(map[string]interface{})
	for i := 0; i < defaultval.NumField(); i++ {
		defaults[defaultval.Type().Field(i).Name] = defaultval.Type().Field(i).Type
	}

	return ServerCreateResponse{DefaultResponse: defaults}, nil
}

// New returns the new ServerCreateInput instance with empty values
func New() *ServerCreateInput {
	net := &ServerCreateInput{}
	return net
}
