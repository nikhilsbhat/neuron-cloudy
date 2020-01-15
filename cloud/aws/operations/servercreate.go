package aws

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// CreateServerInput will implement the methods for creating instances and holds the value for the same.
type CreateServerInput struct {
	InstanceName string
	ImageId      string
	InstanceType string
	KeyName      string
	MaxCount     int64
	MinCount     int64
	SubnetId     string
	SecGroupId   string
	UserData     string
	AssignPubIp  bool
	GetRaw       bool
}

// ServerResponse holds the filtered/unfiltered output of CreateServer from aws.
type ServerResponse struct {
	InstanceName        string                        `json:"InstanceName,omitempty"`
	InstanceId          string                        `json:"InstanceId,omitempty"`
	SubnetId            string                        `json:"SubnetId,omitempty"`
	PrivateIpAddress    string                        `json:"IpAddress,omitempty"`
	PublicIpAddress     string                        `json:"PublicIpAddress,omitempty"`
	PrivateDnsName      string                        `json:"PrivateDnsName,omitempty"`
	CreatedOn           string                        `json:"CreatedOn,omitempty"`
	State               string                        `json:"State,omitempty"`
	InstanceDeleteState string                        `json:"InstanceDeleteState,omitempty"`
	InstanceType        string                        `json:"InstanceType,omitempty"`
	Cloud               string                        `json:"Cloud,omitempty"`
	Region              string                        `json:"Region,omitempty"`
	PreviousState       string                        `json:"PreviousState,omitempty"`
	CurrentState        string                        `json:"CurrentState,omitempty"`
	DefaultResponse     interface{}                   `json:"DefaultResponse,omitempty"`
	Error               error                         `json:"Error,omitempty"`
	CreateInstRaw       *ec2.DescribeInstancesOutput  `json:"CreateInstRaw,omitempty"`
	GetInstRaw          *ec2.DescribeInstancesOutput  `json:"DescribeInstRaw,omitempty"`
	DeleteInstRaw       *ec2.TerminateInstancesOutput `json:"DeleteInstRaw,omitempty"`
	StartInstRaw        *ec2.StartInstancesOutput     `json:"StartInstRaw,omitempty"`
	StopInstRaw         *ec2.StopInstancesOutput      `json:"StopInstRaw,omitempty"`
	CreateImgRaw        *ec2.CreateImageOutput        `json:"CreateImgRaw,omitempty"`
	DescribeImg         *ec2.DescribeImagesOutput     `json:"DescribeImg,omitempty"`
}

// CreateServer will help in creating instances/vms with the configuration passed.
func (csrv *CreateServerInput) CreateServer(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	// I will make a decision which security group to pick
	subInput := GetNetworksInput{SubnetIds: []string{csrv.SubnetId}}
	subResult, suberr := subInput.FindSubnet(con)
	if suberr != nil {
		return nil, suberr
	}

	if subResult != true {
		return nil, fmt.Errorf("Could not find the entered SUBNET, please enter valid/existing SUBNET id")
	}

	inst := new(aws.CreateServerInput)

	switch csrv.SecGroupId {
	case "":
		vpcRes, vpcerr := subInput.GetVpcFromSubnet(con)
		if vpcerr != nil {
			return nil, vpcerr
		}

		secInput := NetworkComponentInput{VpcIds: []string{vpcRes.VpcId}}
		secRes, secerr := secInput.GetSecFromVpc(con)
		if secerr != nil {
			return nil, nil
		}
		inst.SecurityGroups = secRes.SecGroupIds

	default:
		inst.SecurityGroups = []string{csrv.SecGroupId}
	}

	// I will be the spoc for the instance creation with the userdata passed to me
	switch csrv.UserData {
	case "":
		inst.UserData = b64.StdEncoding.EncodeToString([]byte("echo 'nothing'"))
	default:
		inst.UserData = b64.StdEncoding.EncodeToString([]byte(csrv.UserData))
	}

	switch csrv.MinCount {
	case 0:
		inst.MinCount = 1
	default:
		inst.MinCount = csrv.MinCount
	}

	switch csrv.MaxCount {
	case 0:
		inst.MaxCount = 1
	default:
		inst.MaxCount = csrv.MaxCount
	}

	inst.ImageId = csrv.ImageId
	inst.InstanceType = csrv.InstanceType
	inst.KeyName = csrv.KeyName
	inst.AssignPubIp = csrv.AssignPubIp
	inst.SubnetId = csrv.SubnetId
	// support for custom ebs mapping will be rolled out soon
	serverCreateResult, err := ec2.CreateInstance(inst)

	if err != nil {
		return nil, err
	}

	instanceIds := make([]string, 0)
	for _, instance := range serverCreateResult.Instances {
		instanceIds = append(instanceIds, *instance.InstanceId)
	}

	// I will make program wait until instance become running
	waitErr := ec2.WaitTillInstanceAvailable(
		&aws.DescribeComputeInput{
			InstanceIds: instanceIds,
		},
	)
	if waitErr != nil {
		return nil, waitErr
	}

	// creating tags for the server
	for i, instance := range instanceIds {
		tags := new(Tag)
		tags.Resource = instance
		tags.Name = "Name"
		tags.Value = csrv.InstanceName + "-" + strconv.Itoa(i)
		_, tagErr := tags.CreateTags(con)
		if tagErr != nil {
			return nil, tagErr
		}
	}

	//fetching the details of server
	result, serverr := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			InstanceIds: instanceIds,
		},
	)
	if serverr != nil {
		return nil, serverr
	}

	type serverResponse struct {
		name       string
		instanceId string
		ipaddress  string
		privatedns string
		publicIp   string
		createdon  string
	}

	response := make([]serverResponse, 0)
	createServerResponse := make([]ServerResponse, 0)

	if csrv.GetRaw == true {
		createServerResponse = append(createServerResponse, ServerResponse{CreateInstRaw: result, Cloud: "Amazon"})
		return createServerResponse, nil
	}

	// fetching the instance details which is created in previous process
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if csrv.AssignPubIp == true {
				response = append(response, serverResponse{name: *instance.Tags[0].Value, instanceId: *instance.InstanceId, ipaddress: *instance.PrivateIpAddress, privatedns: *instance.PrivateDnsName, publicIp: *instance.PublicIpAddress, createdon: (*instance.LaunchTime).String()})
			} else {
				response = append(response, serverResponse{name: *instance.Tags[0].Value, instanceId: *instance.InstanceId, ipaddress: *instance.PrivateIpAddress, privatedns: *instance.PrivateDnsName, createdon: (*instance.LaunchTime).String()})
			}
		}
	}

	for _, server := range response {
		createServerResponse = append(createServerResponse, ServerResponse{InstanceName: server.name, InstanceId: server.instanceId, SubnetId: csrv.SubnetId, PrivateIpAddress: server.ipaddress, PublicIpAddress: server.publicIp, PrivateDnsName: server.privatedns, CreatedOn: server.createdon, Cloud: "Amazon"})
	}

	return createServerResponse, nil
}
