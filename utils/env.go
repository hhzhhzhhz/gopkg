package utils

import (
	"os"
	"strconv"
)

func Env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

func boolEnv(key string, def bool) bool {
	if x := os.Getenv(key); x != "" {
		if v, err := strconv.ParseBool(x); err == nil {
			return v
		}
	}
	return def
}
