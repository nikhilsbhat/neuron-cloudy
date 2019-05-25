package azurecompute

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"io/ioutil"
	"log"
	"neuron/cloud/azure/access"
	"os"
)

var (
	token, _, subscription = auth.GetServicePrincipalToken()
	fakepubkey             = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCjCX8wh0lnk2KUvoCulBER4TQ+4+repQF5vvQeCVc5eWHNQKIPuSxy4fGcEbar15U4wjEJYDsXUGhW0JIh4peIKFf+dXUtZlMQEo7QvPGGORjVm8Zf+je/cVqGQJOvUP4s1/J8EQ+/n6gidtByBL+4lN/vDp/lgPSZzRgb08zVuW40z6jFrxfwalru10FHzzPmkCEtW54YkdJ2yEnLzk+xZDJXmG7JE4c2yRl+Y35HCzHfeRsUqcF1ErV2KYHcRWqwzD9oDZ5V2uTC4ERHkF102Ve7LOSyYK3cvJ8QSWMoOCOPA/UpdrkJRq9e2eVdpIqvnbu2vp6xazU080ZNu/BB"
	ctx                    = context.Background()
)

func getVMClient() compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(subscription)
	vmClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return vmClient
}

type VMIn struct {
	ResourceGroup    string
	VmName           string `json:"vmname,omitempty"`
	NicID            string `json:"nicid,omitempty"`
	UserName         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	Flavour          string `json:"os,omitempty"`
	SSHPublicKeyPath string `json:"sshkeypath,omitempty"`
	Location         string `json:"location,omitempty"`
}

// CreateVM creates a new virtual machine with the specified name using the specified NIC.
// Username, password, and sshPublicKeyPath determine logon credentials.
func (v VMIn) CreateVM() (vm compute.VirtualMachine, err error) {
	// see the network samples for how to create and get a NIC resource

	osconfig := Image(v.Flavour)

	var sshKeyData string
	if _, err = os.Stat(v.SSHPublicKeyPath); err == nil {
		sshBytes, err := ioutil.ReadFile(v.SSHPublicKeyPath)
		if err != nil {
			log.Fatalf("failed to read SSH key data: %v", err)
		}
		sshKeyData = string(sshBytes)
	} else {
		sshKeyData = fakepubkey
	}

	vmClient := getVMClient()
	future, err := vmClient.CreateOrUpdate(
		ctx,
		v.ResourceGroup,
		v.VmName,
		compute.VirtualMachine{
			Location: to.StringPtr(v.Location),
			VirtualMachineProperties: &compute.VirtualMachineProperties{
				HardwareProfile: &compute.HardwareProfile{
					VMSize: compute.BasicA0,
				},
				StorageProfile: &compute.StorageProfile{
					ImageReference: &compute.ImageReference{
						Publisher: to.StringPtr(osconfig.Publisher),
						Offer:     to.StringPtr(osconfig.Offer),
						Sku:       to.StringPtr(osconfig.Sku),
						Version:   to.StringPtr(osconfig.Version),
					},
				},
				OsProfile: &compute.OSProfile{
					ComputerName:  to.StringPtr(v.VmName),
					AdminUsername: to.StringPtr(v.UserName),
					AdminPassword: to.StringPtr(v.Password),
					LinuxConfiguration: &compute.LinuxConfiguration{
						SSH: &compute.SSHConfiguration{
							PublicKeys: &[]compute.SSHPublicKey{
								{
									Path: to.StringPtr(
										fmt.Sprintf("/home/%s/.ssh/authorized_keys",
											v.UserName)),
									KeyData: to.StringPtr(sshKeyData),
								},
							},
						},
					},
				},
				NetworkProfile: &compute.NetworkProfile{
					NetworkInterfaces: &[]compute.NetworkInterfaceReference{
						{
							ID: to.StringPtr(v.NicID),
							NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return vm, fmt.Errorf("cannot create vm: %v", err)
	}

	err = future.WaitForCompletion(ctx, vmClient.Client)
	if err != nil {
		return vm, fmt.Errorf("cannot get the vm create or update future response: %v", err)
	}

	return future.Result(vmClient)
}

func (v VMIn) DeleteVM() (ar autorest.Response, err error) {

	vmClient := getVMClient()
	future, err := vmClient.Delete(
		ctx,
		v.ResourceGroup,
		v.VmName,
	)
	if err != nil {
		return ar, fmt.Errorf("cannot delete VM: %v", err)
	}

	err = future.WaitForCompletion(ctx, vmClient.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get the VM delete future response: %v", err)
	}

	return future.Result(vmClient)
}

func (v VMIn) GetVM() (vm compute.VirtualMachine, err error) {

	vmClient := getVMClient()
	future, err := vmClient.Get(
		ctx,
		v.ResourceGroup,
		v.VmName,
		"")

	if err != nil {
		return vm, fmt.Errorf("cannot get virtual VM: %v", err)
	}

	return future, err
}

func (v VMIn) ListVM() (vm []compute.VirtualMachine, err error) {

	vmClient := getVMClient()
	future, err := vmClient.List(
		ctx,
		v.ResourceGroup,
	)

	if err != nil {
		return vm, fmt.Errorf("cannot list the VMs in a resourcegroup: %v", err)
	}

	return future.Values(), err
}

func ListAllVM() (vm []compute.VirtualMachine, err error) {

	vmClient := getVMClient()
	future, err := vmClient.ListAll(
		ctx,
	)

	if err != nil {
		return vm, fmt.Errorf("cannot list the VMs: %v", err)
	}

	return future.Values(), err
}
