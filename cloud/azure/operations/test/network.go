package network

import (
	"azure/interface/networkinterface"
	//         "fmt"
	//         "encoding/json"
	"strconv"
)

type NetworkIn struct {
	ResourceGroup string
	Name          string
	Cidr          string
	SubnetCidr    []string
	Port          []string
	Location      string
}

type SubnetOut struct {
	Name  string
	ID    string
	NsgID string
}

type VnetOut struct {
	Name string
	ID   string
}
type NetworkOut struct {
	Vnet   VnetOut
	Subnet []SubnetOut
}

func (n NetworkIn) Create() NetworkOut {

	c := azurenetwork.VnetIn{n.ResourceGroup, n.Name, n.Cidr, n.Location}
	create, _ := c.CreateVirtualNetwork()
	//         json, _ := json.MarshalIndent(create, "", "  ")
	//         fmt.Println(string(json))
	var p VnetOut
	p = VnetOut{*create.Name, *create.ID}
	var q []SubnetOut
	for i, s := range n.SubnetCidr {
		j := strconv.Itoa(i)
		ng := azurenetwork.NsgIn{n.ResourceGroup, n.Name + "_nsg_" + j, n.Location}
		nsg, _ := ng.CreateNetworkSecurityGroup()
		for i, ns := range n.Port {
			ns := azurenetwork.SecurityRuleIn{n.ResourceGroup, n.Name + "_nsg_" + j, "Port_" + ns, ns, int32(100) + int32(i)}
			ns.CreateNetworkSecurityRule()
		}

		sub := azurenetwork.SubnetIn{n.ResourceGroup, n.Name, "subnet_" + j, s, *nsg.ID}
		subnet, _ := sub.CreateVirtualNetworkSubnet()
		q = append(q, SubnetOut{*subnet.Name, *subnet.ID, *subnet.SubnetPropertiesFormat.NetworkSecurityGroup.ID})
	}

	/*                ng := azurenetwork.NsgIn{n.ResourceGroup,n.Name + "_nsg_" + j, *subnet.ID, n.Location}
	                ng.CreateNetworkSecurityGroup()
	                for i, ns := range n.Port {
	                         ns := azurenetwork.SecurityRuleIn{n.ResourceGroup, n.Name + "_nsg_" + j, "Port_" + ns, ns,int32(100) + int32(i)}
	                         ns.CreateNetworkSecurityRule()
	                }
		 }

		 ng := azurenetwork.NsgIn{n.ResourceGroup,n.Name + "_nsg", n.Location}
		 ng.CreateNetworkSecurityGroup()
	         for i, ns := range n.Port {
	                 ns := azurenetwork.SecurityRuleIn{n.ResourceGroup, n.Name + "_nsg", "Port_" + ns, ns,int32(100) + int32(i)}
	                 ns.CreateNetworkSecurityRule()
		 }
	*/
	return NetworkOut{p, q}
}
