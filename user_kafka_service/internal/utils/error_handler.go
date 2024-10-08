package utils

import (
	"errors"
)

var (
	ErrCreateRequest = errors.New("error when trying to create request")
	ErrClientSend = errors.New("error when trying to send http request")
	ErrRepoCreateTrx = errors.New("error when trying to begin transaction")
	ErrPreparedStmt = errors.New("error when trying to create prepared statement")
	ErrDbExec = errors.New("error when trying to execute query")
	ErrTrxCommit = errors.New("error when trying to commit transaction")
	ErrNoSqlRows = errors.New("no result found")
	ErrErrScan = errors.New("error when trying to scan query result")
	ErrJsonDecode = errors.New("error when trying to decode json")
	ErrKafkaConsume = errors.New("error when trying to read message")
	ErrKafkaReaderClose = errors.New("error when trying to close kafka consumer")
	ErrKafkaProducer = errors.New("failed to produce order message")
	ErrMarshal = errors.New("failed to encode")
	ErrHttpRequest = errors.New("failed to create request to user service")
	ErrHttpNotRetryAble = errors.New("http request returned an error with not retryable error and need manual edit and retry")
	ErrMaxRetryReached = errors.New("system tried to retry the request but max retry reached")
	ErrUnexpectedResponse = errors.New("unexpected response, exit the loop")
)
