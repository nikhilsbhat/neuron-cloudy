<p align="center">
   <img alt="Neuron-Cloudy" src="https://raw.githubusercontent.com/nikhilsbhat/neuron-cloudy/development/assets/img/logo.png" height="100" />
    <h3 align="center">Neuron-Cloudy</h3>
    <p align="center">A cloud agnistic SDK.</p>
    <p align="center">
        <a href="https://goreportcard.com/report/github.com/nikhilsbhat/neuron-cloudy"><img src="https://goreportcard.com/badge/github.com/nikhilsbhat/neuron-cloudy"></a>
        <a href="https://github.com/nikhilsbhat/neuron-cloudy/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-apache%20v2-blue.svg"></a>
        <a href="https://godoc.org/github.com/nikhilsbhat/neuron-cloudy"><img src="https://godoc.org/github.com/nikhilsbhat/neuron-cloudy?status.svg" alt="GoDoc"></a>
    </p>
</p>

## Introduction

Cloudy is a golang sdk built to interact with all the cloud with a single function call.
This is used by [neuron](https://github.com/nikhilsbhat/neuron), which has a cli and exposes common api to
create resources accross various cloud. As of now cloudy can talk to AWS, GCP and AZURE(Very minimal resources). Soon more features will be added and extends to various clouds.

One who wants to use this, just include this as library in your line of code and start coding.
Cloudy supports __creation/deletion/updation/retrieval__ of various cloud resources including:
vpc,server,loadbalancer,cluster etc.

## Prerequisites

* It expected that [go](https://golang.org/dl/) to be pre installed on the machine. Installing go can be found [here](https://golang.org/doc/install).
* Understanding of basics of cloud and its components.

### Installation:

Use `go get` to retrieve the SDK to add it to your `GOPATH` workspace, or
project's Go module dependencies.
```bash
go get github.com/nikhilsbhat/neuron-cloudy
```
To update the SDK use `go get -u` to retrieve the latest version of the SDK.
```bash
go get -u github.com/nikhilsbhat/neuron-cloudy
```
Import to use it in code.
```golang
import (
    "github.com/nikhilsbhat/neuron-cloudy"
)
```

## Documentation

This Cloudy rely on dependency injection for the better outcome. As part of this one has to initialize the client and pass on the client for __creation/deletion/updation/retrieval__ of the cloud resources. Below sample code explains them.

initializing content for _GCP_ client
```golang
session := sess.CreateGcpSessionInput {
    CredPath = "path/to/credentials.json"
}
```

initializing content for _AWS_ client
```golang
session := sess.CreateAwsSessionInput {
    KeyId = "KEY_ID_OF_AWS"
    AcessKey = "SECRET_ACCESS_KEY_OF_AWS"
}
```
initializing client for the inputs gathered in above step.
```golang
client, err := session.CreateSession()
if err != nil {
    log.Fatal(err)
}
```

The above initialized client should be passed on to all further calls to cloud.

```golang
client, err := session.CreateSession()
if err != nil {
    log.Fatal(err)
}

input := network.New()
input.Cloud.Client = client           // pass the client initialized in the previous step.
input.Cloud.Name = "AWS"              // select the cloud of your preference for resource creation.
input.Cloud.Region = "ap-south-1"     // select the region of your preference for of the cloud selected.
input.Cloud.GetRaw = true             // set this flag if you prefer unfiltered output, defaults to false.
input.NetworkID = []string{"VPC-ID"}  // ID of the network of which the information to be retrieved.

resp, err := input.GetNetworks()      // returns the details of the network selected.
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%v\n",resp)
```

### Go Modules

If you are using Go modules, your `go get` will default to the latest tagged
release version of the SDK. To get a specific release version of the SDK use
`@<tag>` in your `go get` command.

	go get github.com/nikhilsbhat/neuron-cloudy@v0.0.12

To get the latest SDK repository change use `@latest`.

	go get github.com/nikhilsbhat/neuron-cloudy@latest

## Limitations

Currently it supports basic resources, more resources to be added to cloudy for better usability.

## TODO

* [ ] Writing Tests
* [ ] API Reference Doc
* [ ] Examples
