package neuronaws

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	err "github.com/nikhilsbhat/neuron-cloudy/errors"
)

// Filters holds the required values for creating filter, so that this can be applied on resource to get the appropriate result.
type Filters struct {
	Name  string
	Value []string
}

// CreateNetworkInput implements the methods for creating the network. This also holds the values for the same.
type CreateNetworkInput struct {
	// CIDR block of the network that has to be created.
	Cidr    string
	Tenancy string
	// VpcId is the ID of the VPC/network created.
	VpcId string
	// SubId is the ID of the subnetwork which is present in the network created.
	SubId string
	// IgwId is the ID of internet gateway  created in a particular network.
	IgwId string
	// Zone of the particular resource created or to be created.
	Zone string
	// Name the resource that would be created.
	Name string
	// DestinationCidr is the CIDR used while writing routes to routetable.
	DestinationCidr string
	// RouteTableId is the ID of the routetable created or to be created.
	RouteTableId string
}

// IngressEgressInput holds the required values for creating ingress/egress rule for the specified security group and implements the methods for the same.
type IngressEgressInput struct {
	// Port number that has to opened as part of both ingress/egress rule.
	Port int64
	// SecId refers to ID of security group to which the rule has to be applied.
	SecId string
}

// DescribeNetworkInput holds all the required values for fetching the information about the selected network and its components.
// This is achieved by describing the same.
type DescribeNetworkInput struct {
	// SecIds are the security security id's created in a particular network.
	SecIds []string
	// IgwIds are the internet gateway id's created in a particular network.
	IgwIds []string
	// VpcIds are the id's of your virtual private cloud/network
	VpcIds []string
	// SubnetIds are the id's of the subnets created within the network.
	SubnetIds []string
	// RouteTableIds are the id's of the routetable that are created as part of network creation.
	RouteTableIds []string
	// Filters can be applied over the resource to get more precise data about it.
	Filters Filters
	// AssociationsId are the id's used to identify the RouteTable are are used while detaching RouteTable from the subnetwork.
	AssociationsId string
}

// CreateVpc will create the network/VPC in aws based on the values and session passed to it.
func (sess *EstablishedSession) CreateVpc(v *CreateNetworkInput) (*ec2.CreateVpcOutput, error) {

	if sess.Ec2 != nil {
		if (len(v.Cidr) != 0) || (len(v.Tenancy) != 0) {
			input := &ec2.CreateVpcInput{
				CidrBlock:       aws.String(v.Cidr),
				InstanceTenancy: aws.String(v.Tenancy),
			}
			result, err := (sess.Ec2).CreateVpc(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateVpc", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

// CreateSubnet will create the subnetwork in aws based on the values and session passed to it.
func (sess *EstablishedSession) CreateSubnet(s *CreateNetworkInput) (*ec2.CreateSubnetOutput, error) {

	if sess.Ec2 != nil {
		if (s.Cidr != "") || (s.VpcId != "") || (s.Zone != "") {
			input := &ec2.CreateSubnetInput{
				CidrBlock:        aws.String(s.Cidr),
				VpcId:            aws.String(s.VpcId),
				AvailabilityZone: aws.String(s.Zone),
			}
			result, err := (sess.Ec2).CreateSubnet(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateSubnet", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

// CreateIgw will create the internetgateway so that the compute components can talk to outside world and vise versa when attached to the network where the components is present.
func (sess *EstablishedSession) CreateIgw() (*ec2.CreateInternetGatewayOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.CreateInternetGatewayInput{}
		result, err := (sess.Ec2).CreateInternetGateway(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// AttachIgw will attach the internetgateway to the specified network.
//never pass arrays to attach the intergateways, it never works
func (sess *EstablishedSession) AttachIgw(a *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if (a.IgwIds != nil) || (a.VpcIds != nil) {
			input := &ec2.AttachInternetGatewayInput{
				InternetGatewayId: aws.String(a.IgwIds[0]),
				VpcId:             aws.String(a.VpcIds[0]),
			}
			_, err := (sess.Ec2).AttachInternetGateway(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v AttachIgw", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// CreateSecurityGroup creates and configure the security group with all the parameter passed to it.
func (sess *EstablishedSession) CreateSecurityGroup(s *CreateNetworkInput) (*ec2.CreateSecurityGroupOutput, error) {

	if sess.Ec2 != nil {
		if s.VpcId != "" {
			input := &ec2.CreateSecurityGroupInput{
				Description: aws.String("This security group is created by Neuron api"),
				VpcId:       aws.String(s.VpcId),
				GroupName:   aws.String(s.Name),
			}
			result, err := (sess.Ec2).CreateSecurityGroup(input)

			if err != nil {
				return nil, err
			}
			return result, nil

		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateSecurityGroup", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()

}

// CreateRouteTable creats the routetable, which is one of the network component.
func (sess *EstablishedSession) CreateRouteTable(r *CreateNetworkInput) (*ec2.CreateRouteTableOutput, error) {

	if sess.Ec2 != nil {
		if r.VpcId != "" {
			input := &ec2.CreateRouteTableInput{
				VpcId: aws.String(r.VpcId),
			}
			result, err := (sess.Ec2).CreateRouteTable(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v CreateRouteTable", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

// WriteRoute will write routes into the specified route table.
func (sess *EstablishedSession) WriteRoute(r *CreateNetworkInput) error {

	if sess.Ec2 != nil {
		if r.RouteTableId != "" {
			input := &ec2.CreateRouteInput{
				DestinationCidrBlock: aws.String(r.DestinationCidr),
				GatewayId:            aws.String(r.IgwId),
				RouteTableId:         aws.String(r.RouteTableId),
			}
			_, err := (sess.Ec2).CreateRoute(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v WriteRoute", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

// AttachRouteTable attaches the route table to the specified subnetwork.
func (sess *EstablishedSession) AttachRouteTable(r *CreateNetworkInput) error {

	if sess.Ec2 != nil {
		if r.RouteTableId != "" {
			input := &ec2.AssociateRouteTableInput{
				RouteTableId: aws.String(r.RouteTableId),
				SubnetId:     aws.String(r.SubId),
			}
			_, err := (sess.Ec2).AssociateRouteTable(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v AttachRouteTable", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// DettachRouteTable dettaches the route table from the subnetwork it is attached to.
func (sess *EstablishedSession) DettachRouteTable(r *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if r.AssociationsId != "" {
			input := &ec2.DisassociateRouteTableInput{
				AssociationId: aws.String(r.AssociationsId),
			}
			_, err := (sess.Ec2).DisassociateRouteTable(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v AttachRouteTable", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// CreateEgressRule creates the egress rule with the specified configuration.
func (sess *EstablishedSession) CreateEgressRule(i *IngressEgressInput) error {

	if sess.Ec2 != nil {
		securityIngressInput := &ec2.AuthorizeSecurityGroupEgressInput{
			//    FromPort   : aws.int64(from_port),
			GroupId: aws.String(i.SecId),
			//    ToPort     : aws.int64(to_port),
			IpPermissions: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("-1"),
					IpRanges: []*ec2.IpRange{
						{
							//CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		}
		_, egressErr := (sess.Ec2).AuthorizeSecurityGroupEgress(securityIngressInput)

		if egressErr != nil {
			return egressErr
		}
		return nil
	}
	return err.InvalidSession()

}

// CreateIngressRule creates the ingress rule with the specified configurations.
func (sess *EstablishedSession) CreateIngressRule(i *IngressEgressInput) error {

	if sess.Ec2 != nil {
		securityIngressInput := &ec2.AuthorizeSecurityGroupIngressInput{
			FromPort:   aws.Int64(i.Port),
			IpProtocol: aws.String("tcp"),
			GroupId:    aws.String(i.SecId),
			ToPort:     aws.Int64(i.Port),
			CidrIp:     aws.String("0.0.0.0/0"),
		}
		_, ingressErr := (sess.Ec2).AuthorizeSecurityGroupIngress(securityIngressInput)

		if ingressErr != nil {
			return ingressErr
		}
		return nil
	}
	return err.InvalidSession()

}

// DeleteIgw deletes the specified internetgateway selected.
func (sess *EstablishedSession) DeleteIgw(i *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if i.IgwIds != nil {

			input := &ec2.DeleteInternetGatewayInput{
				InternetGatewayId: aws.String(i.IgwIds[0]),
			}
			_, err := (sess.Ec2).DeleteInternetGateway(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteIgw", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

// DetachIgw helps in dettaching the internetgateway from the attached network.
//never pass arrays to dettach the intergateways, it never works
func (sess *EstablishedSession) DetachIgw(i *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if i.IgwIds != nil {
			input := &ec2.DetachInternetGatewayInput{
				InternetGatewayId: aws.String(i.IgwIds[0]),
				VpcId:             aws.String(i.VpcIds[0]),
			}
			_, err := (sess.Ec2).DetachInternetGateway(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DetachIgw", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

// DescribeIgw fetches the information about the intertnetgateway selected.
// This is achieved by describing the same.
func (sess *EstablishedSession) DescribeIgw(d *DescribeNetworkInput) (*ec2.DescribeInternetGatewaysOutput, error) {

	if sess.Ec2 != nil {
		if d.IgwIds != nil {
			input := &ec2.DescribeInternetGatewaysInput{
				InternetGatewayIds: aws.StringSlice(d.IgwIds),
			}
			result, err := (sess.Ec2).DescribeInternetGateways(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeIgw. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeInternetGatewaysInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeInternetGateways(input)

		if err != nil {
			return nil, err
		}
		return result, nil

	}
	return nil, err.InvalidSession()
}

// DescribeAllIgw describes all the internetgateway in the specified region of the cloud.
func (sess *EstablishedSession) DescribeAllIgw(d *DescribeNetworkInput) (*ec2.DescribeInternetGatewaysOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeInternetGatewaysInput{}
		result, err := (sess.Ec2).DescribeInternetGateways(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DeleteRouteTable deletes the specified route table.
func (sess *EstablishedSession) DeleteRouteTable(r *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if r.RouteTableIds != nil {
			input := &ec2.DeleteRouteTableInput{
				RouteTableId: aws.String(r.RouteTableIds[0]),
			}
			_, err := (sess.Ec2).DeleteRouteTable(input)
			if err != nil {
				return err
			}
			return nil

		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteRouteTable", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// DescribeRouteTable fetches the information about the selected route table by describing it.
func (sess *EstablishedSession) DescribeRouteTable(d *DescribeNetworkInput) (*ec2.DescribeRouteTablesOutput, error) {

	if sess.Ec2 != nil {
		if d.RouteTableIds != nil {
			input := &ec2.DescribeRouteTablesInput{
				RouteTableIds: aws.StringSlice(d.RouteTableIds),
			}
			result, err := (sess.Ec2).DescribeRouteTables(input)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeRouteTable. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeRouteTables(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DeleteSecurityGroup deletes the specified security group.
func (sess *EstablishedSession) DeleteSecurityGroup(s *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if s.SecIds != nil {
			for _, sec := range s.SecIds {
				input := &ec2.DeleteSecurityGroupInput{
					GroupId: aws.String(sec),
				}
				_, err := (sess.Ec2).DeleteSecurityGroup(input)
				if err != nil {
					return err
				}
				return nil
			}
		}

		return fmt.Errorf(fmt.Sprintf("%v DeleteSecurityGroup", err.EmptyStructError()))
	}
	return err.InvalidSession()

}

// DescribeSecurityGroup fetches the information about the selected security group. This is achieved by describing it.
func (sess *EstablishedSession) DescribeSecurityGroup(d *DescribeNetworkInput) (*ec2.DescribeSecurityGroupsOutput, error) {

	if sess.Ec2 != nil {
		if d.SecIds != nil {
			input := &ec2.DescribeSecurityGroupsInput{
				GroupIds: aws.StringSlice(d.SecIds),
			}
			result, err := (sess.Ec2).DescribeSecurityGroups(input)
			if err != nil {
				return nil, err
			}
			return result, nil

		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeSecurityGroup. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeSecurityGroupsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeSecurityGroups(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

// DescribeAllSecurityGroup fetches the information about all the security group present in the specified region.
func (sess *EstablishedSession) DescribeAllSecurityGroup(d *DescribeNetworkInput) (*ec2.DescribeSecurityGroupsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeSecurityGroupsInput{}
		result, err := (sess.Ec2).DescribeSecurityGroups(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DescribeAllSubnet fetches the information about all the subnetnetworks present in the specified region.
func (sess *EstablishedSession) DescribeAllSubnet(d *DescribeNetworkInput) (*ec2.DescribeSubnetsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeSubnetsInput{}
		result, err := (sess.Ec2).DescribeSubnets(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DescribeSubnet fetches the information about the selected subnetwork. This is achieved by describing it.
func (sess *EstablishedSession) DescribeSubnet(d *DescribeNetworkInput) (*ec2.DescribeSubnetsOutput, error) {

	if sess.Ec2 != nil {
		if d.SubnetIds != nil {
			input := &ec2.DescribeSubnetsInput{
				SubnetIds: aws.StringSlice(d.SubnetIds),
			}
			result, err := (sess.Ec2).DescribeSubnets(input)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeSubnet. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeSubnets(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

// DeleteSubnet helps in deletion of the selected subnetwork.
func (sess *EstablishedSession) DeleteSubnet(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if d.SubnetIds != nil {
			input := &ec2.DeleteSubnetInput{
				SubnetId: aws.String(d.SubnetIds[0]),
			}
			_, err := (sess.Ec2).DeleteSubnet(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteSubnet", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// DescribeAllVpc fetches the required information of all network present in the selected region. This is achieved by describing it.
func (sess *EstablishedSession) DescribeAllVpc(d *DescribeNetworkInput) (*ec2.DescribeVpcsOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeVpcsInput{}
		result, err := (sess.Ec2).DescribeVpcs(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DescribeVpc fetch the information about the specified network.
func (sess *EstablishedSession) DescribeVpc(d *DescribeNetworkInput) (*ec2.DescribeVpcsOutput, error) {

	if sess.Ec2 != nil {
		if d.VpcIds != nil {
			input := &ec2.DescribeVpcsInput{
				VpcIds: aws.StringSlice(d.VpcIds),
			}
			result, err := (sess.Ec2).DescribeVpcs(input)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return nil, fmt.Errorf(fmt.Sprintf("%v DescribeVpc. You selected filters for this yet you passed empty", err.EmptyStructError()))
		}
		input := &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeVpcs(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DeleteVpc deletes the specified network in clopud aws.
func (sess *EstablishedSession) DeleteVpc(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if d.VpcIds != nil {
			for _, vpc := range d.VpcIds {
				input := &ec2.DeleteVpcInput{
					VpcId: aws.String(vpc),
				}
				_, err := (sess.Ec2).DeleteVpc(input)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteVpc", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// WaitTillVpcAvailable makes the method called this to wait until the network becomes available.
func (sess *EstablishedSession) WaitTillVpcAvailable(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {

		if reflect.DeepEqual(d.Filters, Filters{}) {
			return fmt.Errorf(fmt.Sprintf("%v WaitTillVpcAvailable", err.EmptyStructError()))
		}
		input := &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		err := (sess.Ec2).WaitUntilVpcAvailable(input)
		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}

// WaitTillSubnetAvailable makes the method called this to wait till the subnet state becomes available.
func (sess *EstablishedSession) WaitTillSubnetAvailable(d *DescribeNetworkInput) error {

	if sess.Ec2 != nil {
		if reflect.DeepEqual(d.Filters, Filters{}) {
			return fmt.Errorf(fmt.Sprintf("%v WaitTillSubnetAvailable", err.EmptyStructError()))
		}
		input := &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String(d.Filters.Name),
					Values: aws.StringSlice(d.Filters.Value),
				},
			},
		}
		err := (sess.Ec2).WaitUntilSubnetAvailable(input)
		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}

// WaitUntilSubnetDeleted makes the method called this to wait till the subnet is successfully deleted.
func (sess *EstablishedSession) WaitUntilSubnetDeleted(d *DescribeNetworkInput) (bool, error) {

	if sess.Ec2 != nil {
		if d.SubnetIds != nil {
			input := &ec2.DescribeSubnetsInput{
				SubnetIds: aws.StringSlice(d.SubnetIds),
			}

			response, deserr := (sess.Ec2).DescribeSubnets(input)
			if response.Subnets != nil {
				start := time.Now()
				for len(response.Subnets) > 0 {
					response, deserr = (sess.Ec2).DescribeSubnets(input)
					if deserr != nil {
						switch deserr.(awserr.Error).Code() {
						case "InvalidSubnetID.NotFound":
							return true, nil
						default:
							return false, deserr
						}
					}
					if time.Since(start) > time.Duration(10*time.Second) {
						return false, fmt.Errorf("Time Out .Oops...!! it took annoyingly more than anticipated time while waiting for subnet to get deleted. Guess I was not called after delete subnet function")
					}
				}
			}
			if deserr != nil {
				switch deserr.(awserr.Error).Code() {
				case "InvalidSubnetID.NotFound":
					return true, nil
				default:
					return false, deserr
				}
			}
			return false, fmt.Errorf("Error occurred while waiting for the subnet deletion")
		}
		return false, fmt.Errorf(fmt.Sprintf("%v WaitUntilSubnetDeleted", err.EmptyStructError()))
	}
	return false, err.InvalidSession()
}

// WaitUntilRoutTableDeleted makes the method called this to wait till the route table is successfully deleted.
func (sess *EstablishedSession) WaitUntilRoutTableDeleted(d *DescribeNetworkInput) (bool, error) {

	if sess.Ec2 != nil {
		if d.RouteTableIds != nil {
			input := &ec2.DescribeRouteTablesInput{
				RouteTableIds: aws.StringSlice(d.RouteTableIds),
			}

			response, deserr := (sess.Ec2).DescribeRouteTables(input)
			if response.RouteTables != nil {
				start := time.Now()
				for len(response.RouteTables) > 0 {
					response, deserr = (sess.Ec2).DescribeRouteTables(input)
					if deserr != nil {
						switch deserr.(awserr.Error).Code() {
						case "InvalidRouteTableID.NotFound":
							return true, nil
						default:
							return false, deserr
						}
					}
					if time.Since(start) > time.Duration(10*time.Second) {
						return false, fmt.Errorf("Time Out .Oops...!! it took annoyingly more than anticipated time while waiting for routetable to get deleted. Guess I was not called after delete routetable function")
					}
				}
			}
			if deserr != nil {
				switch deserr.(awserr.Error).Code() {
				case "InvalidRouteTableID.NotFound":
					return true, nil
				default:
					return false, deserr
				}
			}
			return false, fmt.Errorf("Error occurred while waiting for the routetable deletion")
		}
		return false, fmt.Errorf(fmt.Sprintf("%v WaitUntilRoutTableDeleted", err.EmptyStructError()))
	}
	return false, err.InvalidSession()
}

// WaitUntilIgwDeleted makes the method called this to wait until internetgateway deletion is successful.
func (sess *EstablishedSession) WaitUntilIgwDeleted(d *DescribeNetworkInput) (bool, error) {

	if sess.Ec2 != nil {
		if d.IgwIds != nil {
			input := &ec2.DescribeInternetGatewaysInput{
				InternetGatewayIds: aws.StringSlice(d.IgwIds),
			}

			response, deserr := (sess.Ec2).DescribeInternetGateways(input)
			if response.InternetGateways != nil {
				start := time.Now()
				for len(response.InternetGateways) > 0 {
					response, deserr = (sess.Ec2).DescribeInternetGateways(input)
					if deserr != nil {
						switch deserr.(awserr.Error).Code() {
						case "InvalidInternetGatewayID.NotFound":
							return true, nil
						default:
							return false, deserr
						}
					}
					if time.Since(start) > time.Duration(10*time.Second) {
						return false, fmt.Errorf("Time Out .Oops...!! it took annoyingly more than anticipated time while waiting for igw to get deleted. Guess I was not called after delete igw function")
					}
				}
			}
			if deserr != nil {
				switch deserr.(awserr.Error).Code() {
				case "InvalidInternetGatewayID.NotFound":
					return true, nil
				default:
					return false, deserr
				}
			}
			return false, fmt.Errorf("Error occurred while waiting for the InternetGateways deletion")
		}
		return false, fmt.Errorf(fmt.Sprintf("%v WaitUntilIgwDeleted", err.EmptyStructError()))
	}
	return false, err.InvalidSession()
}
