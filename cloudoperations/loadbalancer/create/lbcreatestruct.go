package createloadbalancer

import (
	cmn "github.com/nikhilsbhat/neuron-cloudy/cloudoperations"
)

// LbCreateInput takes all the inputs required by CreateLoadBalancer.
// These parameters vary based on cloud choosed.
type LbCreateInput struct {
	Name              string   `json:"name"`
	VpcId             string   `json:"vpcid"`
	SubnetIds         []string `json:"subnetids"`
	AvailabilityZones []string `json:"availabilityzones"`
	SecurityGroupIds  []string `json:"securitygroupids"`
	Scheme            string   `json:"scheme"`
	Type              string   `json:"type"` //required only if the LB protocol is HTTPS else can be initiazed with dummy value
	SslCert           string   `json:"sslcert"`
	SslPolicy         string   `json:"sslpolicy"`
	LbPort            int64    `json:"lbport"`
	InstPort          int64    `json:"instport"`
	Lbproto           string   `json:"lbproto"` //required ex: HTTPS, HTTP
	Instproto         string   `json:"instproto"`
	HttpCode          string   `json:"httpcode"`
	HealthPath        string   `json:"healthpath"`
	IpAddressType     string   `json:"ipaddresstype"`
	Cloud             cmn.Cloud
}

//Nothing much from this file. This file contains only the structs for loadbalance/create
