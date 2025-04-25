package main

import (
	"errors"
	"net/http"
)

func GetCurrentUserId() string {
	return "user01"
}

// behind an SP proxy
func GetCurrentUserIdFromHttpHeader(r *http.Request, key string) (string, error) {
	userId := r.Header.Get(key)
	if userId == "" {
		return userId, errors.New("empty user id")
	}
	return userId, nil
}
