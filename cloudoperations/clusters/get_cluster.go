package clusters

import (
	"fmt"
	"strings"

	//"github.com/aws/aws-sdk-go/aws/session"

	gcp "github.com/nikhilsbhat/neuron-cloudy/cloud/gcp/operations"
	common "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/common"
	support "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/support"
)

// GetClusters will collect all the required information of the images specified to it and send back the response.
func (clust *GetClusterInput) GetClusters() (ClusterResponse, error) {

	if status := support.DoesCloudSupports(strings.ToLower(clust.Cloud.Name)); status != true {
		return ClusterResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetCluster")
	}

	switch strings.ToLower(clust.Cloud.Name) {
	case "aws":
		return ClusterResponse{}, fmt.Errorf(common.DefaultAwsResponse)
	case "azure":
		return ClusterResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":

		getCluster := new(gcp.GetClusterInput)
		getCluster.ClusterName = clust.ClusterName
		getCluster.ProjectID = clust.ProjectID
		getCluster.Regions = clust.Regions
		resp, err := getCluster.GetClusters(clust.Cloud.Client)
		if err != nil {
			return ClusterResponse{}, err
		}
		return ClusterResponse{GCPResponse: resp}, nil

	case "openstack":
		return ClusterResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return ClusterResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetCluster")
	}
}

// GetCluster will fetch the details of all the kubernetes cluster in the specified account of region.
func (clust *GetClusterInput) GetCluster() (ClusterResponse, error) {

	switch strings.ToLower(clust.Cloud.Name) {
	case "aws":
		return ClusterResponse{}, fmt.Errorf(common.DefaultAwsResponse)
	case "azure":
		return ClusterResponse{}, fmt.Errorf(common.DefaultAzResponse)
	case "gcp":

		getCluster := new(gcp.GetClusterInput)
		getCluster.ClusterName = clust.ClusterName
		getCluster.ProjectID = clust.ProjectID
		getCluster.Regions = clust.Regions
		resp, err := getCluster.GetClusters(clust.Cloud.Client)
		if err != nil {
			return ClusterResponse{}, err
		}
		return ClusterResponse{GCPResponse: resp}, nil

	case "openstack":
		return ClusterResponse{}, fmt.Errorf(common.DefaultOpResponse)
	default:
		return ClusterResponse{}, fmt.Errorf(common.DefaultCloudResponse + "GetAllImage")
	}
}

// New returns the new instance of GetImagesInput with empty values.
func New() *GetClusterInput {
	net := &GetClusterInput{}
	return net
}
