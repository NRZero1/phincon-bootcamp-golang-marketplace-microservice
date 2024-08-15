package utils

import (
	"errors"
	"fmt"
)

var (
	ErrKafkaConsume error
	ErrKafkaReaderClose error
)

func NewErrKafkaConsume(message string) error {
	errMessage := fmt.Sprintf("Failed to read message from kafka with error message: %v", message)

	ErrKafkaConsume = errors.New(errMessage)
	return ErrKafkaConsume
}

func NewErrKafkaReaderClose(message string) error {
	errMessage := fmt.Sprintf("Failed to close Kafka reader: %s", message)

	ErrKafkaReaderClose = errors.New(errMessage)

	return ErrKafkaReaderClose
}