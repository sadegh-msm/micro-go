package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
)

var Context = context.Background()

func CreateClients(No int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "pass",
		DB:       No,
	})
	return rdb
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	if url == "urlshortner-service" {
		return false
	}

	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == "urlshortner-service" {
		return false
	}

	return true
}
