package neurongcp

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/container/v1"
)

// GetClusterInput holds the required values to fetch the cluster details
type GetClusterInput struct {
	// ResourceUrl refers to link of the resource of which the information has to be fetched
	ResourceURL string
	GcpClient
}

// GcpClient is a established session of GCP with appropriate AUTH scopes.
type GcpClient struct {
	// Client holds the session of GCP which has access.
	Client *http.Client
}

// ListClusters lists all GKE cluster and its details across the regions specified.
func (c *GetClusterInput) ListClusters() ([]*container.Cluster, error) {

	if c.Client != nil {
		ctx := context.Background()
		containerService, err := container.New(c.Client)
		if err != nil {
			return nil, err
		}
		resp, err := containerService.Projects.Locations.Clusters.List(c.ResourceURL).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		return resp.Clusters, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

// GetCluster fetches the details of the GKE cluster specified.
func (c *GetClusterInput) GetCluster() (*container.Cluster, error) {

	if c.Client != nil {
		ctx := context.Background()
		containerService, err := container.New(c.Client)
		if err != nil {
			return nil, err
		}
		resp, err := containerService.Projects.Locations.Clusters.Get(c.ResourceURL).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}
