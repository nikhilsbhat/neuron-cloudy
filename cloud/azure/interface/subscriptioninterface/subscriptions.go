package azuresubscription

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/go-autorest/autorest"
	"neuron/cloud/azure/access"
)

var (
	token, _, _ = auth.GetServicePrincipalToken()
	ctx         = context.Background()
)

type SubcriptionIn struct {
	Subscription string
}

func getSubscriptionClient() subscription.SubscriptionsClient {
	subscriptionClient := subscription.NewSubscriptionsClient()
	subscriptionClient.Authorizer = autorest.NewBearerAuthorizer(token)

	return subscriptionClient
}

// This function will get the subcription from subscriptions that are associated with a tenant in azure account.

func (s SubcriptionIn) GetSubscription() (sub subscription.Model, err error) {
	subscriptionClient := getSubscriptionClient()
	future, err := subscriptionClient.Get(
		ctx,
		s.Subscription,
	)

	if err != nil {
		return sub, fmt.Errorf("cannot get subscription: %v", err)
	}

	return future, err
}

// This function will list the subcriptions that are associated with a tenant in azure account.

func ListSubscription() (sub []subscription.Model, err error) {
	subscriptionClient := getSubscriptionClient()
	future, err := subscriptionClient.List(
		ctx,
	)

	if err != nil {
		return sub, fmt.Errorf("cannot list subscriptions: %v", err)
	}

	return future.Values(), err
}
