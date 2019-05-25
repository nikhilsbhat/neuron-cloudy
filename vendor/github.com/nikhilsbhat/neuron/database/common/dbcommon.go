// Package dbcommon will help user to fetch
package dbcommon

import (
	//"github.com/globalsign/mgo/bson"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/nikhilsbhat/neuron/database"
	"github.com/nikhilsbhat/neuron/database/fs"
	"github.com/nikhilsbhat/neuron/database/mongodb"
	err "github.com/nikhilsbhat/neuron/error"
)

// ConfigDb will help neuron to set the database, so that neuron can use it in further operations.
func ConfigDb(d database.Storage) (interface{}, error) {

	if d.Db != nil {
		switch (d.Db).(type) {
		case *mgo.Session:
			database.Db = d.Db
			return nil, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	if d.Fs != "" {
		database.Db = d.Fs
		return nil, nil
	}
	return nil, fmt.Errorf("Oops..!! an error occurred. We did not receive valid input to configure DB")
}

// StoreCIdata will help neuron store information with regards to CI tool on to the database configured.
func StoreCIdata(d database.DataDetail, data database.CiData) (interface{}, error) {

	if database.Db != nil {
		switch database.Db.(type) {
		case *mgo.Session:
			status, staterr := mongo.StoreCIdata(d, data)
			if staterr != nil {
				return nil, staterr
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

// GetCiData will hep in fething the data of CI stored in neuron.
func GetCiData(ci string, d ...database.DataDetail) (database.CiData, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, staterr := mongo.GetCiData(ci, d[0])
			if staterr != nil {
				return database.CiData{}, staterr
			}
			return status, nil
		case string:
			status, staterr := fs.GetCiData(ci)
			if staterr != nil {
				return database.CiData{}, staterr
			}
			return status, nil
		default:
			return database.CiData{}, err.UnknownDbType()
		}
	}
	return database.CiData{}, fmt.Errorf("Database is not configured, we are not supporting filesystem now")
}

// CreateUser helps in creating user for neuron so that created user has appripriate permission to carry out the operations.
func CreateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, staterr := mongo.CreateUser(d, data)
			if staterr != nil {
				return nil, staterr
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

// UpdateUser helps in updating details of user created in neuron. This comes handy if you would like to add cloud profiles to the existing user.
// This not the only use, but is just a use case.
func UpdateUser(session interface{}, d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, staterr := mongo.UpdateUser(d, data)
			if staterr != nil {
				return nil, staterr
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

// GetUserDetails helps in fetching the user details stored as part of neuron.
func GetUserDetails(d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, staterr := mongo.GetUserDetails(d, data)
			if staterr != nil {
				return nil, staterr
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

/*func GetUsersDetails(d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, staterr := mongo.GetUsersDetails(d, data)
			if staterr != nil {
				return nil, staterr
			}
			return status, nil
		case string:
			status, staterr := fs.GetUsersDetails(data, cred)
			if staterr != nil {
				return database.CloudProfiles{}, staterr
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}*/

// GetCloudCredentails helps in fetching access details of cloud stored under a particular cloud-profile of user.
func GetCloudCredentails(data database.UserData, cred database.GetCloudAccess, d ...database.DataDetail) (database.CloudProfiles, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, staterr := mongo.GetCloudCredentails(d[0], data, cred)
			if staterr != nil {
				return database.CloudProfiles{}, staterr
			}
			return status, nil
		case string:
			status, staterr := fs.GetCloudCredentails(data, cred)
			if staterr != nil {
				return database.CloudProfiles{}, staterr
			}
			return status, nil
		default:
			return database.CloudProfiles{}, err.UnknownDbType()
		}
	}

	return database.CloudProfiles{}, err.DbNotConfiguredError()
}
