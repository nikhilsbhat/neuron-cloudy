// Package support will make a call whether this tool supports the passed cloud.
package support

// Listing the names of the cloud that this tool supports,
// once the compactable for new cloud is made just make an entry here.
var clouds = []string{"aws", "azure", "gcp", "openstack"}

// DoesCloudSupports is the place where the actual decision for the clous is made and will return status to the called method.
// This has to be used by cloudoperations.
func DoesCloudSupports(input string) bool {
	for _, value := range clouds {
		if input == value {
			return true
		}
	}
	return false
}
