package utils

import "net/http"

var RetryAbleStatusCode = map[int]bool{
	http.StatusRequestTimeout: true,
	http.StatusTooManyRequests: true,
	http.StatusInternalServerError: true,
	http.StatusBadGateway: true,
	http.StatusServiceUnavailable: true,
	http.StatusGatewayTimeout: true,
}

func IsRetryAbleStatusCode(statusCode int) bool {
    return RetryAbleStatusCode[statusCode]
}
