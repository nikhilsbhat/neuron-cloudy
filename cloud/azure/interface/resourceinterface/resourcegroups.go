package azurenetwork

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"neuron/cloud/azure/access"
)

var (
	token, _, subscription = auth.GetServicePrincipalToken()
	ctx                    = context.Background()
)

type GroupsIn struct {
	ResourceGroup string
	Location      string `json:"location,omitempty"`
}

func getGroupsClient() resources.GroupsClient {
	groupsClient := resources.NewGroupsClient(subscription)
	groupsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return groupsClient
}

// Creates a new resource group

func (g GroupsIn) CreateResourceGroup() (resources.Group, error) {
	groupsClient := getGroupsClient()
	fmt.Printf("\n creating resource group '%s' on location: %v", g.ResourceGroup, g.Location)
	return groupsClient.CreateOrUpdate(
		ctx,
		g.ResourceGroup,
		resources.Group{
			Location: to.StringPtr(g.Location),
		})
}
