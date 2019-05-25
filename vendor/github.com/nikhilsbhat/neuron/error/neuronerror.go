package neuronerror

import (
	"errors"
	"fmt"
)

const (
	uiConfigNotFound    = "You mentioned directory of UI in config but couldn't find one in the path"
	configNotValid      = "Oops..! An error occurred while reading config file you passed, the configuration file provided is not vaild. Please provide a valid json"
	configNotFound      = "I couldn't find any config file by namr which you mentioned"
	configDataBase      = "Oops..! An error occurred while configuring databse with neuron, check application log (path/to/log/neuronapp.log) for more"
	dbNotConfigured     = "Oops..! the action cannot be completed as database/filesystem to data is not configured"
	decodeJsonError     = "Oops..! An error occurred while decoding the json to datastructure"
	dbSessionError      = "Oops..! An error occurred while establishig session to database"
	dirCreateError      = "Oops..! An error occurred while creating directory"
	errInvalidParent    = "Invalid Parent"
	logpathNotFound     = "Could not find logpath in the configuration you passed"
	uilogCreationError  = "Oops..!! An error occurred while creating a UI logfile"
	uiLogOpenError      = "Oops..!! An error occurred while opening UI log file"
	logCreationError    = "Oops..!! An error occurred while creating a neuron logfile"
	logOpenError        = "Oops..!! An error occurred while opening neuron log file"
	logDirError         = "Oops..!! An error occurred while creating a neuron log directory"
	logSetupError       = "Oops..!! we ran into an error while setting up log, with logs we cannot proceed"
	neuronStartError    = "Oops..!! We have encountered error while starting neuron App, check if the specified port is free"
	unableToStart       = "We faced series of error because of which we are unable to start application"
	readFilError        = "Oops..!! An error occurred while reading a file, either file is not valid no sufficient permission to open it"
	testMessage         = "+++++++++++++++This is test error+++++++++++++++++"
	unknownDatabaseType = "Oops..!! we received a session of unkown database type"
	usersFileNotFound   = "We were unable to locate the file containing users data"
	invalidUsersFile    = "Oops..!! Users file is not in a state of decoding its data. Either it contains data in wrong format or file is corrupted"
	emptyStructError    = "You provided empty struct to retrive the data, this is not acceptable"
	notValidSession     = "Did not get session to perform action, cannot proceed further"
	instanceNotFound    = "Couldn't find server equivalent to the Id which you entered. Without server I cannot capture image"
	imageNotFound       = "Couldn't find the image with the Id passed, and hence cannot proceed further with the action specified."
	fileNotFound        = "We were unable to find the file you specified"
	clierror            = "Oops..!! error occurred while initializing neuron cli"
	cliproceederror     = "Without basic configuration we cannot proceed further"
	clinotstarting      = "Neuron CLI was not initialized properly, please check and call it again"
)

type error interface {
	Error() string
}

type NoLogFound struct {
	Message string
}

func (e NoLogFound) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

func InvalidConfig() error {
	return errors.New(configNotValid)
}

func ConfigNotfound() error {
	return errors.New(configNotFound)
}

func UiNotFound() error {
	return errors.New(uiConfigNotFound)
}

func JsonDecodeError() error {
	return errors.New(decodeJsonError)
}

func ConfigDbError() error {
	return errors.New(configDataBase)
}

func DbSessionError() error {
	return errors.New(dbSessionError)
}

func UnknownDbType() error {
	return errors.New(unknownDatabaseType)
}

func DbNotConfiguredError() error {
	return errors.New(dbNotConfigured)
}

func DirCreateError() error {
	return errors.New(dirCreateError)
}

func UiLogCreationError() error {
	return errors.New(uilogCreationError)
}

func UiLogOpenError() error {
	return errors.New(uiLogOpenError)
}

func LogCreationError() error {
	return errors.New(logCreationError)
}

func LogOpenError() error {
	return errors.New(logOpenError)
}

func LogDirError() error {
	return errors.New(logDirError)
}

func SetupLogError() error {
	return errors.New(logSetupError)
}

func StartNeuronError() error {
	return errors.New(neuronStartError)
}

func FailStartError() error {
	return errors.New(unableToStart)
}

func LogNotFound() error {
	return NoLogFound{logpathNotFound}
}

func ReadFileError() error {
	return errors.New(readFilError)
}

func UsersNotFound() error {
	return errors.New(usersFileNotFound)
}

func InvalidUsersFile() error {
	return errors.New(invalidUsersFile)
}

func InvalidSession() error {
	return errors.New(notValidSession)
}

func EmptyStructError() error {
	return errors.New(emptyStructError)
}

func ServerNotFound() error {
	return errors.New(instanceNotFound)
}

func ImageNotFound() error {
	return errors.New(imageNotFound)
}

func InvalidCiDataFile() error {
	return errors.New(invalidUsersFile)
}

func NoFileFoundError() error {
	return errors.New(fileNotFound)
}

func UninitializedCli() error {
	return errors.New(clierror)
}

func CliFailure() error {
	return errors.New(cliproceederror)
}

func CliNoStart() error {
	return errors.New(clinotstarting)
}
