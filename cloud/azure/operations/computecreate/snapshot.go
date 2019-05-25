package main

import (
	"dengine/interface/azureinterface/computeinterface"
	"encoding/json"
	"fmt"
)

func main() {

	m := azurecompute.VMIn{ResourceGroup: "M1038273", VmName: "ubuntu"}
	vm, _ := m.GetVM()

	osDisk := vm.VirtualMachineProperties.StorageProfile.OsDisk.ManagedDisk
	//        fmt.Println(*osDisk.ID)

	s := azurecompute.SnapshotIn{"M1038273", "test-snap1", *osDisk.ID, "CentralIndia"}
	//        snapshot, _ := s.CreateSnapshot()
	s.CreateSnapshot()
	//        jsonMarshall(snapshot)

	r := azurecompute.SnapshotIn{ResourceGroup: "M1038273", SnapshotName: "test-snap1"}
	sn, _ := r.GetSnapshot()
	jsonMarshall(sn)

	//        resp, _ := r.DeleteSnapshot()
	//        if resp.Response.StatusCode == 200 { fmt.Println("Deleted snapshot "+r.SnapshotName+" successfully")}

}

func jsonMarshall(i interface{}) {
	json, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(json))
}
