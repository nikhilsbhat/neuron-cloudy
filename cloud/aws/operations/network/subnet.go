package awsnetwork

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	common "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/common"
)

// SubnetReponse is a struct that will be the response type of almost all the subnet related activities under cloud/operations.
type SubnetReponse struct {
	Name            string                     `json:"Name,omitempty"`
	Id              string                     `json:"Id,omitempty"`
	State           string                     `json:"State,omitempty"`
	VpcId           string                     `json:"VpcId,omitempty"`
	CreateSubnetRaw *ec2.CreateSubnetOutput    `json:"CreateSubnetRaw,omitempty"`
	GetSubnetRaw    *ec2.DescribeSubnetsOutput `json:"GetSubnetRaw,omitempty"`
}

// CreateSubnet is a customized method for subnet creation, if one needs plain subnet creation then he/she has to call the GOD, interface which talks to cloud.
func (subin *NetworkCreateInput) CreateSubnet(con neuronaws.EstablishConnectionInput) (SubnetReponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return SubnetReponse{}, seserr
	}
	// I am gathering inputs since create subnets needs it
	sub, suberr := ec2.CreateSubnet(
		&neuronaws.CreateNetworkInput{
			Cidr:  subin.SubCidr,
			VpcId: subin.VpcId,
			Zone:  subin.Zone,
		},
	)

	// handling the error if it throws while subnet is under creation process
	if suberr != nil {
		return SubnetReponse{}, suberr
	}

	// I will make program wait until subnet become available
	waiterr := ec2.WaitTillSubnetAvailable(
		&neuronaws.DescribeNetworkInput{
			Filters: neuronaws.Filters{
				Name:  "subnet-id",
				Value: []string{*sub.Subnet.SubnetId},
			},
		},
	)

	if waiterr != nil {
		return SubnetReponse{}, waiterr
	}

	// I will be the spock for tags creation.
	tags := common.Tag{*sub.Subnet.SubnetId, "Name", subin.Name}
	subtag, tagerr := tags.CreateTags(con)
	if tagerr != nil {
		return SubnetReponse{}, tagerr
	}

	routes := new(NetworkComponentInput)
	routes.Name = subin.Name + "_route"
	routes.VpcIds = []string{subin.VpcId}
	routes.SubId = *sub.Subnet.SubnetId
	routes.IgwId = subin.IgwId
	routes.SubType = subin.Type

	routeerr := routes.CreateRouteTable(con)

	if routeerr != nil {
		return SubnetReponse{}, routeerr
	}

	if subin.GetRaw == true {
		return SubnetReponse{CreateSubnetRaw: sub}, nil
	}

	return SubnetReponse{Name: subtag, Id: *sub.Subnet.SubnetId}, nil
}

// GetAllSubnets is a customized method for fetching details of all subnets for a given region, if one needs plain get subnet then he/she has to call the GOD, interface which talks to cloud.
func (net *GetNetworksInput) GetAllSubnets(con neuronaws.EstablishConnectionInput) (NetworkResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkResponse{}, seserr
	}

	result, err := ec2.DescribeAllSubnet(
		&neuronaws.DescribeNetworkInput{},
	)
	if err != nil {
		return NetworkResponse{}, err
	}

	if net.GetRaw == true {
		return NetworkResponse{GetSubnetRaw: result}, nil
	}

	subnets := make([]SubnetReponse, 0)
	for _, subnet := range result.Subnets {
		subnets = append(subnets, SubnetReponse{Name: *subnet.Tags[0].Value, Id: *subnet.SubnetId, State: *subnet.State, VpcId: *subnet.VpcId})
	}
	return NetworkResponse{Subnets: subnets}, nil

}

// GetSubnets is a customized method for fetching details of a particular subnet for a given region, if one needs plain get subnet then he/she has to call the GOD, interface which talks to cloud.
func (net *GetNetworksInput) GetSubnets(con neuronaws.EstablishConnectionInput) (NetworkResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkResponse{}, seserr
	}
	result, err := ec2.DescribeSubnet(
		&neuronaws.DescribeNetworkInput{
			SubnetIds: net.SubnetIds,
		},
	)
	if err != nil {
		return NetworkResponse{}, err
	}

	if net.GetRaw == true {
		return NetworkResponse{GetSubnetRaw: result}, nil
	}

	subnets := make([]SubnetReponse, 0)
	for _, subnet := range result.Subnets {
		if subnet.Tags[0] != nil {
			subnets = append(subnets, SubnetReponse{Name: *subnet.Tags[0].Value, Id: *subnet.SubnetId, State: *subnet.State, VpcId: *subnet.VpcId})
		} else {
			subnets = append(subnets, SubnetReponse{Id: *subnet.SubnetId, State: *subnet.State, VpcId: *subnet.VpcId})
		}
	}
	return NetworkResponse{Subnets: subnets}, nil

}

// GetSubnetsFromVpc is method which gets the list of available subnets from a asked network.
// Passing multiple values in vpcids array makes no difference here as we use only first element of it, this is customized function for raw data refer interface
func (net *GetNetworksInput) GetSubnetsFromVpc(con neuronaws.EstablishConnectionInput) (NetworkResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkResponse{}, seserr
	}
	result, err := ec2.DescribeSubnet(
		&neuronaws.DescribeNetworkInput{
			Filters: neuronaws.Filters{
				Name:  "vpc-id",
				Value: net.VpcIds,
			},
		},
	)
	if err != nil {
		return NetworkResponse{}, err
	}

	if net.GetRaw == true {
		return NetworkResponse{GetSubnetRaw: result}, nil
	}

	subnets := make([]SubnetReponse, 0)
	for _, subnet := range result.Subnets {
		if subnet.Tags != nil {
			subnets = append(subnets, SubnetReponse{Name: *subnet.Tags[0].Value, Id: *subnet.SubnetId, State: *subnet.State})
		} else {
			subnets = append(subnets, SubnetReponse{Id: *subnet.SubnetId, State: *subnet.State})
		}
	}
	return NetworkResponse{VpcId: net.VpcIds[0], Subnets: subnets}, nil

}

// DeleteSubnets is a customized method for deleting subnets, if one needs plain subnet deletion then he/she has to call the GOD, interface which talks to cloud.
func (s *DeleteNetworkInput) DeleteSubnets(con neuronaws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}

	for _, subnet := range s.SubnetIds {
		err := ec2.DeleteSubnet(
			&neuronaws.DescribeNetworkInput{
				SubnetIds: []string{subnet},
			},
		)
		if err != nil {
			return err
		}

		//Waiting till subnets deletion is successfully completed
		subwait, subwaiterr := ec2.WaitUntilSubnetDeleted(
			&neuronaws.DescribeNetworkInput{
				SubnetIds: []string{subnet},
			},
		)
		if subwaiterr != nil {
			return subwaiterr
		}
		if subwait == false {
			return fmt.Errorf("An error occurred while deleting a subnet")
		}
	}
	return nil
}

// FindSubnet is a customized method which sends back the response to the caller about the existence of subnet asked for.
func (net *GetNetworksInput) FindSubnet(con neuronaws.EstablishConnectionInput) (bool, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return false, seserr
	}
	result, err := ec2.DescribeSubnet(
		&neuronaws.DescribeNetworkInput{
			SubnetIds: net.SubnetIds,
		},
	)
	if err != nil {
		return false, err
	}
	if len(result.Subnets) != 0 {
		return true, nil
	}
	return false, nil
}

// GetVpcFromSubnet is a customized method which helps in fetching VPC from subnet asked for.
// Passing multi valued array make no difference as this is customized
func (net *GetNetworksInput) GetVpcFromSubnet(con neuronaws.EstablishConnectionInput) (SubnetReponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return SubnetReponse{}, seserr
	}
	result, err := ec2.DescribeSubnet(
		&neuronaws.DescribeNetworkInput{
			SubnetIds: net.SubnetIds,
		},
	)
	if err != nil {
		return SubnetReponse{}, err
	}
	return SubnetReponse{VpcId: *result.Subnets[0].VpcId}, nil
}
