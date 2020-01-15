package neuronaws

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateServerInput holds all required values to create a server/instance in aws.
type CreateServerInput struct {
	// ImageId is the ID of the image of which the information has to be retrived.
	ImageId string
	// InstanceType of the the instance that has to be created ex: t2.micro, t2.medium etc.
	InstanceType string
	// KeyName refers to the name of the key-value pair that has to be assossiated with the instance that would be created.
	KeyName string
	// MaxCount of the instance to be created (defaults to 1).
	MaxCount int64
	// MinCount of the instance to be created (defaults to 1).
	MinCount int64
	// SubnetId is the ID of the subnetwork in which the instance has to be created.
	SubnetId string
	// SecurityGroups is the ID of the securoty group which has to be associated with instance.
	SecurityGroups []string
	// UserData are the command/script has to be passed while instance bootup.
	UserData string
	// AssignPubIp is the deciding factor for opening instance to public (defaults to false making instance accessible only at private network).
	AssignPubIp bool
}

// DescribeComputeInput holds all the required values to describe the instance/vm or any compute resources in aws.
type DescribeComputeInput struct {
	// InstanceIds are the IDs of the instances of whom the information has to be retrived.
	InstanceIds []string
	// ImageIds are the IDs of the images of whom the information has to be retrived.
	ImageIds []string
	// Filters can be applied on the resource to fetch more appropriate information.
	Filters Filters
}

// UpdateComputeInput holds all the required values to update the compute resources in aws.
type UpdateComputeInput struct {
	// InstanceIds are the IDs of the instances which reqires updation.
	InstanceIds []string
	// Force will forcefully apply the updates on the resource.
	Force string
}

// DeleteComputeInput holds the required details for deleting compute resource.
type DeleteComputeInput struct {
	// ImageId is the ID of the image which has to be deleted.
	ImageId string
	// SnapshotId is the ID of the snapshot associated to an image which has to be deleted.
	SnapshotId string
	// InstanceIds are the IDs of the instances that has to be deleted.
	InstanceIds []string
}

// ImageCreateInput holds the details for creation of imagfe in aws.
type ImageCreateInput struct {
	// Description to the image would be created.
	Description string
	// ServerName is the name of the instance who's image has to be captured (this would be fetched automatically by other apis).
	ServerName string
	// InstanceId is the ID of the instance of which the image has to be captured.
	InstanceId string
}

// CreateInstance will create instance/vm as per the configuration specified.
func (sess *EstablishedSession) CreateInstance(ins *CreateServerInput) (*ec2.Reservation, error) {

	if sess.Ec2 != nil {
		if (ins.ImageId != "") || (ins.InstanceType != "") || (ins.KeyName != "") || (ins.MinCount != 0) || (ins.MaxCount != 0) || (ins.UserData != "") || (ins.SubnetId != "") || (ins.SecurityGroups != nil) {
			// support for custom ebs mapping will be rolled out soon
			createServerInput := &ec2.RunInstancesInput{
				ImageId:      aws.String(ins.ImageId),
				InstanceType: aws.String(ins.InstanceType),
				KeyName:      aws.String(ins.KeyName),
				MaxCount:     aws.Int64(ins.MaxCount),
				MinCount:     aws.Int64(ins.MinCount),
				UserData:     aws.String(ins.UserData),
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{{
					AssociatePublicIpAddress: aws.Bool(ins.AssignPubIp),
					DeviceIndex:              aws.Int64(0),
					DeleteOnTermination:      aws.Bool(true),
					SubnetId:                 aws.String(ins.SubnetId),
					Groups:                   aws.StringSlice(ins.SecurityGroups),
				}},
			}
			serverCreateResult, err := (sess.Ec2).RunInstances(createServerInput)
			// handling the error if it throws while subnet is under creation process
			if err != nil {
				return nil, err
			}
			return serverCreateResult, nil
		}
		return nil, fmt.Errorf("You provided empty/wrong details to CreateInstance, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

// DescribeInstance will help in fetching the information about the instance selected, by describing it.
func (sess *EstablishedSession) DescribeInstance(des *DescribeComputeInput) (*ec2.DescribeInstancesOutput, error) {

	if sess.Ec2 != nil {
		if des.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(des.InstanceIds),
			}
			result, err := (sess.Ec2).DescribeInstances(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(des.Filters, Filters{}) {
			return nil, fmt.Errorf("You provided empty struct to DescribeInstance, this is not acceptable")
		}
		if (des.Filters.Name == "") || (des.Filters.Value == nil) {
			return nil, fmt.Errorf("You chose Filters to fetch server details and did not provided required value for Filters")
		}
		input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{Name: aws.String(des.Filters.Name),
					Values: aws.StringSlice(des.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeInstances(input)

		if err != nil {
			return nil, err
		}
		return result, nil

	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// DescribeAllInstances will describe all the av available instance in the specified region in aws.
func (sess *EstablishedSession) DescribeAllInstances(des *DescribeComputeInput) (*ec2.DescribeInstancesOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeInstancesInput{}
		result, err := (sess.Ec2).DescribeInstances(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// DeleteInstance will delete the instance who's Id is specified.
func (sess *EstablishedSession) DeleteInstance(d *DeleteComputeInput) (*ec2.TerminateInstancesOutput, error) {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			terminateInstanceInput := &ec2.TerminateInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			_, err := (sess.Ec2).TerminateInstances(terminateInstanceInput)

			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		return nil, fmt.Errorf("You provided empty struct to DeleteInstance, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// StartInstances helps in stating the stopped instances in aws.
func (sess *EstablishedSession) StartInstances(s *UpdateComputeInput) (*ec2.StartInstancesOutput, error) {

	if sess.Ec2 != nil {
		if s.InstanceIds != nil {
			input := &ec2.StartInstancesInput{
				InstanceIds: aws.StringSlice(s.InstanceIds),
			}
			result, err := (sess.Ec2).StartInstances(input)

			if err != nil {
				return nil, err
			}
			return result, nil

		}
		return nil, fmt.Errorf("You provided empty struct to StartInstances, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// StopInstances will help in stopping the currently runnig instance.
func (sess *EstablishedSession) StopInstances(s *UpdateComputeInput) (*ec2.StopInstancesOutput, error) {

	if sess.Ec2 != nil {
		if s.InstanceIds != nil {
			input := &ec2.StopInstancesInput{
				InstanceIds: aws.StringSlice(s.InstanceIds),
			}
			result, err := (sess.Ec2).StopInstances(input)

			if err != nil {
				return nil, err
			}
			return result, nil

		}
		return nil, fmt.Errorf("You provided empty struct to StopInstances, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// CreateImage is responsible for capturing the image of the server who's Id is passed to it.
func (sess *EstablishedSession) CreateImage(img *ImageCreateInput) (*ec2.CreateImageOutput, error) {

	if sess.Ec2 != nil {
		if (img.ServerName != "") || (img.InstanceId != "") || (img.Description != "") {
			input := &ec2.CreateImageInput{
				Description: aws.String(img.Description),
				InstanceId:  aws.String(img.InstanceId),
				Name:        aws.String(img.ServerName),
			}
			result, err := (sess.Ec2).CreateImage(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf("You provided empty struct to CreateImage, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// DeregisterImage along with DeleteSnapshot has to be used to delete an image.
func (sess *EstablishedSession) DeregisterImage(img *DeleteComputeInput) error {

	if sess.Ec2 != nil {
		if img.ImageId != "" {
			// deregistering image will be done by below code
			input := &ec2.DeregisterImageInput{ImageId: aws.String(img.ImageId)}
			_, err := (sess.Ec2).DeregisterImage(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to DeregisterImage, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// DeleteSnapshot helps in deleting the snapshot, but this is more effective while deleting image from aws.
func (sess *EstablishedSession) DeleteSnapshot(img *DeleteComputeInput) error {

	if sess.Ec2 != nil {
		if img.SnapshotId != "" {
			// Deletion of snapshot will addressed by below code
			input := &ec2.DeleteSnapshotInput{SnapshotId: aws.String(img.SnapshotId)}
			_, err := (sess.Ec2).DeleteSnapshot(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to DeleteSnapshot, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// DescribeImages helps in fetching information about the images who's Id is passed. This is achieved by describing it.
func (sess *EstablishedSession) DescribeImages(img *DescribeComputeInput) (*ec2.DescribeImagesOutput, error) {

	if sess.Ec2 != nil {
		if img.ImageIds != nil {
			// desribing image to check if image exists
			searchImageInput := &ec2.DescribeImagesInput{
				ImageIds: aws.StringSlice(img.ImageIds),
			}
			result, err := (sess.Ec2).DescribeImages(searchImageInput)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf("You provided empty struct to DescribeImages, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// DescribeAllImages helps in describing all the images present in region, account or any depending on the filter you apply.
func (sess *EstablishedSession) DescribeAllImages(img *DescribeComputeInput) (*ec2.DescribeImagesOutput, error) {

	if sess.Ec2 != nil {
		// desribing image to check if image exists
		input := &ec2.DescribeImagesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{Name: aws.String("is-public"),
					Values: aws.StringSlice([]string{"false"}),
				},
			},
		}
		result, err := (sess.Ec2).DescribeImages(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// WaitTillInstanceAvailable makes the called method to wait till the created instance becomes available.
func (sess *EstablishedSession) WaitTillInstanceAvailable(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceRunning(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceAvailable, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// WaitTillInstanceRunning makes the called method to wait till the created/started instance enters to runnig state.
func (sess *EstablishedSession) WaitTillInstanceRunning(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceRunning(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceRunning, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// WaitTillInstanceStopped makes the called method to wait till the specified instance enters the stop state.
func (sess *EstablishedSession) WaitTillInstanceStopped(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceStopped(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceStopped, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// WaitTillInstanceTerminated makes the called method to wait till the specified instance is terminated successfully.
func (sess *EstablishedSession) WaitTillInstanceTerminated(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceTerminated(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceTerminated, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}
