package utils

import "net/http"

var NotRetryAbleStatusCode = map[int]bool{
	http.StatusBadRequest: true,
	http.StatusNotFound: true,
	http.StatusUnsupportedMediaType: true,
	http.StatusUnprocessableEntity: true,
}

func IsNotRetryAble(statusCode int) bool {
	return NotRetryAbleStatusCode[statusCode]
}
