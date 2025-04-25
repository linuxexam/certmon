package main

import (
	"errors"
	"net/http"
)

func GetCurrentUserId(r *http.Request) string {
	userId, err := GetCurrentUserIdFromHttpHeader(r, "UTSSO-Proxy-utorid")
	if err != nil {
		return ""
	}
	return userId
}

// behind an SP proxy
func GetCurrentUserIdFromHttpHeader(r *http.Request, key string) (string, error) {
	userId := r.Header.Get(key)
	if userId == "" {
		return userId, errors.New("empty user id")
	}
	return userId, nil
}
