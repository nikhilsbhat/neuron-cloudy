// Package gcp will help one in creating/deleting/updating/fetching informaation of various resources in gcp.
// But this is tailor made for this application if one needs
// customized result, he/she has to write one similar to this for them by calling the (master) interface.
// This package is capable of returining both custom response
// and raw from cloud depending on what you pass.
package gcp

import (
	"fmt"

	neurongcp "github.com/nikhilsbhat/neuron-cloudy/cloud/gcp/interface"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
)

// GetClusterInput holds the required values to fetch the cluster details
type GetClusterInput struct {
	// ClusterName refers to the name of the GKE cluster of which the information has to be fetched.
	ClusterName string
	// ProjectID refers to the ID of the GCP project of whom the cluster belongs to.
	ProjectID string
	// Regions holds the names of the region to which the resourc ebelongs to.
	Regions []string
	// GetRaw makes sure that function returns unfiltered response if it is set.
	GetRaw bool
	CredMode
}

// ClusterResponse contains filtered/unfiltered response from GCP on cluster details.
type ClusterResponse struct {
	// ClusterLocation refers to the list locations where the cluser components are spread across.
	ClusterLocation []string
	// ClusterNetwork refers to the google network of which the cluster is part of.
	ClusterNetwork string
	// ClusterName refers to the name of the cluster.
	ClusterName string
	// CreateTime is time when the cluster is created.
	CreateTime string
	// CurrentMasterVersion refers to the version of the master GKE node.
	CurrentMasterVersion string
	// ClusterAuth holds the authentication details of the cluster seleccted.
	ClusterAuth *container.MasterAuth
	// CurrentNodeVersion refers to the version of the GKE node.
	CurrentNodeVersion string
	// KubeEndpoint contains the url where kubernetes is accessible.
	KubeEndpoint string
	// ClusterNodePool contains the details of the nodepool of the cluster specified.
	ClusterNodePool []*container.NodePool
	// InstanceGroupUrls contains the details of the compute instances which are part of cluster.
	ClusterInstanceGroupUrls []string
	// Autoscaling holds more information on autoscaling configurations of the specified cluster.
	//ClusterAutoscaling *container.ClusterAutoscaling
	// GetClustersRaw contains unfiltered response from GCP on cluster
	GetClustersRaw []*container.Cluster
}

// CredMode holds the type of credentials passed.
type CredMode struct {
	// JSONPath refers to the path of GCP JSON credential file.
	JSONPath string
	// RawJSON could alternative to JSONPath, if one prefers to inject credential directly this is the way.
	RawJSON []byte
}

// GetClusters gets the list of kube clusters along with its other details
// This is region specific, mind the region you pass to it.
func (clust GetClusterInput) GetClusters(client interface{}) ([]ClusterResponse, error) {

	// Initialization of gcp client
	sess := getClientFromBase(client, []string{compute.CloudPlatformScope})

	if len(clust.ProjectID) == 0 {
		return nil, fmt.Errorf("Project ID cannot be nil")
	}

	if clust.Regions == nil {
		clust.Regions = append(clust.Regions, "-")
	}

	// The parent (project and location) where the clusters will be listed.
	// Specified in the format 'projects/*/locations/*'.
	// Location "-" matches all zones and all regions.
	//parent := "projects/core-search-dev/locations/us-central1" // TODO: Update placeholder value.
	clusters := make([]*container.Cluster, 0)
	for _, r := range clust.Regions {
		input := new(neurongcp.GetClusterInput)
		input.ResourceURL = fmt.Sprintf("projects/%s/locations/%s", clust.ProjectID, r)
		input.Client = sess
		resp, err := input.ListClusters()
		if err != nil {
			return nil, err
		}
		for _, cluster := range resp {
			clusters = append(clusters, cluster)
		}
	}
	response := make([]ClusterResponse, 0)
	if clust.GetRaw == true {
		return append(response, ClusterResponse{GetClustersRaw: clusters}), nil
	}

	// filters the response obtained from the GCP
	for _, cluster := range clusters {
		clust := new(ClusterResponse)
		clust.ClusterName = cluster.Name
		clust.KubeEndpoint = cluster.Endpoint
		clust.CurrentMasterVersion = cluster.CurrentMasterVersion
		clust.CurrentNodeVersion = cluster.CurrentNodeVersion
		clust.ClusterNetwork = cluster.Network
		clust.CreateTime = cluster.CreateTime
		clust.ClusterAuth = cluster.MasterAuth
		clust.ClusterLocation = cluster.Locations
		clust.ClusterNodePool = cluster.NodePools
		clust.ClusterInstanceGroupUrls = cluster.InstanceGroupUrls
		//clust.ClusterAutoscaling = cluster.Autoscaling
		response = append(response, *clust)
	}
	return response, nil
}

// GetCluster gets the details of specific kube cluster along with its other details
// This is region specific, mind the region you pass to it.
func (clust GetClusterInput) GetCluster(client interface{}) ([]ClusterResponse, error) {

	// Initialization of gcp client
	sess := getClientFromBase(client, []string{compute.CloudPlatformScope})

	if len(clust.ProjectID) == 0 {
		return nil, fmt.Errorf("Project ID cannot be empty while fetching the details of cluster")
	}

	if len(clust.Regions) == 0 {
		return nil, fmt.Errorf("Region cannot be empty while fetching the details of cluster")
	}

	clusters := make([]*container.Cluster, 0)
	input := new(neurongcp.GetClusterInput)
	input.ResourceURL = fmt.Sprintf("projects/%s/locations/%s/clusters/%s", clust.ProjectID, clust.Regions[0], clust.ClusterName)
	input.Client = sess
	cluster, err := input.GetCluster()
	if err != nil {
		return nil, err
	}
	clusters = append(clusters, cluster)

	response := make([]ClusterResponse, 0)
	if clust.GetRaw == true {
		return append(response, ClusterResponse{GetClustersRaw: clusters}), nil
	}

	// filters the response obtained from the GCP
	for _, cluster := range clusters {
		clust := new(ClusterResponse)
		clust.ClusterName = cluster.Name
		clust.KubeEndpoint = cluster.Endpoint
		clust.CurrentMasterVersion = cluster.CurrentMasterVersion
		clust.CurrentNodeVersion = cluster.CurrentNodeVersion
		clust.ClusterNetwork = cluster.Network
		clust.CreateTime = cluster.CreateTime
		clust.ClusterAuth = cluster.MasterAuth
		clust.ClusterLocation = cluster.Locations
		clust.ClusterNodePool = cluster.NodePools
		clust.ClusterInstanceGroupUrls = cluster.InstanceGroupUrls
		//clust.ClusterAutoscaling = cluster.Autoscaling
		response = append(response, *clust)
	}
	return response, nil
}
