package azurenetwork

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

func getNicClient() network.InterfacesClient {
	nicClient := network.NewInterfacesClient(subscription)
	nicClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return nicClient
}

type NicIn struct {
	ResourceGroup string
	NicName       string `json:"nicname,omitempty"`
	NsgID         string `json:"nsgid,omitempty"`
	SubnetID      string `json:"subnetid,omitempty"`
	IpID          string `json:"ipid,omitempty"`
	Location      string `json:"location,omitempty"`
}

// CreateNIC creates a new network interface.

func (n NicIn) CreateNIC() (nic network.Interface, err error) {

	nicParams := network.Interface{
		Name:     to.StringPtr(n.NicName),
		Location: to.StringPtr(n.Location),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &[]network.InterfaceIPConfiguration{
				{
					Name: to.StringPtr(n.NicName + "-ipConfig1"),
					InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
						Subnet: &network.Subnet{
							ID: to.StringPtr(n.SubnetID),
						},
						PrivateIPAllocationMethod: network.Dynamic,
						PublicIPAddress: &network.PublicIPAddress{
							ID: to.StringPtr(n.IpID),
						},
					},
				},
			},
		},
	}

	if n.NsgID != "" {
		nicParams.NetworkSecurityGroup = &network.SecurityGroup{
			ID: to.StringPtr(n.NsgID),
		}
	}

	nicClient := getNicClient()
	future, err := nicClient.CreateOrUpdate(
		ctx,
		n.ResourceGroup,
		n.NicName,
		nicParams,
	)

	if err != nil {
		return nic, fmt.Errorf("cannot create nic: %v", err)
	}

	err = future.WaitForCompletion(ctx, nicClient.Client)
	if err != nil {
		return nic, fmt.Errorf("cannot get nic create or update future response: %v", err)
	}

	return future.Result(nicClient)
}

func (n NicIn) DeleteNIC() (ar autorest.Response, err error) {
	nicClient := getNicClient()
	future, err := nicClient.Delete(
		ctx,
		n.ResourceGroup,
		n.NicName,
	)
	if err != nil {
		return ar, fmt.Errorf("cannot delete nic: %v", err)
	}

	err = future.WaitForCompletion(ctx, nicClient.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get nic delete future response: %v", err)
	}

	return future.Result(nicClient)
}

func (n NicIn) GetNIC() (nic network.Interface, err error) {
	nicClient := getNicClient()
	future, err := nicClient.Get(
		ctx,
		n.ResourceGroup,
		n.NicName,
		"")
	if err != nil {
		return nic, fmt.Errorf("cannot list get: %v", err)
	}

	return future, err
}

func (n NicIn) ListNIC() (nic []network.Interface, err error) {
	nicClient := getNicClient()
	future, err := nicClient.List(
		ctx,
		n.ResourceGroup,
	)

	if err != nil {
		return nic, fmt.Errorf("cannot list nic: %v", err)
	}

	return future.Values(), err
}

func ListAllNIC() (nic []network.Interface, err error) {
	nicClient := getNicClient()
	future, err := nicClient.ListAll(
		ctx,
	)

	if err != nil {
		return nic, fmt.Errorf("cannot list NICs: %v", err)
	}

	return future.Values(), err
}
