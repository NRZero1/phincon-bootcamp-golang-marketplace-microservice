package utils

import (
	"errors"
	"fmt"
)

var (
	ErrFindUserById error
	ErrFindUserByUsername error
	ErrHash = errors.New("error hashing password")
	ErrDecode = errors.New("error decoding JSON")
	ErrValidation = errors.New("failed to save user because didn't pass the validation")
	ErrRequestTimedOut = errors.New("request timed out")
	ErrPathVar = errors.New("error when converting string to int for path variable")
	ErrWrongPass = errors.New("wrong password")
	ErrCreateSignature = errors.New("error creating signature")
)

func NewErrFindUserById(id int) (error) {
	ErrFindUserById = fmt.Errorf("no user found with ID %d", id)
	return ErrFindUserById
}

func NewErrFindByUsername(username string) (error) {
	ErrFindUserByUsername = fmt.Errorf("no user found with username %s", username)
	return ErrFindUserByUsername
}
