package main

import (
	"azure/interface/networkinterface"
	"encoding/json"
	"fmt"
)

func main() {
	c := azurenetwork.SubnetIn{"M1038273", "test", "test-sub", "192.168.10.0/24"}
	create, _ := c.CreateVirtualNetworkSubnet()
	//         json, _ := json.Marshal(create)
	//         fmt.Println(string(json))
	fmt.Println(create)

	l := azurenetwork.SubnetIn{ResourceGroup: "M1038273", VnetName: "test"}
	list, _ := l.ListVirtualNetworkSubnet()
	json_val, _ := json.MarshalIndent(list, "", "  ")
	fmt.Println(string(json_val))

	g := azurenetwork.SubnetIn{ResourceGroup: "M1038273", VnetName: "test", SubnetName: "test-sub"}
	val1, _ := g.GetVirtualNetworkSubnet()
	fmt.Println(val1)

	r := azurenetwork.SubnetIn{ResourceGroup: "M1038273", VnetName: "test", SubnetName: "test-sub"}
	resp, _ := r.DeleteVirtualNetworkSubnet()
	if resp.Response.StatusCode == 200 {
		fmt.Println("Deleted network " + r.SubnetName + " successfully")
	}

}
