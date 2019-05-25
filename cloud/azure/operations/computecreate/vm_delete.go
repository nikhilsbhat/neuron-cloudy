// this package will delete the VM along with the components associated with it.

package main

import (
	"dengine/interface/azureinterface/computeinterface"
	"dengine/interface/azureinterface/networkinterface"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {

	m := azurecompute.VMIn{ResourceGroup: "test", VmName: "test-vm"}
	vm, _ := m.GetVM()
	//        jsonMarshall(vm)

	var interfaceName string
	net := vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces
	for _, nictemp := range *net {
		interfaceName = nameFromID(*nictemp.ID)
	}

	n := azurenetwork.NicIn{ResourceGroup: "test", NicName: interfaceName}
	nic, _ := n.GetNIC()
	//        fmt.Println(*nic.Name)

	nsg := nameFromID(*nic.NetworkSecurityGroup.ID)
	//        fmt.Println(nsg)

	var ip string
	pubip := nic.InterfacePropertiesFormat.IPConfigurations
	for _, iptemp := range *pubip {
		ip = nameFromID(*iptemp.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress.ID)
	}

	//        fmt.Println(publicIp)
	osDisk := *vm.VirtualMachineProperties.StorageProfile.OsDisk.Name
	//        fmt.Println(osDisk)

	//************************@@@VM DELETE@@@***********************
	vd, _ := m.DeleteVM()
	if vd.Response.StatusCode == 200 {
		fmt.Println("Deleted VM " + m.VmName + " successfully")
	}
	//***********************@@@NIC DELETE@@@***********************
	nd, _ := n.DeleteNIC()
	if nd.Response.StatusCode == 200 {
		fmt.Println("Deleted NIC " + n.NicName + " successfully")
	}
	//***********************@@@NSG DELETE@@@***********************
	s := azurenetwork.NsgIn{ResourceGroup: "test", NsgName: nsg}
	sd, _ := s.DeleteNetworkSecurityGroup()
	if sd.Response.StatusCode == 200 {
		fmt.Println("Deleted NSG " + s.NsgName + " successfully")
	}
	//************************@@@IP DELETE@@@***********************
	i := azurenetwork.IpIn{ResourceGroup: "test", IpName: ip}
	id, _ := i.DeletePublicIP()
	if id.Response.StatusCode == 200 {
		fmt.Println("Deleted IP " + i.IpName + " successfully")
	}
	//***********************@@@DISK DELETE@@***********************
	d := azurecompute.DisksIn{ResourceGroup: "test", DiskName: osDisk}
	dd, _ := d.DeleteDisk()
	if dd.Response.StatusCode == 200 {
		fmt.Println("Deleted DISK " + d.DiskName + " successfully")
	}
}

func nameFromID(s string) string {
	slice := strings.Split(s, "/")
	return slice[len(slice)-1]
}

func jsonMarshall(i interface{}) {
	json, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(json))
}
