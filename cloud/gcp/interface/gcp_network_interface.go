package neurongcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

// GetNetworkInput holds the required values to get the details of the requied network.
type GetNetworkInput struct {
	// ProjectID refers to the ID of the GCP project in which the selected resource exists.
	ProjectID string
	// NetworkID refers to the name/id of the network of which the information has to be retrieved.
	NetworkID string
	GcpClient
}

// GetNetwork helps in retriving the information of the selected network.
func (net *GetNetworkInput) GetNetwork() (*compute.Network, error) {

	if net.Client != nil {
		ctx := context.Background()
		computeService, err := compute.New(net.Client)
		if err != nil {
			return nil, err
		}
		resp, err := computeService.Networks.Get(net.ProjectID, net.NetworkID).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

// GetNetworks helps in retriving the information of all networks in the selected project.
func (net *GetNetworkInput) GetNetworks() ([]*compute.Network, error) {

	if net.Client != nil {
		ctx := context.Background()
		computeService, err := compute.New(net.Client)
		if err != nil {
			return nil, err
		}
		req := computeService.Networks.List(net.ProjectID)
		networks := make([]*compute.Network, 0)
		if err := req.Pages(ctx, func(page *compute.NetworkList) error {
			for _, network := range page.Items {
				networks = append(networks, network)
			}
			return nil
		}); err != nil {
			return nil, err
		}

		return networks, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}
