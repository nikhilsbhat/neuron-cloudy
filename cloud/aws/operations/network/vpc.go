package awsnetwork

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	common "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/common"
)

// VpcResponse is a struct that will be the response type of almost all the VPC related activities under cloud/operations.
type VpcResponse struct {
	Name              string                           `json:"Name,omitempty"`
	Type              string                           `json:"Type,omitempty"`
	VpcId             string                           `json:"VpcId,omitempty"`
	IgwId             string                           `json:"IgwId,omitempty"`
	SecGroupIds       []string                         `json:"SecGroupId,omitempty"`
	IsDefault         bool                             `json:"IsDefault,omitempty"`
	State             string                           `json:"State,omitempty"`
	CreateVpcRaw      *ec2.CreateVpcOutput             `json:"CreateVpcRaw,omitempty"`
	GetVpcRaw         *ec2.DescribeVpcsOutput          `json:"GetVpcRaw,omitempty"`
	CreateIgwRaw      *ec2.CreateInternetGatewayOutput `json:"CreateIgwRaw,omitempty,omitempty"`
	CreateSecurityRaw *ec2.CreateSecurityGroupOutput   `json:"CreateSecRaw,omitempty,omitempty"`
}

// CreateVpc is a customized method for vpc creation, if one needs plain vpc creation then he/she has to call the GOD, interface which talks to cloud.
func (vpc *NetworkCreateInput) CreateVpc(con aws.EstablishConnectionInput) (VpcResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return VpcResponse{}, seserr
	}
	// I am gathering inputs since create vpc needs it
	vpcResult, vpcErr := ec2.CreateVpc(
		&aws.CreateNetworkInput{
			Cidr:    vpc.VpcCidr,
			Tenancy: "default",
		})

	// handling the error if it throws while vpc is under creation process
	if vpcErr != nil {
		return VpcResponse{}, vpcErr
	}

	// I will program wait until vpc become available
	waitErr := ec2.WaitTillVpcAvailable(
		&aws.DescribeNetworkInput{
			Filters: aws.Filters{
				Name: "vpc-id", Value: []string{*vpcResult.Vpc.VpcId},
			},
		},
	)
	if waitErr != nil {
		return VpcResponse{}, waitErr
	}

	// I will pass name to create_tags to set a name to the vpc
	vpctagin := common.Tag{*vpcResult.Vpc.VpcId, "Name", vpc.Name}
	vpctag, tagErr := vpctagin.CreateTags(con)
	if tagErr != nil {
		return VpcResponse{}, tagErr
	}

	// Decision will be made here, whether we need public network or private, based on the input I receive
	netcomp := new(NetworkComponentInput)
	netcomp.Name = vpc.Name
	netcomp.VpcIds = []string{*vpcResult.Vpc.VpcId}
	netcomp.GetRaw = vpc.GetRaw
	vpcresponse := new(VpcResponse)

	if (strings.ToLower(vpc.Type) == "public") || (strings.ToLower(vpc.Type) == "") {
		ig, igErr := netcomp.CreateIgw(con)
		if igErr != nil {
			return VpcResponse{}, igErr
		}
		if vpc.GetRaw != true {
			vpcresponse.IgwId = ig.IgwIds[0]
		} else {
			vpcresponse.CreateIgwRaw = ig.CreateIgwRaw
		}
	} else if strings.ToLower(vpc.Type) == "private" {
		vpcresponse.IgwId = ""
	} else {
		return VpcResponse{}, fmt.Errorf("You provided unknown network type. There are two possibility, either we do not support this type else you would have misspelled")
	}

	// I will initialize data required to create security group and pass it to respective person to create one
	netcomp.Ports = vpc.Ports
	sec, secErr := netcomp.CreateSecurityGroup(con)
	if secErr != nil {
		return VpcResponse{}, secErr
	}

	if vpc.GetRaw == true {
		vpcresponse.CreateSecurityRaw = sec.CreateSecurityRaw
		vpcresponse.CreateVpcRaw = vpcResult
		return *vpcresponse, nil
	}

	vpcresponse.SecGroupIds = sec.SecGroupIds
	vpcresponse.Name = vpctag
	vpcresponse.VpcId = *vpcResult.Vpc.VpcId
	vpcresponse.Type = vpc.Type
	return *vpcresponse, nil
}

// DeleteVpc is a customized method for vpc deletion, if one needs plain vpc deletion then he/she has to call the GOD, interface which talks to cloud.
func (vpc *DeleteNetworkInput) DeleteVpc(con aws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}

	err := ec2.DeleteVpc(
		&aws.DescribeNetworkInput{
			VpcIds: vpc.VpcIds,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// GetVpcs is a customized method for fetching details of all vpc for a given region, if one needs plain get subnet then he/she has to call the GOD, interface which talks to cloud.
func (v *GetNetworksInput) GetVpcs(con aws.EstablishConnectionInput) (NetworkResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkResponse{}, seserr
	}
	response, err := ec2.DescribeAllVpc(
		&aws.DescribeNetworkInput{
			VpcIds: v.VpcIds,
		},
	)
	if err != nil {
		return NetworkResponse{}, err
	}

	if v.GetRaw == true {
		return NetworkResponse{GetVpcsRaw: response}, nil
	}

	vpcs := make([]VpcResponse, 0)
	for _, vpc := range response.Vpcs {
		vpcs = append(vpcs, VpcResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, State: *vpc.State, IsDefault: *vpc.IsDefault})
	}
	return NetworkResponse{Vpcs: vpcs}, nil
}

// FindVpcs is a customized method which sends back the response to the caller about the existence of vpc asked for.
func (v *GetNetworksInput) FindVpcs(con aws.EstablishConnectionInput) (bool, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return false, seserr
	}
	response, err := ec2.DescribeVpc(
		&aws.DescribeNetworkInput{
			VpcIds: v.VpcIds,
		},
	)
	if err != nil {
		return false, err
	}
	if len(response.Vpcs) != 0 {
		return true, nil
	}
	return false, fmt.Errorf("Could not find the VPC's you asked for")
}
