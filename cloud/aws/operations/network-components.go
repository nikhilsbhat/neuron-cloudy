package aws

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// NetworkComponentInput  will implement the methods which deals with the creation/deletion of network components under cloud/operations.
type NetworkComponentInput struct {
	Name            string   `json:"name"`
	VpcIds          []string `json:"vpcid"`
	SubId           string   `json:"subid"`
	IgwId           string   `json:"igwid"`
	IgwIds          []string `json:"igwids"`
	SubType         string   `json:"subtype"`
	Ports           []string `json:"ports"`
	Filters         Filters  `json:"filters"`
	SecGroupIds     []string `json:"secgroupids"`
	RouteTableIds   []string `json:"routetableids"`
	DestinationCidr string   `json:"destinationcidr"`
	GetRaw          bool     `json:"getraw"`
}

// NetworkComponentResponse will be the response type of almost all the network components related activities under cloud/operations.
type NetworkComponentResponse struct {
	IgwIds            []string                            `json:"IgwId,omitempty"`
	SecGroupIds       []string                            `json:"SecGroupIds,omitempty"`
	RouteTableIds     []string                            `json:"RouteTableIds,omitempty"`
	CreateIgwRaw      *ec2.CreateInternetGatewayOutput    `json:"CreateIgwRaw,omitempty"`
	GetIgwRaw         *ec2.DescribeInternetGatewaysOutput `json:"GetIgwRaw,omitempty"`
	CreateSecurityRaw *ec2.CreateSecurityGroupOutput      `json:"CreateSecRaw,omitempty"`
	GetRouteTableRaw  *ec2.DescribeRouteTablesOutput      `json:"DescribeRouteTableRaw,omitempty"`
	GetSecurityRaw    *ec2.DescribeSecurityGroupsOutput   `json:"DescribeSecurityRaw,omitempty"`
}

// CreateIgw is customized internet-gateway creation, if one needs plain internet-gateway creation he/she has call interface the GOD which talks to cloud.
func (net *NetworkComponentInput) CreateIgw(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	ig, igerr := ec2.CreateIgw()
	if igerr != nil {
		return NetworkComponentResponse{}, igerr
	}

	if net.VpcIds != nil {
		aterr := ec2.AttachIgw(
			&aws.DescribeNetworkInput{
				IgwIds: []string{*ig.InternetGateway.InternetGatewayId},
				VpcIds: net.VpcIds,
			},
		)
		if aterr != nil {
			return NetworkComponentResponse{}, aterr
		}
	}

	igtags := new(Tag)
	igtags.Resource = *ig.InternetGateway.InternetGatewayId
	igtags.Name = "Name"
	igtags.Value = net.Name + "_igw"
	_, igtagerr := igtags.CreateTags(con)
	if igtagerr != nil {
		return NetworkComponentResponse{}, igtagerr
	}

	if net.GetRaw == true {
		return NetworkComponentResponse{CreateIgwRaw: ig}, nil
	}
	return NetworkComponentResponse{IgwIds: []string{*ig.InternetGateway.InternetGatewayId}}, nil
}

// GetIgwFromVpc will help one in fetching IGW from the VPC which they specify.
func (net *NetworkComponentInput) GetIgwFromVpc(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	response, err := ec2.DescribeIgw(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "attachment.vpc-id",
				Value: net.VpcIds,
			},
		},
	)
	if err != nil {
		return NetworkComponentResponse{}, err
	}

	if net.GetRaw == true {
		return NetworkComponentResponse{GetIgwRaw: response}, nil
	}

	igwids := make([]string, 0)
	for _, igw := range response.InternetGateways {
		igwids = append(igwids, *igw.InternetGatewayId)
	}

	return NetworkComponentResponse{IgwIds: igwids}, nil
}

// CreateSecurityGroup is customized security-group creation, if one needs plain security-group creation he/she has call interface the GOD which talks to cloud.
func (net *NetworkComponentInput) CreateSecurityGroup(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	security, secerr := ec2.CreateSecurityGroup(
		&aws.CreateNetworkInput{
			VpcId: net.VpcIds[0],
			Name:  net.Name + "_sec",
		},
	)
	if secerr != nil {
		return NetworkComponentResponse{}, secerr
	}

	sctags := new(Tag)
	sctags.Resource = *security.GroupId
	sctags.Name = "Name"
	sctags.Value = net.Name + "_sec"
	_, sctagerr := sctags.CreateTags(con)
	if sctagerr != nil {
		return NetworkComponentResponse{}, sctagerr
	}

	//creating egree and ingres rules for the security group which I created just now
	for _, port := range net.Ports {
		intport, _ := strconv.ParseInt(port, 10, 64)
		ingreserr := ec2.CreateIngressRule(
			&aws.IngressEgressInput{
				Port:  intport,
				SecId: *security.GroupId,
			},
		)
		if ingreserr != nil {
			return NetworkComponentResponse{}, ingreserr
		}
	}
	egreserr := ec2.CreateEgressRule(
		&aws.IngressEgressInput{
			SecId: *security.GroupId,
		},
	)
	if egreserr != nil {
		return NetworkComponentResponse{}, egreserr
	}

	if net.GetRaw == true {
		return NetworkComponentResponse{CreateSecurityRaw: security}, nil
	}
	return NetworkComponentResponse{SecGroupIds: []string{*security.GroupId}}, nil
}

// GetSecFromVpc will help one in fetching security-group from the VPC which they specify.
func (net *NetworkComponentInput) GetSecFromVpc(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	response, err := ec2.DescribeSecurityGroup(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: net.VpcIds,
			},
		},
	)
	if err != nil {
		return NetworkComponentResponse{}, err
	}
	secids := make([]string, 0)
	if net.GetRaw == true {
		return NetworkComponentResponse{GetSecurityRaw: response}, nil
	}
	for _, sec := range response.SecurityGroups {
		secids = append(secids, *sec.GroupId)
	}
	return NetworkComponentResponse{SecGroupIds: secids}, nil
}

// DeleteSecutiryGroup will help one in deleting security-group which they specify, for deletion of other resources refer other methods.
func (net *NetworkComponentInput) DeleteSecutiryGroup(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}
	err := ec2.DeleteSecurityGroup(
		&aws.DescribeNetworkInput{
			SecIds: net.SecGroupIds,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// CreateRouteTable is customized route-table creation, if one needs plain route-table creation he/she has call interface the GOD which talks to cloud.
func (net *NetworkComponentInput) CreateRouteTable(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}
	routetable, routetableerr := ec2.CreateRouteTable(
		&aws.CreateNetworkInput{
			VpcId: net.VpcIds[0],
		},
	)

	if routetableerr != nil {
		return routetableerr
	}

	if net.IgwId != "" {
		if strings.ToLower(net.SubType) == "public" {
			routeerr := ec2.WriteRoute(
				&aws.CreateNetworkInput{
					DestinationCidr: "0.0.0.0/0",
					IgwId:           net.IgwId,
					RouteTableId:    *routetable.RouteTable.RouteTableId,
				},
			)
			if routeerr != nil {
				return routeerr
			}

			routeattacherr := ec2.AttachRouteTable(
				&aws.CreateNetworkInput{
					RouteTableId: *routetable.RouteTable.RouteTableId,
					SubId:        net.SubId,
				},
			)
			if routeattacherr != nil {
				return routeattacherr
			}

			return nil
		} else {
			// Releasing Soon...!!. We are not supporting writing custom routes into route tables as of now.
			return nil
		}
	} else {
		if strings.ToLower(net.SubType) == "public" {
			igws, igwerr := ec2.DescribeAllIgw(
				&aws.DescribeNetworkInput{},
			)
			if igwerr != nil {
				return igwerr
			}
			for _, igw := range igws.InternetGateways {
				if *igw.Attachments[0].VpcId == net.VpcIds[0] {
					routeerr := ec2.WriteRoute(
						&aws.CreateNetworkInput{
							DestinationCidr: "0.0.0.0/0",
							IgwId:           *igw.InternetGatewayId,
							RouteTableId:    *routetable.RouteTable.RouteTableId,
						},
					)
					if routeerr != nil {
						return routeerr
					}

					routeattacherr := ec2.AttachRouteTable(
						&aws.CreateNetworkInput{
							RouteTableId: *routetable.RouteTable.RouteTableId,
							SubId:        net.SubId,
						},
					)
					if routeattacherr != nil {
						return routeattacherr
					}

					return nil
				}
			}

		} else {
			// Releasing Soon...!!. We are not supporting writing custom routes into route tables as of now.
			return nil
		}
		return nil
	}
}

// DisassociateRouteTable will help one in disassociating the route-table which you specify, from the subnet to which it is attached, for disassociating other resources refer other methods.
func (net *NetworkComponentInput) DisassociateRouteTable(con aws.EstablishConnectionInput) (bool, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return false, seserr
	}

	response, reserr := ec2.DescribeRouteTable(
		&aws.DescribeNetworkInput{
			RouteTableIds: net.RouteTableIds,
		},
	)
	if reserr != nil {
		return false, reserr
	}

	associationId := make([]string, 0)
	for _, routetable := range response.RouteTables {
		if routetable.Associations != nil {
			associationId = append(associationId, *routetable.Associations[0].RouteTableAssociationId)
		}
	}

	for _, id := range associationId {
		deterr := ec2.DettachRouteTable(
			&aws.DescribeNetworkInput{
				AssociationsId: id,
			},
		)
		if deterr != nil {
			return false, deterr
		}
	}
	return true, nil
}

// GetRouteTableFromVpc will help one in fetching route-tables from the VPC which they specify.
func (net *NetworkComponentInput) GetRouteTableFromVpc(con aws.EstablishConnectionInput) (NetworkComponentResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkComponentResponse{}, seserr
	}
	response, err := ec2.DescribeRouteTable(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: net.VpcIds,
			},
		},
	)
	if err != nil {
		return NetworkComponentResponse{}, err
	}

	routeids := make([]string, 0)

	if net.GetRaw == true {
		return NetworkComponentResponse{GetRouteTableRaw: response}, nil
	}

	for _, route := range response.RouteTables {
		routeids = append(routeids, *route.RouteTableId)
	}
	return NetworkComponentResponse{RouteTableIds: routeids}, nil
}

// DeleteRouteTable will help one in deleting route-table which they specify, for deletion of other resources refer other methods.
func (net *NetworkComponentInput) DeleteRouteTable(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}

	for _, route := range net.RouteTableIds {
		err := ec2.DeleteRouteTable(
			&aws.DescribeNetworkInput{
				RouteTableIds: []string{route},
			},
		)
		if err != nil {
			return err
		}

		//Waiting till reoutetables deletion is successfully completed
		routewait, routwaiterr := ec2.WaitUntilRoutTableDeleted(
			&aws.DescribeNetworkInput{
				RouteTableIds: []string{route},
			},
		)
		if routwaiterr != nil {
			return routwaiterr
		}
		if routewait == false {
			return fmt.Errorf("An error occurred while deleting a routetable")
		}
	}
	return nil
}

// DetachIgws will help one in dettaching the IGW's which you specify, from the network to which it is attached, for dettaching other resources refer other methods.
func (net *NetworkComponentInput) DetachIgws(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}
	for _, igw := range net.IgwIds {
		err := ec2.DetachIgw(
			&aws.DescribeNetworkInput{
				IgwIds: []string{igw},
				VpcIds: net.VpcIds,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteIgws will help one in deleting IGW's which they specify, for deletion of other resources refer other methods.
func (net *NetworkComponentInput) DeleteIgws(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}

	for _, igw := range net.IgwIds {
		err := ec2.DeleteIgw(
			&aws.DescribeNetworkInput{
				IgwIds: []string{igw},
				VpcIds: net.VpcIds,
			},
		)
		if err != nil {
			return err
		}

		//Waiting till internetgateways deletion is successfully completed
		igwwait, igwwaiterr := ec2.WaitUntilIgwDeleted(
			&aws.DescribeNetworkInput{
				IgwIds: []string{igw},
			},
		)
		if igwwaiterr != nil {
			return igwwaiterr
		}
		if igwwait == false {
			return fmt.Errorf("An error occurred while deleting a igws")
		}
	}
	return nil
}
