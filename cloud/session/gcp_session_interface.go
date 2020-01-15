package session

import (
	"fmt"
	"net/http"

	"github.com/nikhilsbhat/config/decode"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/compute/v1"
)

type gcpSVCred struct {
	Type                string `json:"type,omitempty"`
	ProjectID           string `json:"project_id,omitempty"`
	PrivateKeyID        string `json:"private_key_id,omitempty"`
	PrivateKey          string `json:"private_key,omitempty"`
	ClientEmail         string `json:"client_email,omitempty"`
	ClientID            string `json:"client_id,omitempty"`
	AuthURI             string `json:"auth_uri,omitempty"`
	TokenURI            string `json:"token_uri,omitempty"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientCertURL       string `json:"client_x509_cert_url,omitempty"`
}

// CreateGcpSessionInput holds the client of GCP which will be helpful to connect with Google Cloud.
type CreateGcpSessionInput struct {
	gcpsvcauth *gcpSVCred
	// ProjectID of the google cloud of which the resource has to be dealt with.
	ProjectID string
	// AuthScopes has to be passed based on the resource that was selected to dealt with.
	AuthScopes []string
	// CredPath of GCP credentails, the preferred way of initializing client for GCP.
	CredPath string
	// Zone of the resource.
	Zone string
	// RawJSON could be alternative to CredPath, if one prefers to inject credential directly this is the way
	// this contradicts CredPath, both cannot be used simultaneously.
	RawJSON []byte
	// Client is the actual client of the Google Cloud.
	// Client *http.Client
}

// Session is an interface that implements session for all clouds.
type Session interface {
	CreateSession()
}

// CreateSession returns the valid gcp client which has permissions to perform operation in GCP.
func (auth *CreateGcpSessionInput) CreateSession() (interface{}, error) {
	if (len(auth.CredPath) == 0) && (auth.RawJSON == nil) {
		client, err := auth.getDefalutGCPClient()
		if err != nil {
			return nil, fmt.Errorf("Unable to initialize the default client for GCP")
		}
		return client, nil
	}

	if err := auth.getGCPCred(); err != nil {
		return nil, err
	}
	if client := auth.getCustomGCPClient(); client != nil {
		return client, nil
	}
	return nil, fmt.Errorf("Unable to initialize client for GCP")
}

func (auth *CreateGcpSessionInput) getDefalutGCPClient() (*http.Client, error) {

	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (auth *CreateGcpSessionInput) getCustomGCPClient() *jwt.Config {

	conf := &jwt.Config{
		Email:      auth.gcpsvcauth.ClientEmail,
		PrivateKey: []byte(auth.gcpsvcauth.PrivateKey),
		Scopes:     auth.AuthScopes,
		TokenURL:   auth.gcpsvcauth.TokenURI,
		Subject:    auth.gcpsvcauth.ClientEmail,
	}
	return conf
}

// GetClient returns the fully working cutom client
func GetClient(conf *jwt.Config) *http.Client {
	return conf.Client(oauth2.NoContext)
}

func (auth *CreateGcpSessionInput) getGCPCred() error {

	if (len(auth.CredPath) != 0) && (auth.RawJSON != nil) {
		return fmt.Errorf("Cannot use both RawJSON and Path to credential file together")
	}

	if len(auth.CredPath) != 0 {
		jsonCont, err := decode.ReadFile(auth.CredPath)
		if err != nil {
			return err
		}
		auth.RawJSON = jsonCont
	}

	if err := auth.decodeGcpSVCred(); err != nil {
		return err
	}
	return nil
}

func (auth *CreateGcpSessionInput) decodeGcpSVCred() error {

	jsonAuth := new(gcpSVCred)
	if decodneuerr := decode.JsonDecode(auth.RawJSON, &jsonAuth); decodneuerr != nil {
		return decodneuerr
	}
	auth.gcpsvcauth = jsonAuth
	return nil
}
