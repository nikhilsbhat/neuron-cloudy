package main

import (
	"azure/interface/networkinterface"
	"encoding/json"
	"fmt"
)

func main() {
	c := azurenetwork.SecurityRuleIn{"M1038273", "test-nsg", "ssh", "22", 123}
	create, _ := c.CreateNetworkSecurityRule()
	json, _ := json.MarshalIndent(create, "", "  ")

	fmt.Println(string(json))
}
