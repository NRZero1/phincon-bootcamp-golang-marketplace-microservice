package utils

import "net/http"

var SuccessStatusCode = map[int]bool{
	http.StatusOK: true,
}

func IsSuccess(statusCode int) bool {
	return SuccessStatusCode[statusCode]
}
