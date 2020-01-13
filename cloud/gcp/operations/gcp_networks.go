package gcp

import (
	"fmt"
	"time"

	neurongcp "github.com/nikhilsbhat/neuron-cloudy/cloud/gcp/interface"
	"google.golang.org/api/compute/v1"
)

// GetNetworkInput holds the required values to fetch the net details
type GetNetworkInput struct {
	neurongcp.GetNetworkInput
	// GetRaw makes sure that function returns unfiltered response if it is set.
	GetRaw bool
	CredMode
}

// NetworkResponse contains filtered/unfiltered response from GCP on cluster details.
type NetworkResponse struct {
	// Name refers to the name of the network of which the information is retrived.
	Name string
	// IPv4Range refers to the range of IP address assogned to the particluar network.
	IPv4Range string
	// ID meaning the unique ID assgned to the network.
	ID uint64
	// Subnetworks holds the list of subnetworks present in the retrived network.
	Subnetworks []string
	// AutoCreateSubnetworks defines whether to create subnets automatically in a particular network.
	AutoCreateSubnetworks bool
	// Duration let to know the time since the network been created (The life of the network).
	Duration string
	// Peerings holds the configuaration details of the peered networks.
	Peerings []*compute.NetworkPeering
	// RoutingConfig holds the information of router attached to the network.
	RoutingConfig *compute.NetworkRoutingConfig
	// SelfLink holds the link which refers to a particluar network of which the information retrived.
	SelfLink string
	// GetNetworkRaw contains unfiltered response from GCP on network
	GetNetworkRaw []*compute.Network
}

// GetNetworks gets the list of kube network along with its other details
// This is region specific, mind the region you pass to it.
func (net *GetNetworkInput) GetNetworks(client interface{}) ([]NetworkResponse, error) {

	if len(net.ProjectID) == 0 {
		return nil, fmt.Errorf("Project ID cannot be nil")
	}

	// Initialization of gcp client
	sess := getClientFromBase(client, []string{compute.CloudPlatformScope})
	input := new(neurongcp.GetNetworkInput)
	input.ProjectID = net.ProjectID
	input.Client = sess
	networks, err := input.GetNetworks()
	if err != nil {
		return nil, err
	}

	response := make([]NetworkResponse, 0)
	if net.GetRaw == true {
		return append(response, NetworkResponse{GetNetworkRaw: networks}), nil
	}

	for _, network := range networks {
		netw := new(NetworkResponse)
		netw.Name = network.Name
		netw.ID = network.Id
		netw.IPv4Range = network.IPv4Range
		netw.AutoCreateSubnetworks = network.AutoCreateSubnetworks
		netw.Subnetworks = network.Subnetworks
		netw.Peerings = network.Peerings
		netw.RoutingConfig = network.RoutingConfig
		netw.SelfLink = network.SelfLink
		netw.Duration = network.CreationTimestamp

		// duration, err := getElapsedTime(network.CreationTimestamp)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to fetch the duration for the network")
		// 	response = append(response, *netw)
		// 	return response, err
		// }
		// netw.Duration = duration
	}
	return response, nil
}

// GetNetwork gets the list of kube network along with its other details
// This is region specific, mind the region you pass to it.
func (net GetNetworkInput) GetNetwork(client interface{}) ([]NetworkResponse, error) {

	if len(net.ProjectID) == 0 {
		return nil, fmt.Errorf("Project ID cannot be nil")
	}

	// Initialization of gcp client
	sess := getClientFromBase(client, []string{compute.CloudPlatformScope})
	// Returns list of networks created in the project if networkID is missed while invoking the function.
	if len(net.NetworkID) == 0 {
		network, err := net.GetNetworks(sess)
		if err != nil {
			return nil, err
		}
		return network, nil
	}

	networks := make([]*compute.Network, 0)
	input := new(neurongcp.GetNetworkInput)
	input.ProjectID = net.ProjectID
	input.NetworkID = net.NetworkID
	input.Client = sess
	network, err := input.GetNetwork()
	if err != nil {
		return nil, err
	}
	networks = append(networks, network)

	response := make([]NetworkResponse, 0)
	if net.GetRaw == true {
		return append(response, NetworkResponse{GetNetworkRaw: networks}), nil
	}

	for _, network := range networks {
		netw := new(NetworkResponse)
		netw.Name = network.Name
		netw.ID = network.Id
		netw.IPv4Range = network.IPv4Range
		netw.AutoCreateSubnetworks = network.AutoCreateSubnetworks
		netw.Subnetworks = network.Subnetworks
		netw.Peerings = network.Peerings
		netw.RoutingConfig = network.RoutingConfig
		netw.SelfLink = network.SelfLink
		netw.Duration = network.CreationTimestamp

		// duration, err := getElapsedTime(network.CreationTimestamp)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to fetch the duration for the network")
		// 	response = append(response, *netw)
		// 	return response, err
		// }
		// netw.Duration = duration
	}

	return response, nil
}

func getElapsedTime(duration string) (time.Duration, error) {
	Timelayout := "2019-01-12T09:20:21.277-08:00"
	creationTime, err := time.Parse(Timelayout, duration)
	if err != nil {
		return *new(time.Duration), err
	}

	currentTime := time.Now()
	elapsedTime := currentTime.Sub(creationTime)
	return elapsedTime, nil
}
