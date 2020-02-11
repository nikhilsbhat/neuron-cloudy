// Package support will make a call whether this tool supports the invoked cloud.
package support

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	cloud "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

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

// ValidateClient validates the client validity against the cloud passed.
func ValidateClient(c *cloud.Cloud) bool {
	switch (c.Client).(type) {
	case *session.Session:
		if strings.ToLower(c.Name) == "aws" {
			return true
		}
		return false
	case *http.Client:
		if strings.ToLower(c.Name) == "gcp" {
			return true
		}
		return false
	// to be enabled once azure is configured
	// case <azure_client>:
	// 	if strings.ToLower(c.Name) == "azure" {
	// 		return true
	// 	}
	// 	return false
	default:
		return false
	}
}
