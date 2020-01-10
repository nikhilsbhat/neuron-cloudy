package commonoperations

import (
	db "github.com/nikhilsbhat/neuron/database"
	dbcommon "github.com/nikhilsbhat/neuron/database/common"
)

const (
	//
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
)

// GetCredentialsInput holds the information of profile and cloud that has to be fetched from database.
type GetCredentialsInput struct {
	Profile string
	Cloud   string
}

// GetCredentials helps in fetching the of the credentials of the specified user along with the cloud details asked for.
func GetCredentials(gcred *GetCredentialsInput) (db.CloudProfiles, error) {

	//fetchinig credentials from loged-in user to establish the connection with appropriate cloud.
	creds, crderr := dbcommon.GetCloudCredentails(
		db.UserData{UserName: "nikhibt434@gmail", Password: "42bhat24"},
		db.GetCloudAccess{ProfileName: gcred.Profile, Cloud: gcred.Cloud},
		db.DataDetail{"neuron", "users"},
	)
	if crderr != nil {
		return db.CloudProfiles{}, crderr
	}

	return creds, nil
}
