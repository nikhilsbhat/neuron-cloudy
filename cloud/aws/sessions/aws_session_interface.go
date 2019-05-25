// Package sessionaws is helpful in establishing connection with aws.
// One will be able to use all other api's of aws from this SDK only after this package is initialized.
package sessionaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// CreateSessionInput will implement CreateAwsSession and hold values which will be used in establishing connection.
type CreateSessionInput struct {
	// Region is required to establish session with aws. Without default value session cannot be established.
	Region string

	// KeyId is one of the component of AWS credentials and is necessary and not optional.
	KeyId string

	// AcessKey is one more component of AWS credentials and is necessary and not optional.
	AcessKey string
}

// CreateAwsSession will actually establish connection to aws with the credentials passed to it.
// Make sure you pass the correct/working credentials.
func (con *CreateSessionInput) CreateAwsSession() *session.Session {

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentialsFromCreds(
			credentials.Value{
				AccessKeyID:     con.KeyId,
				SecretAccessKey: con.AcessKey,
			}),
		Region: aws.String(con.Region),
	}))
	return sess
}
