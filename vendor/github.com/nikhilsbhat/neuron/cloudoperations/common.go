package commoncloud

// Cloud is the common structure which is called in all cloudoperations.
type Cloud struct {
	// Pass the cloud in which the resource has to be created. usage: "aws","azure" etc.
	Name string `json:"name"`

	// Along with cloud, pass region in which resource has to be created.
	Region string `json:"region"`

	// Passing the profile is important, because this will help in fetching the the credentials
	// of cloud stored along with user details.
	Profile string `json:"profile"`

	// Use this option if in case you need unfiltered output from cloud.
	GetRaw bool `json:"getraw"`
}
