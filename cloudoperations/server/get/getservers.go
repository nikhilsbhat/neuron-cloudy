package getservers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	auth "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	awscommon "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/common"
	server "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/server"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetServerResponse will return the filtered/unfiltered responses of variuos clouds.
type GetServerResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse []server.ServerResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"Response,omitempty"`
}

// GetServersDetails will fetch the details of servers with the instructions passed to it
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *GetServersInput) GetServersDetails() (GetServerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return GetServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":

		// I will establish session so that we can carry out the process in cloud
		sess := (serv.Cloud.Client).(*session.Session)

		//authorizing to request further
		authinpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}
		// I will call CreateServer of interface and get the things done

		if serv.InstanceIds != nil {
			serverin := server.DescribeInstanceInput{}
			serverin.InstanceIds = serv.InstanceIds
			serverin.GetRaw = serv.Cloud.GetRaw
			serverResponse, serverr := serverin.GetServersDetails(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: serverResponse}, nil
		} else if serv.SubnetIds != nil {
			serverin := server.DescribeInstanceInput{}
			serverin.SubnetIds = serv.SubnetIds
			serverin.GetRaw = serv.Cloud.GetRaw
			serverResponse, serverr := serverin.GetServersFromSubnet(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: serverResponse}, nil
		} else if serv.VpcIds != nil {
			serverin := server.DescribeInstanceInput{}
			serverin.VpcIds = serv.VpcIds
			serverin.GetRaw = serv.Cloud.GetRaw
			serverResponse, serverr := serverin.GetServersFromNetwork(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: serverResponse}, nil
		} else {
			serverin := server.DescribeInstanceInput{GetRaw: serv.Cloud.GetRaw}
			serverResponse, serverr := serverin.GetAllServers(authinpt)
			if serverr != nil {
				return GetServerResponse{}, serverr
			}
			return GetServerResponse{AwsResponse: serverResponse}, nil
		}
	case "azure":
		return GetServerResponse{DefaultResponse: common.DefaultAzResponse}, nil
	case "gcp":
		return GetServerResponse{DefaultResponse: common.DefaultGcpResponse}, nil
	case "openstack":
		return GetServerResponse{DefaultResponse: common.DefaultOpResponse}, nil
	default:
		return GetServerResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetServers")
	}
}

// GetAllServers will fetch the details of all servers across the cloud
// appropriate user and his cloud profile details which was passed while calling it.
func (serv *GetServersInput) GetAllServers() ([]GetServerResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(serv.Cloud.Name)); status != true {
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetNetworks")
	}

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":

		// Gets the established session so that it can carry out the process in cloud.
		sess := (serv.Cloud.Client).(*session.Session)

		//authorize
		authinpt := auth.EstablishConnectionInput{Region: serv.Cloud.Region, Resource: "ec2", Session: sess}

		// Fetching list of regions to get details  of server across the account
		regionin := awscommon.CommonInput{}
		regions, regerr := regionin.GetRegions(authinpt)
		if regerr != nil {
			return nil, regerr
		}

		reg := make(chan []server.ServerResponse, len(regions.Regions))
		serv.getservers(regions.Regions, reg)

		serverResponse := make([]GetServerResponse, 0)
		for regionDetail := range reg {
			if len(regionDetail) != 0 {
				serverResponse = append(serverResponse, GetServerResponse{AwsResponse: regionDetail})
			}
		}
		return serverResponse, nil

	case "azure":
		return nil, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":
		return nil, fmt.Errorf(common.DefaultGcpResponse)
	case "openstack":
		return nil, fmt.Errorf(common.DefaultOpResponse)
	default:
		return nil, fmt.Errorf(common.DefaultCloudResponse + "GetAllServers")
	}
}

// this will be called by getallservers, he is the one who gets the details of all the servers,
// and send over a channel.
func (serv *GetServersInput) getservers(regions []string, reg chan []server.ServerResponse) {

	switch strings.ToLower(serv.Cloud.Name) {
	case "aws":
		// I will establish session so that we can carry out the process in cloud
		sess := (serv.Cloud.Client).(*session.Session)
		var wg sync.WaitGroup
		wg.Add(len(regions))
		for _, region := range regions {
			go func(region string) {
				defer wg.Done()

				//authorize
				authinpt := auth.EstablishConnectionInput{Region: region, Resource: "ec2", Session: sess}
				serverin := server.DescribeInstanceInput{GetRaw: serv.Cloud.GetRaw}
				serverResponse, _ := serverin.GetAllServers(authinpt)
				reg <- serverResponse
			}(region)
		}
		wg.Wait()
		close(reg)

	case "azure":
	case "gcp":
	case "openstack":
	default:
	}
}

// New returns the new GetServersInput instance with empty values
func New() *GetServersInput {
	net := &GetServersInput{}
	return net
}
