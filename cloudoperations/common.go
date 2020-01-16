package commoncloud

// Cloud is the common structure which is called in all cloudoperations.
type Cloud struct {
	// Pass the cloud in which the resource has to be created. usage: "aws","azure" etc.
	Name string `json:"name"`
	// Along with cloud, pass region in which resource has to be created.
	Region string `json:"region"`
	// Profile is important, because this will help in fetching the the credentials
	// of cloud stored along with user details.
	Profile string `json:"profile"`
	// Use this option if in case you need unfiltered output from cloud.
	GetRaw bool `json:"getraw"`
	// Path to credential file to be specified under this.
	CredPath string `json:"credpath"`
	// Client for the appropriate cloud, without this one cannot interact with the various resource of neuron-cloudy.
	Client interface{}
}
