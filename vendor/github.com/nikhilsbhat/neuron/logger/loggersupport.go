package neuronlogger

import (
	"bytes"
	"encoding/json"
	err "github.com/nikhilsbhat/neuron/error"
	"io/ioutil"
	"os"
)

const (
	pritnDash       = "+++++++++++++++++++++++++++++++++++++++++++++++++++++++"
	loginitializing = "Logging is atmost important than anything other step hence preparing for it"
	printemptyline  = ""
)

func getlog() (string, error) {

	if _, direrr := os.Stat("/var/lib/neuron/neuron.json"); os.IsNotExist(direrr) {

		Info(pritnDash)
		Info(loginitializing)
		Info(printemptyline)
		return "/var/log/neuron", nil
	} else {

		configdata, conferr := getloglocation()
		if conferr != nil {
			switch conferr.(type) {
			case err.NoLogFound:
				return "/var/log/neuron", nil
			default:
				return "", conferr
			}
		}
		return configdata, nil
	}
}

func getloglocation() (string, error) {
	configfile, conferr := ioutil.ReadFile("/var/lib/neuron/neuron.json")
	if conferr != nil {
		return "", err.ReadFileError()
	}
	decoder := json.NewDecoder(bytes.NewReader([]byte(configfile)))

	var confdata map[string]interface{}
	if decoderr := decoder.Decode(&confdata); decoderr != nil {
		Error(err.JsonDecodeError())
		return "", err.InvalidConfig()
	}

	if confdata["loglocation"] != nil {
		return (confdata["loglocation"]).(string), nil
	}
	return "", err.LogNotFound()
}
