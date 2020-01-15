package aws

import (
	"fmt"
	"strings"

	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// DescribeInstanceInput implements various methods such as GetServersFromNetwork, GetServersFromSubnet, GetAllServers and GetServersDetails.
type DescribeInstanceInput struct {
	InstanceIds []string
	VpcIds      []string
	SubnetIds   []string
	Filters     Filters
	GetRaw      bool
}

// Filters holds the value for the filter to be applied on the servers.
type Filters struct {
	Name  string
	Value []string
}

// GetServersFromNetwork is tailored to fectch the servers from network, to fetch the data one has to pass either subnet-id else vpc-id to filters to make the life easy.
func (d *DescribeInstanceInput) GetServersFromNetwork(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, desInstErr := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)

	if desInstErr != nil {
		return nil, desInstErr
	}

	serverResponse := make([]ServerResponse, 0)

	if d.GetRaw == true {
		serverResponse = append(serverResponse, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return serverResponse, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			switch strings.ToLower(*instance.State.Name) {
			case "running":
				serverResponse = append(serverResponse, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PublicIpAddress: *instance.PublicIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "stopped":
				serverResponse = append(serverResponse, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "terminated":
				serverResponse = append(serverResponse, ServerResponse{State: *instance.State.Name, Cloud: "Amazon"})
			default:
				return nil, fmt.Errorf("Oops...!!!!. few instances are not in a state of fetching its details, check back after few seconds")
			}
		}
	}
	return serverResponse, nil
}

// GetServersFromSubnet will help in fetching the details of the server from the particular subnetwork paseed/selected.
func (d *DescribeInstanceInput) GetServersFromSubnet(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, desInstErr := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			Filters: aws.Filters{
				Name:  "subnet-id",
				Value: d.SubnetIds,
			},
		},
	)

	if desInstErr != nil {
		return nil, desInstErr
	}

	serverResponse := make([]ServerResponse, 0)

	if d.GetRaw == true {
		serverResponse = append(serverResponse, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return serverResponse, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {

			switch strings.ToLower(*instance.State.Name) {
			case "running":
				serverResponse = append(serverResponse, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PublicIpAddress: *instance.PublicIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "stopped":
				serverResponse = append(serverResponse, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "terminated":
				serverResponse = append(serverResponse, ServerResponse{State: *instance.State.Name, Cloud: "Amazon"})
			default:
				return nil, fmt.Errorf("Oops...!!!!. instances are not in a state of fetching the details of it, check back after few minutes")
			}
		}
	}
	return serverResponse, nil
}

// GetAllServers will fetch information of all the servers from a particular region.
func (d *DescribeInstanceInput) GetAllServers(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, err := ec2.DescribeAllInstances(
		&aws.DescribeComputeInput{},
	)
	if err != nil {
		return nil, err
	}

	serverResponse := make([]ServerResponse, 0)

	if d.GetRaw == true {
		serverResponse = append(serverResponse, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return serverResponse, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if (*instance.State.Name == "running") || (*instance.State.Name == "stopped") {
				serverResponse = append(serverResponse, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, InstanceType: *instance.InstanceType, Cloud: "Amazon", Region: *instance.Placement.AvailabilityZone})

			} else {
				// change has to be made here (introduction of omitempty is required)
				serverResponse = append(serverResponse, ServerResponse{State: "terminated", Cloud: "Amazon"})
			}
		}
	}
	return serverResponse, nil
}

//GetServersDetails will fetch the details of the random servers selected.
func (d *DescribeInstanceInput) GetServersDetails(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, desInstErr := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if desInstErr != nil {
		return nil, desInstErr
	}

	serverResponse := make([]ServerResponse, 0)

	if d.GetRaw == true {
		serverResponse = append(serverResponse, ServerResponse{GetInstRaw: result, Cloud: "Amazon"})
		return serverResponse, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {

			switch strings.ToLower(*instance.State.Name) {
			case "running":
				serverResponse = append(serverResponse, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PublicIpAddress: *instance.PublicIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "stopped":
				serverResponse = append(serverResponse, ServerResponse{InstanceName: *instance.Tags[0].Value, InstanceId: *instance.InstanceId, SubnetId: *instance.SubnetId, PrivateIpAddress: *instance.PrivateIpAddress, PrivateDnsName: *instance.PrivateDnsName, CreatedOn: (*instance.LaunchTime).String(), State: *instance.State.Name, Cloud: "Amazon"})
			case "terminated":
				serverResponse = append(serverResponse, ServerResponse{State: *instance.State.Name, Cloud: "Amazon"})
			default:
				return nil, fmt.Errorf("Oops...!!!!. instances are not in a state of fetching the details of it, check back after few minutes")
			}
		}
	}
	return serverResponse, nil
}
