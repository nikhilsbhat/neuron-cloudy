// This package is used to initialize the credentials from a credentials file (~/.azure/credentials).
// file shoud be in Json and the credentials will be initialized based on the username passed.

package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"io/ioutil"
	"log"
	"os/user"
)

type Credentials struct {
	Profile        string
	ClientID       string
	SubscriptionID string
	TenantID       string
	ClientSecret   string
}

var (
	result Credentials
)

func init() {

	user, _ := user.Current()
	file := user.HomeDir + "/.azure/credentials"

	plan, _ := ioutil.ReadFile(file) // filename is the JSON file to read

	var data []Credentials
	err := json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Errorf("Cannot unmarshal the json ", err)
	}
	for _, t := range data {

		if t.Profile == "ranjith" { // for now hardcoded as 'ranjith'
			result = t
			break
		} else if t.Profile != "ranjith" {
			continue
		}
	}
	if (Credentials{}) == result {
		fmt.Println("I have no credentials with the user you passed")
	}
}

func GetServicePrincipalToken() (adal.OAuthTokenProvider, error, string) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, result.TenantID)
	code, err := adal.NewServicePrincipalToken(
		*oauthConfig,
		result.ClientID,
		result.ClientSecret,
		azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("%s: %v\n", "failed to initiate device auth", err)
	}

	return code, err, result.SubscriptionID
}
