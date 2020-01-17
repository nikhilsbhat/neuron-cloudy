// Package cloudy helps user to connect various cloud and provision/delete/updated/retrive the resource details cloud agnostically.
//
// The Cloudy SDK is bilt on Go and has APIs and utilities that one can use to
// build Go applications that uses various cloud services which includes AWS,GCP,AZURE
// Majority of them being in beta.
//
// This SDK helps in resolving the complexity for a developer who would end up invoking
// multiple SDKs of appropriate cloud. For this, one has to invoke single api to provision/delete/updated/retrive
// with cloud as a parameter to it.
//
// This SDK cloudy also eliminates the burden of developers remembering the name various services
// across the cloud inspite of the fact that they will endup consuming similar services.
// Ex: The network would be addressed as VPC in aws where are it is vnet in AZURE.
//
// To consume the APIs of cloudy one has to follow two steps.
// First initilize the client for the appropriate cloud,
// then invoke an api by passing client while provisioning/deletion/updation/retrival of the cloud services.
//
// // The doc https://godoc.org/github.com/nikhilsbhat/neuron-cloudy/cloud/session would help in initializing the client for cloudy:
//  Ex: GCP
//  session := sess.CreateGcpSessionInput {
//  		CredPath = "path/to/credentials.json"
//  }
//
//  Ex: AWS
//  session := sess.CreateAwsSessionInput {
//  		KeyId = "KEY_ID_OF_AWS"
//  		AcessKey = "SECRET_ACCESS_KEY_OF_AWS"
//  }
//
// Failing to pass the requied parameters while initializing client creates the default session fetching it from environment varibale.
// More info on what default session can found at:
// [GCP]https://godoc.org/golang.org/x/oauth2/google#DefaultClient
// [AWS](https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/#NewEnvCredentials)
//
package cloudy
