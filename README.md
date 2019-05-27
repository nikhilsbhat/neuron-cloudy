<p align="center">
    <h3 align="center">Neuron-Cloudy</h3>
    <p align="center">A cloud agnistic SDK.</p>
    <p align="center">
        <a href="https://goreportcard.com/report/github.com/nikhilsbhat/neuron-cloudy"><img src="https://goreportcard.com/badge/github.com/nikhilsbhat/neuron-cloudy"></a>
        <a href="https://github.com/nikhilsbhat/neuron-cloudy/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-APACHE2-blue.svg"></a>
    </p>
</p>

Neuron-CLoudy is a golang sdk buit to interact with all the cloud with a single function call.
This is used by [neuron](https://github.com/nikhilsbhat/neuron), which has a cli and exposes api to
create resources accross various cloud. As of now cloudy can talk to AWS and AZURE(Very minimal resources). Soon more features will be added and extends to various clouds.

One who wants to use this, just include this as library and start using it your lines of code.
Cloudy supports creation/deletion/updation and fetching details of various cloud resources including:
vpc,server,loadbalancer etc.

## Documentation

### Installation:

```golang
go get -u github.com/nikhilsbhat/neuron-cloudy
```
```golang
import (
    "github.com/nikhilsbhat/neuron-cloudy"
)
```
