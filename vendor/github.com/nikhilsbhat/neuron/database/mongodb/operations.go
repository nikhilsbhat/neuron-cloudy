package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/nikhilsbhat/neuron/database"
	"strings"
)

// StoreCIdata helps in storing data with regards to ci into the database to which neuron talks to.
func StoreCIdata(d database.DataDetail, data database.CiData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"ciname": data.CiName, "ciurl": data.CiURL})
	resp := []bson.M{}
	qryerr := query.All(&resp)
	if qryerr != nil {
		return nil, qryerr
	}
	if len(resp) == 0 {
		inserr := c.Insert(data)
		if inserr != nil {
			return nil, inserr
		}
		return "Records Created Successfully, check previous page for the details", nil
	}
	return nil, fmt.Errorf("The details you enetered matches with existing records")
}

// GetCiData helps in fetching the details of ci stored in neuron.
func GetCiData(n string, d database.DataDetail) (database.CiData, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{})
	resp := []bson.M{}
	qryerr := query.All(&resp)
	if qryerr != nil {
		return database.CiData{}, qryerr
	}
	_ = resp
	return database.CiData{}, nil
}

// CreateUser helps in creating user under neuron.
func CreateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qryerr := query.All(&resp)
	if qryerr != nil {
		return nil, qryerr
	}
	if len(resp) == 0 {
		// making indexing dynamic
		indexquery := c.Find(bson.M{})
		indexresp := []bson.M{}
		indexqryerr := indexquery.Sort("_id").All(&indexresp)
		if indexqryerr != nil {
			return nil, indexqryerr
		}
		var insertvalue database.UserData
		if data.CloudProfiles != nil {
			cldprf := make([]database.CloudProfiles, 0)
			//cldprf = append(cldprf, data.CloudProfiles)
			insertvalue = database.UserData{Id: (indexresp[len(indexresp)-1]["_id"].(int)) + 1, UserName: data.UserName, Password: data.Password, CloudProfiles: cldprf}
		} else {
			insertvalue = database.UserData{Id: (indexresp[len(indexresp)-1]["_id"].(int)) + 1, UserName: data.UserName, Password: data.Password}
		}
		inserr := c.Insert(insertvalue)
		if inserr != nil {
			return nil, inserr
		}
		return "User details saved successfully", nil

	}
	if len(resp) > 1 {
	}
	return "We cannot take the data you entered, Because we found the data matchces your entries", nil
}

// UpdateUser helps in updating user stored under neuron.
func UpdateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qryerr := query.All(&resp)
	if qryerr != nil {
		return nil, qryerr
	}
	existing := bson.M{"username": data.UserName, "password": data.Password}

	var valappnd []interface{}
	for _, value := range resp {
		for k, v := range value {
			if k == "cloudprofiles" {
				for _, v1 := range v.([]interface{}) {
					valappnd = append(valappnd, v1)
				}
			}
		}
	}

	valappnd = append(valappnd, data.CloudProfiles)
	change := bson.M{"$set": bson.M{"cloudprofiles": valappnd}}
	_, uperr := c.Upsert(existing, change)
	if uperr != nil {
		return nil, uperr
	}
	return "User profile updated successfully", nil
}

// GetUserDetails helps in fetching details of the user stored under neuron under the database mongo.
func GetUserDetails(d database.DataDetail, data database.UserData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qryerr := query.Sort("_id").All(&resp)
	if qryerr != nil {
		return nil, qryerr
	}
	for _, value := range resp {
		return database.UserData{UserName: value["username"].(string), Password: value["password"].(string)}, nil
	}
	return nil, fmt.Errorf("Something went wrong while fetching user details")
}

// GetCloudCredentails helps in fetching access credentials of cloud stored under neuron.
func GetCloudCredentails(d database.DataDetail, data database.UserData, cred database.GetCloudAccess) (database.CloudProfiles, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qryerr := query.Sort("_id").All(&resp)

	if qryerr != nil {
		fmt.Println(qryerr)
	}

	for _, value := range resp {
		for k, v := range value {
			if k == "cloudprofiles" {
				for _, v1 := range v.([]interface{}) {
					if (v1.(bson.M)["name"].(string) == cred.ProfileName) && (strings.ToLower(v1.(bson.M)["cloud"].(string)) == cred.Cloud) {
						return database.CloudProfiles{Name: v1.(bson.M)["name"].(string), Cloud: v1.(bson.M)["cloud"].(string), KeyId: v1.(bson.M)["keyid"].(string), SecretAccess: v1.(bson.M)["secretaccess"].(string)}, nil
					}
				}
			}
		}
	}
	return database.CloudProfiles{}, fmt.Errorf("Unable to find cloud credentials for the profile enetered")
}
