package commonoperations

const (
	// DefaultAwsResponse holds message when person choose aws as their cloud.
	DefaultAwsResponse = "We have not reached to aws yet on this resource"
	// DefaultOpResponse holds message when person choose openstack as their cloud.
	DefaultOpResponse = "We have not reached to openstack yet on this resource"
	// DefaultAzResponse holds message when person choose azure as their cloud.
	DefaultAzResponse = "We have not reached to azure yet on this resource"
	// DefaultGcpResponse holds message when person choose google as the cloud.
	DefaultGcpResponse = "We have not reached to google cloud yet on this resource"
	// DefaultCloudResponse holds default message.
	DefaultCloudResponse = "I feel we are lost in performing the action, guess you have entered wrong cloud. The action was: "
	// BetaResponse helps in constructing response for beta resources.
	BetaResponse = "%s is in Beta and supports very minimal support."
	// AlphaResponse helps in constructing response for Alpha resources.
	AlphaResponse = "%s is in Alpha and supports very minimal support."
	// InvalidClientResponse has the valid response if the client is invalid.
	InvalidClientResponse = "The client passed for invoking %s api is invalid"
)
