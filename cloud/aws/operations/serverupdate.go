package aws

import (
	"fmt"
	"strings"

	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// UpdateServerInput implements UpdateServer and other related activities, it holds all necessary details to update server.
type UpdateServerInput struct {
	// InstanceIds are the list of instances which has to updated.
	InstanceIds []string
	// Action to be performned on instances as part of updation.
	Action string
	GetRaw bool
}

// UpdateServer updates the server (start/stop and other operations).
func (u *UpdateServerInput) UpdateServer(con aws.EstablishConnectionInput) ([]ServerResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	searchInput := CommonComputeInput{InstanceIds: u.InstanceIds}
	search, serverr := searchInput.SearchInstance(con)
	if serverr != nil {
		return nil, serverr
	}

	if search != true {
		return nil, fmt.Errorf("Could not find the entered Instances, please enter valid/existing InstanceIds")
	}
	serverResponse := make([]ServerResponse, 0)

	switch strings.ToLower(u.Action) {
	case "start":
		result, startErr := ec2.StartInstances(
			&aws.UpdateComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)

		if startErr != nil {
			return nil, startErr
		}

		waitErr := ec2.WaitTillInstanceRunning(
			&aws.DescribeComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)
		if waitErr != nil {
			return nil, waitErr
		}

		if u.GetRaw == true {
			serverResponse = append(serverResponse, ServerResponse{StartInstRaw: result, Cloud: "Amazon"})
			return serverResponse, nil
		}

		for _, inst := range result.StartingInstances {
			serverResponse = append(serverResponse, ServerResponse{InstanceId: *inst.InstanceId, CurrentState: "running", PreviousState: *inst.PreviousState.Name})
		}
		return serverResponse, nil

	case "stop":
		result, stopErr := ec2.StopInstances(
			&aws.UpdateComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)

		if stopErr != nil {
			return nil, stopErr
		}
		waitErr := ec2.WaitTillInstanceStopped(
			&aws.DescribeComputeInput{
				InstanceIds: u.InstanceIds,
			},
		)

		if waitErr != nil {
			return nil, waitErr
		}

		if u.GetRaw == true {
			serverResponse = append(serverResponse, ServerResponse{StopInstRaw: result, Cloud: "Amazon"})
			return serverResponse, nil
		}

		for _, inst := range result.StoppingInstances {
			serverResponse = append(serverResponse, ServerResponse{InstanceId: *inst.InstanceId, CurrentState: "stopped", PreviousState: *inst.PreviousState.Name})
		}
		return serverResponse, nil

	default:
		return nil, fmt.Errorf("Sorry...!!!!. I am not aware of the action you asked me to perform, please enter the action which we support. The available actions are: start/stop")
	}
}
