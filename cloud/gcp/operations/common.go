package gcp

import (
	"net/http"

	"github.com/nikhilsbhat/neuron-cloudy/cloud/session"
	"golang.org/x/oauth2/jwt"
)

func isBaseClient(client interface{}) bool {
	switch client.(type) {
	case *jwt.Config:
		return true
	case *http.Client:
		return false
	default:
		return false
	}
}

func getClientFromBase(baseClient interface{}, scope []string) *http.Client {
	client := new(http.Client)
	if isBaseClient(baseClient) {
		baseclient := (baseClient).(*jwt.Config)
		baseclient.Scopes = scope
		client = session.GetClient(baseclient)
	} else {
		client = (baseClient).(*http.Client)
	}
	return client
}
