package aws

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// NetworkCreateInput will implement almost all the creation of network and its components under cloud/operations.
type NetworkCreateInput struct {
	// VpcCidr would be the CIDR bolck for the network that would be created.
	VpcCidr string `json:"vpccidr"`
	// SubCidrs would be the list of CIDR bolcks for the subnetworks that would be created.
	SubCidrs []string `json:"subcidrs"`
	// SubCidr would be the CIDR bolck for the subnetwork that would be created.
	SubCidr string `json:"subcidr"`
	// Name of the network that would be created.
	Name string `json:"name"`
	// Type of network to be created ex: public, private.
	Type string `json:"type"`
	// Ports to be opened on the network that would be created.
	Ports []string `json:"ports"`
	// Zone name in which the network has to reside.
	Zone string `json:"zone"`
	// VpcId refers to the ID of the network in which the subnet should be created.
	VpcId string `json:"vpcid"`
	// IgwId refers to the ID of the internet gateway which should be updated/deleted.
	IgwId  string `json:"igwid"`
	GetRaw bool   `json:"getraw"`
}

// NetworkResponse will be the response type of almost all the network related activities under cloud/operations.
type NetworkResponse struct {
	// Name of the network/subnetwork that would be creted/updated/deleted.
	Name string `json:"name,omitempty"`
	// VpcId refers to the ID of network/subnetwork that would be creted/updated/deleted.
	VpcId string `json:"vpcid,omitempty"`
	// Subnets holds the response which was obtained as part of subnetwork cretion/updation/deletion.
	Subnets []SubnetReponse `json:"subnets,omitempty"`
	// Vpcs holds the response which was obtained as part of network cretion/updation/deletion.
	Vpcs []VpcResponse `json:"vpcs,omitempty"`
	// Type of the network/subnetwork that has to be creted/updated/deleted.
	Type string `json:"type,omitempty"`
	// State refers to current state of the network would be creted/updated/deleted.
	State string `json:"state,omitempty"`
	// IgwId refers to the ID of internet gateway that would be creted/updated/deleted.
	IgwId string `json:"igw,omitempty"`
	// IsDefault will define if the network/subnetwork or its components are pre created.
	IsDefault bool `json:"isdefault,omitempty"`
	// SecGroupIds are the list of security groups IDs that is associated with the network/subnetwork.
	SecGroupIds []string `json:"secgroupid,omitempty"`
	// Region name in which the network/subnetwork or its component were created.
	Region                string                              `json:"region,omitempty"`
	GetVpcsRaw            *ec2.DescribeVpcsOutput             `json:"getvpcsraw,omitempty"`
	GetVpcRaw             *ec2.Vpc                            `json:"getvpcraw,omitempty"`
	GetSubnetRaw          *ec2.DescribeSubnetsOutput          `json:"getsubnetraw,omitempty"`
	CreateVpcRaw          VpcResponse                         `json:"createvpcraw,omitempty"`
	CreateSubnetRaw       []SubnetReponse                     `json:"createsubnetraw,omitempty"`
	CreateIgwRaw          *ec2.CreateInternetGatewayOutput    `json:"createigwraw,omitempty"`
	CreateSecRaw          *ec2.CreateSecurityGroupOutput      `json:"createsecraw,omitempty"`
	DescribeRouteTableRaw *ec2.DescribeRouteTablesOutput      `json:"describeroutetableraw,omitempty"`
	DescribeSecurityRaw   *ec2.DescribeSecurityGroupsOutput   `json:"describesecurityraw,omitempty"`
	DescribeIgwRaw        *ec2.DescribeInternetGatewaysOutput `json:"describeigwraw,omitempty"`
}

// DeleteNetworkInput will implement almost all the deletion of network and its components under cloud/operations.
type DeleteNetworkInput struct {
	// VpcIds are the list of network IDs that would be deleted.
	VpcIds []string `json:"region"`
	// SubnetIds are the list of subnetwork IDs that would be deleted.
	SubnetIds []string `json:"vpcids"`
	// SecIds are the list of security group IDs that would be deleted.
	SecIds []string `json:"secids"`
	// IgwIds are the list of internet gateway IDs that would be deleted.
	IgwIds []string `json:"igwid"`
	// RouteTableIds are the list of routetable that are associated with subnetwork and has to be deleted.
	RouteTableIds []string `json:"routetableids"`
	GetRaw        bool     `json:"getraw"`
}

// GetNetworksInput will implement almost all the methods of fetching network and its components under cloud/operations.
type GetNetworksInput struct {
	// VpcIds are the list of network IDs of which the information has to be retrieved.
	VpcIds []string `json:"vpcids"`
	// SubnetIds are the list of subnetwork IDs of which the information has to be retrieved.
	SubnetIds []string `json:"subnetids"`
	// Filters could be applied in the resource to get more refined response.
	Filters Filters `json:"filters"`
	// Region name to which network/subnetwork belongs to.
	Region string `json:"region"`
	GetRaw bool   `json:"getraw"`
}

// DeleteNetworkResponse will be the response type of almost all the network and its components while terminating them.
type DeleteNetworkResponse struct {
	// Subnets holds the response on the status of subnetwork deleted.
	Subnets string `json:"subnets,omitempty"`
	// SecurityGroups holds the response on the status of security groups deleted.
	SecurityGroups string `json:"securities,omitempty"`
	// Routetables holds the response on the status of routetables deleted.
	Routetables string `json:"routetables,omitempty"`
	// Gateways holds the response on the status of internet gateways deleted.
	Gateways string `json:"gateways,omitempty"`
	// Vpcs holds the response on the status of network deleted.
	Vpcs            string `json:"vpcs,omitempty"`
	DefaultResponse string `json:"defaultresponse,omitempty"`
	Status          string `json:"status,omitempty"`
}

// UpdateNetworkInput will implement the methods that will update the network and its components under cloud/operations.
type UpdateNetworkInput struct {
	// Resource name which has to be updated.
	Resource string `json:"resource"`
	// Network collects the input for creation of network.
	Network NetworkCreateInput `json:"network"`
	// Action to be performed on the resource selected.
	Action string `json:"action"`
	GetRaw bool   `json:"getRaw"`
}

// Filters will help one to have a hold on the call that they make, will help to filter the quiries.
// Ex: use this while fetching list of subnet from appropriate vpc, by using vpc as a 'Value' of filter.
// type Filters struct {
// 	Name  string
// 	Value interface{}
// }

// CreateNetwork is a customized method for network creation, if one needs to create the individual components of network then call the appropriate methods.
func (net *NetworkCreateInput) CreateNetwork(con aws.EstablishConnectionInput) (NetworkResponse, error) {

	if (net.VpcCidr == "") || (net.Name == "") {
		return NetworkResponse{}, fmt.Errorf("You have not provided either CIDR or name for VPC, cannot proceed further")
	}

	/*get the relative sessions before proceeding further
	  ec2, seserr := con.EstablishConnection()
	  if seserr != nil {
	          return NetworkResponse{}, seserr
	  }*/

	netin := new(NetworkCreateInput)
	netin.VpcCidr = net.VpcCidr
	netin.Name = net.Name
	netin.Type = net.Type
	netin.Ports = net.Ports
	netin.GetRaw = net.GetRaw

	vpc, err := netin.CreateVpc(con)
	if err != nil {
		return NetworkResponse{}, err
	}

	zonein := CommonInput{}
	zones, zonerr := zonein.GetAvailabilityZones(con)
	if zonerr != nil {
		return NetworkResponse{}, zonerr
	}

	// This takes care creation of required number of subnets.
	subnets := make([]SubnetReponse, 0)

	zonenum := len(zones) - 1
	for i, sub := range net.SubCidrs {

		if zonenum < 0 {
			zonenum = len(zones) - 1
		}

		// Creating subnet by calling appropriate object
		netin.SubCidr = sub
		netin.Name = net.Name + "_sub" + strconv.Itoa(i)
		netin.Zone = zones[zonenum]
		if net.GetRaw == true {
			netin.VpcId = *vpc.CreateVpcRaw.Vpc.VpcId
			netin.IgwId = *vpc.CreateIgwRaw.InternetGateway.InternetGatewayId
		} else {
			netin.VpcId = vpc.VpcId
			netin.IgwId = vpc.IgwId
		}

		subnet, suberr := netin.CreateSubnet(con)
		if suberr != nil {
			return NetworkResponse{}, suberr
		}
		subnets = append(subnets, subnet)

		zonenum--
	}
	if net.GetRaw == true {
		return NetworkResponse{CreateVpcRaw: vpc, CreateSubnetRaw: subnets}, nil
	}
	return NetworkResponse{Name: vpc.Name, VpcId: vpc.VpcId, Subnets: subnets, Type: vpc.Type, IgwId: vpc.IgwId, SecGroupIds: vpc.SecGroupIds}, nil

}

// DeleteNetwork is a customized method for deletion of network, if one needs to delete the individual components of network then call the appropriate methods.
func (d *DeleteNetworkInput) DeleteNetwork(con aws.EstablishConnectionInput) (DeleteNetworkResponse, error) {

	vpcin := GetNetworksInput{VpcIds: d.VpcIds}
	vpc, err := vpcin.FindVpcs(con)
	if err != nil {
		return DeleteNetworkResponse{}, err
	}

	if vpc != true {
		return DeleteNetworkResponse{}, fmt.Errorf("Could not find the entered VPC, please enter valid/existing VPC id")
	}

	networkdel, neterr := d.getNetworkDeletables(con)
	if neterr != nil {
		return DeleteNetworkResponse{}, neterr
	}

	deletestatus, netdelerr := networkdel.deleteNetworkDeletables(con)
	if netdelerr != nil {
		return DeleteNetworkResponse{}, netdelerr
	}
	return deletestatus, nil
}

func (d *DeleteNetworkInput) deleteNetworkDeletables(con aws.EstablishConnectionInput) (DeleteNetworkResponse, error) {

	/*ec2, seserr := con.EstablishConnection()
	  if seserr != nil {
	          return DeleteNetworkResponse{}, seserr
	  }*/

	if len(d.SecIds) != 0 {
		//Deletion of security groups
		delsecin := NetworkComponentInput{SecGroupIds: d.SecIds}
		secdelerr := delsecin.DeleteSecutiryGroup(con)
		if secdelerr != nil {
			return DeleteNetworkResponse{}, secdelerr
		}
	}

	if len(d.RouteTableIds) != 0 {
		//dissassociating the routetable before deleteing it.
		route := NetworkComponentInput{RouteTableIds: d.RouteTableIds}
		dessroutetable, dessrouterr := route.DisassociateRouteTable(con)
		if dessrouterr != nil {
			return DeleteNetworkResponse{}, dessrouterr
		}
		if dessroutetable != true {
			return DeleteNetworkResponse{}, fmt.Errorf("An error occurred while dettaching routetable from subnet")
		}

		//deletion of routetable is handled by below loop.
		delrouterr := route.DeleteRouteTable(con)
		if delrouterr != nil {
			return DeleteNetworkResponse{}, delrouterr
		}
	}

	if len(d.IgwIds) != 0 {
		//dettachment of igw is been done by below snippet.
		dettachgateway := NetworkComponentInput{IgwIds: d.IgwIds, VpcIds: d.VpcIds}
		detacherr := dettachgateway.DetachIgws(con)
		if detacherr != nil {
			return DeleteNetworkResponse{}, detacherr
		}

		//deletion of igw is been done by below snippet.
		deletegateway := NetworkComponentInput{IgwIds: d.IgwIds}
		deleteigwerr := deletegateway.DeleteIgws(con)
		if deleteigwerr != nil {
			return DeleteNetworkResponse{}, deleteigwerr
		}
	}

	if len(d.SubnetIds) != 0 {
		subdelin := DeleteNetworkInput{SubnetIds: d.SubnetIds}
		subdelerr := subdelin.DeleteSubnets(con)
		if subdelerr != nil {
			return DeleteNetworkResponse{}, subdelerr
		}
	}

	//deletion of vpc is handled by below snippet
	deletevpc := DeleteNetworkInput{VpcIds: d.VpcIds}
	deletevpcerr := deletevpc.DeleteVpc(con)
	if deletevpcerr != nil {
		return DeleteNetworkResponse{}, deletevpcerr
	}
	return DeleteNetworkResponse{Status: "Network and all its components has been deleted successfully"}, nil
}

func (d *DeleteNetworkInput) getNetworkDeletables(con aws.EstablishConnectionInput) (DeleteNetworkInput, error) {

	//creating a session to perform actions
	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return DeleteNetworkInput{}, seserr
	}

	//Getting list of all subnets available in the network
	subnetres, suberr := ec2.DescribeSubnet(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if suberr != nil {
		return DeleteNetworkInput{}, suberr
	}

	subnets := make([]string, 0)
	for _, subnet := range subnetres.Subnets {
		if *subnet.DefaultForAz != true {
			subnets = append(subnets, *subnet.SubnetId)
		}
	}

	//Getting list of all secutiry groups in the entered vpc.
	secres, secerr := ec2.DescribeSecurityGroup(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if secerr != nil {
		return DeleteNetworkInput{}, secerr
	}
	secids := make([]string, 0)
	for _, sec := range secres.SecurityGroups {
		if *sec.GroupName != "default" {
			secids = append(secids, *sec.GroupId)
		}
	}

	//describing all the routetables to fetch the right ones.
	routres, routerr := ec2.DescribeRouteTable(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if routerr != nil {
		return DeleteNetworkInput{}, routerr
	}
	routeids := make([]string, 0)
	for _, route := range routres.RouteTables {
		if route.Associations != nil {
			if *route.Associations[0].Main != true {
				routeids = append(routeids, *route.RouteTableId)
			}
		} else {
			routeids = append(routeids, *route.RouteTableId)
		}
	}

	//describing all internet-gateways to get right one.
	response, err := ec2.DescribeIgw(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name:  "attachment.vpc-id",
				Value: d.VpcIds,
			},
		},
	)
	if err != nil {
		return DeleteNetworkInput{}, err
	}
	igwids := make([]string, 0)
	for _, igw := range response.InternetGateways {
		igwids = append(igwids, *igw.InternetGatewayId)
	}

	//collating the data of entire network which was collected.
	deleteResponse := new(DeleteNetworkInput)
	deleteResponse.SubnetIds = subnets
	deleteResponse.SecIds = secids
	deleteResponse.RouteTableIds = routeids
	deleteResponse.IgwIds = igwids
	deleteResponse.VpcIds = d.VpcIds

	return *deleteResponse, nil

}

// GetNetwork is a customized method for fetching the details of network, if one needs to fetch the information of the individual components of network then call the appropriate methods else call the GOD, interface which talks to cloud.
func (net *GetNetworksInput) GetNetwork(con aws.EstablishConnectionInput) ([]NetworkResponse, error) {

	networkresponse := make([]NetworkResponse, 0)
	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return nil, seserr
	}

	findvpcresult, vpcerr := ec2.DescribeVpc(
		&aws.DescribeNetworkInput{
			VpcIds: net.VpcIds,
		},
	)

	if vpcerr != nil {
		return nil, vpcerr
	}

	for _, vpc := range findvpcresult.Vpcs {

		// getting list of subnets in network
		subnetin := GetNetworksInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		subnets, suberr := subnetin.GetSubnetsFromVpc(con)
		if suberr != nil {
			return nil, suberr
		}

		// getting list of igws in network
		igwin := NetworkComponentInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		igw, igwerr := igwin.GetIgwFromVpc(con)
		if igwerr != nil {
			return nil, igwerr
		}

		// getting list of security group in network
		sec, secerr := igwin.GetSecFromVpc(con)
		if secerr != nil {
			return nil, secerr
		}

		if net.GetRaw == true {
			subnets.GetVpcRaw = vpc
			subnets.DescribeSecurityRaw = sec.GetSecurityRaw
			subnets.DescribeIgwRaw = igw.GetIgwRaw
			networkresponse = append(networkresponse, subnets)
		} else {
			if vpc.Tags != nil {
				networkresponse = append(networkresponse, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], IsDefault: *vpc.IsDefault, SecGroupIds: sec.SecGroupIds})
			} else {
				networkresponse = append(networkresponse, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], IsDefault: *vpc.IsDefault, SecGroupIds: sec.SecGroupIds})
			}
		}
	}
	return networkresponse, nil
}

// GetAllNetworks is a customized method for fetching the details of all the network in the specified region, if one needs to fetch the information of the individual network then call the appropriate method else call the GOD, interface which talks to cloud.
func (net *GetNetworksInput) GetAllNetworks(con aws.EstablishConnectionInput) ([]NetworkResponse, error) {

	networkresponse := make([]NetworkResponse, 0)
	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return nil, seserr
	}

	findvpcresult, vpcerr := ec2.DescribeAllVpc(
		&aws.DescribeNetworkInput{},
	)
	if vpcerr != nil {
		return nil, vpcerr
	}
	for _, vpc := range findvpcresult.Vpcs {

		// getting list of subnets in network
		subnetin := GetNetworksInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		subnets, suberr := subnetin.GetSubnetsFromVpc(con)
		if suberr != nil {
			return nil, suberr
		}

		// getting list of igws in network
		igwin := NetworkComponentInput{VpcIds: []string{*vpc.VpcId}, GetRaw: net.GetRaw}
		igw, igwerr := igwin.GetIgwFromVpc(con)
		if igwerr != nil {
			return nil, igwerr
		}

		// getting list of security group in network
		sec, secerr := igwin.GetSecFromVpc(con)
		if secerr != nil {
			return nil, secerr
		}

		if net.GetRaw == true {
			netres := new(NetworkResponse)
			netres.GetVpcRaw = vpc
			netres.GetSubnetRaw = subnets.GetSubnetRaw
			netres.DescribeSecurityRaw = sec.GetSecurityRaw
			netres.DescribeIgwRaw = igw.GetIgwRaw
			networkresponse = append(networkresponse, *netres)
		} else {
			if vpc.Tags != nil {
				networkresponse = append(networkresponse, NetworkResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], SecGroupIds: sec.SecGroupIds, IsDefault: *vpc.IsDefault, Region: con.Region})
			} else {
				networkresponse = append(networkresponse, NetworkResponse{VpcId: *vpc.VpcId, Subnets: subnets.Subnets, State: *vpc.State, IgwId: igw.IgwIds[0], SecGroupIds: sec.SecGroupIds, IsDefault: *vpc.IsDefault, Region: con.Region})
			}
		}
	}
	return networkresponse, nil
}

// UpdateNetwork is a customized method for updating the network and its components, if one needs to update the individual components network then this method does just that. For more operations call GOD, interface which talks to cloud.
func (net *UpdateNetworkInput) UpdateNetwork(con aws.EstablishConnectionInput) (NetworkResponse, error) {

	/*ec2, seserr := con.EstablishConnection()
	  if seserr != nil {
	          return NetworkResponse{}, seserr
	  }*/

	switch strings.ToLower(net.Resource) {
	case "subnets":

		switch strings.ToLower(net.Action) {
		case "create":
			// Collects all the available availability zones
			var zones []string
			if net.Network.Zone == "" {
				zonein := CommonInput{}
				zone, zonerr := zonein.GetAvailabilityZones(con)
				if zonerr != nil {
					return NetworkResponse{}, zonerr
				}
				zones = zone
			} else {
				zones = []string{net.Network.Zone}
			}
			// I will be the spoc for subnets creation in the loop as per the request made
			subnetresponse := make([]SubnetReponse, 0)
			zonenum := len(zones) - 1

			//Fetching unique number to give our subnet a unique name
			subnets := make([]string, 0)
			subnetin := GetNetworksInput{VpcIds: []string{net.Network.VpcId}}
			subnetlist, suberr := subnetin.GetSubnetsFromVpc(con)
			if suberr != nil {
				return NetworkResponse{}, suberr
			}
			for _, subnet := range subnetlist.Subnets {
				if subnet.Name != "" {
					subnets = append(subnets, subnet.Name)
				}
			}

			// Getting Unique digit to name subnet uniquly
			uqnin := CommonInput{SortInput: subnets}
			uqnchr, unerr := uqnin.GetUniqueNumberFromTags()
			if unerr != nil {
				return NetworkResponse{}, unerr
			}
			for _, sub := range net.Network.SubCidrs {
				if zonenum < 0 {
					zonenum = len(zones) - 1
				}

				// Creating subnet by calling appropriate object
				subin := NetworkCreateInput{
					SubCidr: sub,
					Name:    net.Network.Name + "_sub" + strconv.Itoa(uqnchr),
					Zone:    zones[zonenum],
					VpcId:   net.Network.VpcId,
					GetRaw:  net.GetRaw,
				}
				subnet, suberr := subin.CreateSubnet(con)
				if suberr != nil {
					return NetworkResponse{}, suberr
				}
				subnetresponse = append(subnetresponse, subnet)

				zonenum--
				uqnchr++
			}
			if net.GetRaw == true {
				return NetworkResponse{CreateSubnetRaw: subnetresponse}, nil
			}
			return NetworkResponse{Subnets: subnetresponse}, nil
		case "delete":
			return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting the action %s of the resource %s or you entered wrong name. The action you selected was: %s", net.Action, net.Resource, net.Action))
		default:
			return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting the action %s of the resource %s or you entered wrong name. The action you selected was: %s", net.Action, net.Resource, net.Action))
		}

	case "vpc":
		return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting updation of the resource you entered or you entered wrong name. The resource you enetered was: %s", net.Resource))
	case "igw":
		return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting updation of the resource you entered or you entered wrong name. The resource you enetered was: %s", net.Resource))
	default:
		return NetworkResponse{}, fmt.Errorf(fmt.Sprintf("Either we are not supporting updation of the resource you entered or you entered wrong name. The resource you enetered was: %s", net.Resource))
	}
}
