// Package session is helpful in establishing connection with aws.
// One will be able to use all other api's of aws from this SDK only after this package is initialized.
package session

import (
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nikhilsbhat/config/decode"
)

// CreateAwsSessionInput will implement CreateAwsSession and hold values which will be used in establishing connection.
type CreateAwsSessionInput struct {
	// Region is required to establish session with aws. Without default value session cannot be established.
	Region string
	// KeyId is one of the component of AWS credentials and is necessary and not optional.
	KeyId string `json:"access_key_id,omitempty"`
	// AcessKey is one more component of AWS credentials and is necessary and not optional.
	AcessKey string `json:"secret_access_key,omitempty"`
	// CredPath is path to the credentials file which is used to connect to AWS.
	CredPath string
	// Profile name of the credentail, this will go good when CredPath is used to initialize the client.
	Profile string `json:"profile,omitempty"`
	// RawJSON can take the raw input of the cloud credential passed (this contradicts CredPath, both cannot be used simultaneously).
	RawJSON []byte
	// CustomFile let user choose the json file for credentials
	CustomFile bool
}

type awsSVCred struct {
	// KeyId is one of the component of AWS credentials and is necessary and not optional.
	KeyId string
	// AcessKey is one more component of AWS credentials and is necessary and not optional.
	AcessKey string
}

// CreateSession will actually establish connection to aws with the credentials passed to it.
// Make sure you pass the correct/working credentials.
func (auth *CreateAwsSessionInput) CreateSession() (*session.Session, error) {

	if reflect.DeepEqual(auth, CreateAwsSessionInput{}) {
		return nil, fmt.Errorf("Cannot create session for AWS with empty input")
	}

	if (len(auth.KeyId) != 0) && (len(auth.AcessKey) != 0) {
		return auth.getCustomAWSClient(), nil
	}

	if len(auth.CredPath) != 0 {
		cred, err := auth.getAWSClient()
		if err != nil {
			return nil, err
		}
		return cred, nil
	}

	return auth.getDefalutAWSClient(), nil
}

// getDefalutAWSClient gets the client/session form the environment variable set (it is the default way of initialize the session in AWS).
func (auth *CreateAwsSessionInput) getDefalutAWSClient() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Region:      aws.String(auth.Region),
	}))
}

// getCustomAWSClient gets the client/session form the inputs passed to neuron (one of the way to initialize the session in AWS).
func (auth *CreateAwsSessionInput) getCustomAWSClient() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentialsFromCreds(
			credentials.Value{
				AccessKeyID:     auth.KeyId,
				SecretAccessKey: auth.AcessKey,
			}),
		Region: aws.String(auth.Region),
	}))
}

// getAWSClient gets the client/session form the file fed to cloudy (one of the way to initialize the session in AWS).
func (auth *CreateAwsSessionInput) getAWSClient() (*session.Session, error) {

	if len(auth.CredPath) != 0 {
		return nil, fmt.Errorf("CredPath cannot be empty")
	}

	if stat := auth.awsCredExists(); stat == false {
		return nil, fmt.Errorf("CredPath seems faulty, there are no such credential file")
	}

	return session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials(auth.CredPath, auth.Profile),
		Region:      aws.String(auth.Region),
	})), nil
}

// getCustomAWSClientFile gets the custom client/session form the custom file fed to cloudy.
// This is for future reference currently this does not have role to play.
func (auth *CreateAwsSessionInput) getCustomFileAWSClient() (*session.Session, error) {

	if stat := auth.awsCredExists(); stat == false {
		return nil, fmt.Errorf("CredPath seems faulty, there are no such credential file")
	}

	jsonCont, err := decode.ReadFile(auth.CredPath)
	if err != nil {
		return nil, err
	}
	auth.RawJSON = jsonCont

	cred, err := auth.decodeAwsSVCred()
	if err != nil {
		return nil, err
	}
	auth.KeyId = cred.KeyId
	auth.AcessKey = cred.AcessKey
	return auth.getCustomAWSClient(), nil
}

func (auth *CreateAwsSessionInput) decodeAwsSVCred() (*CreateAwsSessionInput, error) {

	jsonAuth := new(CreateAwsSessionInput)
	if decodneuerr := decode.JsonDecode(auth.RawJSON, &jsonAuth); decodneuerr != nil {
		return nil, decodneuerr
	}

	return jsonAuth, nil
}

func (auth *CreateAwsSessionInput) awsCredExists() bool {
	if _, err := os.Stat(auth.CredPath); os.IsNotExist(err) {
		return false
	}
	return true
}
