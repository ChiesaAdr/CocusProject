package config

import (
	"os"
	"strconv"
)

func GetScalable() bool {

	scalable, exists := os.LookupEnv("COCUS_SCALABLE")
	if !exists {
		return false
	}

	return scalable == "true"
}

func GetRedisLimit() int {

	limit, exists := os.LookupEnv("REDIS_LIMIT_ENTRIES")
	if !exists {
		//Default Limit
		return 20
	}
	v, e := strconv.Atoi(limit)
	if e != nil {
		//Default Limit
		return 20
	}
	return v
}

func GetRedisURL() string {

	url, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		return "localhost:6379"
	}
	return url
}

func GetRedisPassword() string {

	url, exists := os.LookupEnv("REDIS_PASSWORD")
	if !exists {
		return ""
	}

	return url
}

func GetAddressServer() string {

	addr, exists := os.LookupEnv("ADDRESS_SERVER")
	if !exists {
		return "127.0.0.1:0"
	}
	return addr
}
