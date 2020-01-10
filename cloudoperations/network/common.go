// Package netcommon hosts all the responses that cloudoperations give out.
// This is for future development currently this is nothing to do with.
package netcommon

import (
	network "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/network"
)

// NetworkResponse is a struct that will return the filtered/unfiltered responses of variuos clouds.
type NetworkResponse struct {
	// Contains filtered/unfiltered response of AWS.
	AwsResponse network.NetworkResponse `json:"AwsResponse,omitempty"`

	// Contains filtered/unfiltered response of Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`

	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}
