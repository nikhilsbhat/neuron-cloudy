// Package database has a collections of structures which will commanly used across the subpackages of database.
package database

import (
	//"github.com/globalsign/mgo/bson"
	//"github.com/globalsign/mgo"
	"time"
)

var (
	// Db holds database type to which neuron makes a call.
	Db interface{}
)

// Storage holds the different types of the database to which neuron can talk to.
type Storage struct {
	Db interface{} `json:"Db,omitempty"`
	Fs string      `json:"Fs,omitempty"`
}

// DataDetail holds the information of database and collections, from which neuron collects information.
type DataDetail struct {
	Database   string
	Collection string
}

// UserData holds the details of the user stored as part of neuron.
type UserData struct {
	Id            int             `bson:"_id,omitempty" json:"id"`
	UserName      string          `json:"UserName" bson:"username"`
	Password      string          `json:"Password" bson:"password"`
	CloudProfiles []CloudProfiles `json:"CloudProfiles" bson:"cloudprofiles"`
}

// CloudProfiles holds the details of the profiles which is stored as part of user under neuron.
type CloudProfiles struct {
	Name           string    `json:"Name" bson:"name,omitempty"`
	Cloud          string    `json:"Cloud" bson:"cloud,omitempty"`
	KeyId          string    `json:"KeyId" bson:"keyid,omitempty"`
	SecretAccess   string    `json:"SecretAccess" bson:"secretaccess,omitempty"`
	ClientId       string    `json:"ClientID" bson:"clientid,omitempty"`
	SubscriptionId string    `json:"SubscriptionID" bson:"subscriptionid,omitempty"`
	TenantId       string    `json:"TenantID" bson:"tenantid,omitempty"`
	ClientSecret   string    `json:"ClientSecret" bson:"clientsecret,omitempty"`
	CreationTime   time.Time `json:"CreationTime" bson:"creationtime,omitempty"`
}

// CiData holds the information about the CI, to which neuron have a conversation.
type CiData struct {
	Id         int       `json:"id" bson:"_id,omitempty"`
	CiName     string    `json:"CiName" bson:"ciname"`
	CiURL      string    `json:"CiURL" bson:"ciurl"`
	CiUsername string    `json:"CiUsername" bson:"ciusername"`
	CiPassword string    `json:"CiPassword" bson:"cipassword"`
	Timestamp  time.Time `json:"Timestamp" bson:"timestamp"`
}

// GetCloudAccess will help one in fetching the profiles when he asks for.
type GetCloudAccess struct {
	ProfileName string
	Cloud       string
}
