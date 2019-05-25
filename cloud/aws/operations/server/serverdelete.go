package awsinstance

import (
	"fmt"

	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// DeleteServerInput will implement the methods to delete servers/vm and other related activities.
type DeleteServerInput struct {
	VpcId       string
	InstanceIds []string
	GetRaw      bool
}

// DeleteServer will delete the server who's Id was selected.
func (d *DeleteServerInput) DeleteServer(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	instanceSearchInput := CommonComputeInput{InstanceIds: d.InstanceIds}
	searchInstance, serverr := instanceSearchInput.SearchInstance(con)

	if serverr != nil {
		return nil, serverr
	}

	if searchInstance != true {
		return nil, fmt.Errorf("Could not find the entered Instances, please enter valid/existing InstanceIds")
	}
	deleteResult, insTermErr := ec2.DeleteInstance(
		&aws.DeleteComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if insTermErr != nil {
		return nil, insTermErr
	}

	waiterr := ec2.WaitTillInstanceTerminated(
		&aws.DescribeComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if waiterr != nil {
		return nil, waiterr
	}

	result, err := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			InstanceIds: d.InstanceIds,
		},
	)
	if err != nil {
		return nil, err
	}

	deleteResponse := make([]ServerResponse, 0)

	if d.GetRaw == true {
		deleteResponse = append(deleteResponse, ServerResponse{DeleteInstRaw: deleteResult, Cloud: "Amazon"})
		return deleteResponse, nil
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			deleteResponse = append(deleteResponse, ServerResponse{InstanceId: *instance.InstanceId, CurrentState: *instance.State.Name})
		}
	}
	return deleteResponse, nil
}

// DeleteServerFromVpc deletes every single instances from the vpc selected.
func (d *DeleteServerInput) DeleteServerFromVpc(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	instanceSearchInput := DescribeInstanceInput{
		VpcIds: []string{d.VpcId},
	}
	searchInstance, serverr := instanceSearchInput.GetServersFromNetwork(con)
	if serverr != nil {
		return nil, serverr
	}

	insatanceids := make([]string, 0)
	for _, instanceid := range searchInstance {
		insatanceids = append(insatanceids, instanceid.InstanceId)
	}

	result, servDelErr := ec2.DeleteInstance(
		&aws.DeleteComputeInput{
			InstanceIds: insatanceids,
		},
	)
	if servDelErr != nil {
		return nil, servDelErr
	}

	deleteResponse := make([]ServerResponse, 0)

	if d.GetRaw == true {
		deleteResponse = append(deleteResponse, ServerResponse{DeleteInstRaw: result, Cloud: "Amazon"})
		return deleteResponse, nil
	}

	for _, instance := range result.TerminatingInstances {
		deleteResponse = append(deleteResponse, ServerResponse{InstanceId: *instance.InstanceId, CurrentState: *instance.CurrentState.Name})
	}

	return deleteResponse, nil
}
