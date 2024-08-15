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
	ErrErrScan = errors.New("error when trying to scan query result to struct")
	ErrJsonDecode = errors.New("error when trying to decode json")
	ErrKafkaConsume = errors.New("error when trying to read message")
	ErrKafkaReaderClose = errors.New("error when trying to close kafka consumer")
)