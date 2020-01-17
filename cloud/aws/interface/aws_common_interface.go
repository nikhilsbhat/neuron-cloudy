package neuronaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AwsCommonInput holds the common input values to perform the operations.
type AwsCommonInput struct {
	// AvailabilityZone holds the name of the zone who's information has to be retrieved.
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
}

// CreateTagsInput holds the input values which helps in creating tags for the resource created in aws.
type CreateTagsInput struct {
	// Resource name to which requries the naming.
	Resource string
	// Name or the Tag name that has to be assigned to the resource that was created.
	Name string
	// Value for the tag that has to be created.
	Value string
}

// DescribeAllAvailabilityZones describes all the availability zones present in the aws.
func (sess *EstablishedSession) DescribeAllAvailabilityZones(a *AwsCommonInput) (*ec2.DescribeAvailabilityZonesOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeAvailabilityZonesInput{}
		result, err := (sess.Ec2).DescribeAvailabilityZones(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

// DescribeAvailabilityZones describes the particular availability zoa=ne selected.
func (sess *EstablishedSession) DescribeAvailabilityZones(a *AwsCommonInput) (*ec2.DescribeAvailabilityZonesOutput, error) {

	if sess.Ec2 != nil {
		if a.AvailabilityZone != "" {
			input := &ec2.DescribeAvailabilityZonesInput{
				ZoneNames: aws.StringSlice([]string{a.AvailabilityZone}),
			}
			result, err := (sess.Ec2).DescribeAvailabilityZones(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

// CreateTags creates tag for the resources created. It works only with the help of resource Id, make sure you pass right reource Id.
func (sess *EstablishedSession) CreateTags(t *CreateTagsInput) error {

	if sess.Ec2 != nil {
		input := &ec2.CreateTagsInput{
			Resources: []*string{
				aws.String(t.Resource),
			},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(t.Name),
					Value: aws.String(t.Value),
				},
			},
		}
		_, err := (sess.Ec2).CreateTags(input)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// GetRegions describes the regions available in the cloud aws.
func (sess *EstablishedSession) GetRegions() (*ec2.DescribeRegionsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeRegionsInput{}
		result, err := (sess.Ec2).DescribeRegions(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}
