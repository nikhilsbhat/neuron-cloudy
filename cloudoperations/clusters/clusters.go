// Package clusters makes the tool cloud agnostic with respect to the resource cluster.
// The decision will be made here to route the request to respective package based on input.
package clusters

import (
	gcp "github.com/nikhilsbhat/neuron-cloudy/cloud/gcp/operations"
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// GetClusterInput is the way in to fetch the details of the cluster.
type GetClusterInput struct {
	// Inputs to fetch the information of appropriate cluster.
	//Gcp   gcp.GetClusterInput
	// ClusterName refers to the name of the Kubernetes cluster of which the information has to be fetched.
	ClusterName string
	// ProjectID refers to the ID of the project of whom the cluster belongs to.
	ProjectID string
	// Regions holds the names of the region to which the resourc ebelongs to.
	Regions []string
	Cloud   cmn.Cloud
}

// ClusterResponse returns the filtered/unfiltered responses of variuos clouds.
type ClusterResponse struct {
	// Contains filtered/unfiltered response from AWS.
	AwsResponse string `json:"AwsResponse,omitempty"`
	// Contains filtered/unfiltered response from Azure.
	AzureResponse string `json:"AzureResponse,omitempty"`
	// Contains filtered/unfiltered response from GCP.
	GCPResponse []gcp.ClusterResponse `json:"GcpResponse,omitempty"`
	// Default response if no inputs or matching the values required.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
}

//Nothing much from this file. This file contains only the structs for network/create
