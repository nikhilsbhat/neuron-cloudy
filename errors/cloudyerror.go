package cloudyerror

import (
	"errors"
)

const (
	emptyStructError = "You provided empty struct to retrive the data, this is not acceptable"
	notValidSession  = "Did not get session to perform action, cannot proceed further"
)

// InvalidSession returns an error if the methods finds the session in not vaild or illegal.
func InvalidSession() error {
	return errors.New(notValidSession)
}

// EmptyStructError will be thrown if the method encounters an empty structs where it should not have.
func EmptyStructError() error {
	return errors.New(emptyStructError)
}
