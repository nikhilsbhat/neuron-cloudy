package cloudyerror

import (
	"errors"
)

const (
	emptyStructError = "You provided empty struct to retrive the data, this is not acceptable"
	notValidSession  = "Did not get session to perform action, cannot proceed further"
	imageNotFound    = "Couldn't find the image with the Id passed, and hence cannot proceed further with the action specified."
	instanceNotFound = "Couldn't find server equivalent to the Id which you entered. Without server I cannot capture image"
)

// InvalidSession returns an error if the methods finds the session in not vaild or illegal.
func InvalidSession() error {
	return errors.New(notValidSession)
}

// EmptyStructError will be thrown if the method encounters an empty structs where it should not have.
func EmptyStructError() error {
	return errors.New(emptyStructError)
}

// ImageNotFound will be thrown if entered image does not exists
func ImageNotFound() error {
	return errors.New(imageNotFound)
}

// ServerNotFound will be thrown if the entered server does not exists
func ServerNotFound() error {
	return errors.New(instanceNotFound)
}
