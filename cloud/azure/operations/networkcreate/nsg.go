package main

import (
	"azure/interface/networkinterface"
	"encoding/json"
	"fmt"
)

func main() {
	/*         c := azurenetwork.NsgIn{"M1038273", "test-nsg", "CentralIndia"}
	           create, _ := c.CreateNetworkSecurityGroup()
	           json, _ := json.MarshalIndent(create, "", "  ")
	           fmt.Println(string(json))
	*/
	l := azurenetwork.NsgIn{ResourceGroup: "M1038273"}
	yo, _ := l.ListNetworkSecurityGroup()
	jo, _ := json.MarshalIndent(yo, "", "  ")
	fmt.Println(string(jo))
	/*
	   r := azurenetwork.NsgIn{ResourceGroup: "M1038273", NsgName: "test-nsg"}
	   resp, _ := r.DeleteNetworkSecurityGroup()
	   if resp.Response.StatusCode == 200 { fmt.Println("Deleted network "+r.NsgName+" successfully")}
	*/
}
