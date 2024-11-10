package kv

import (
	"fmt"
	"time"
)

type cachedValue[T any] struct {
	Value T         `json:"value"`
	Exp   time.Time `json:"exp"`
}

// UseCachedResult caches the result of a function
// result: the function to cache
// storageKey: the key to store the result
// maxDuration: the duration for which the result is valid
// invalidate: if true, the cache is invalidated
func UseCachedResult(
	result func() (string, error),
	storageKey string,
	maxDuration time.Duration,
	invalidate bool,
) (string, error) {
	var empty string

	if !invalidate {
		cachedStr := GetTemporary(storageKey)
		if cachedStr != "" {
			return cachedStr, nil
		}
	}

	resultValue, err := result()
	if err != nil {
		return empty, err
	}

	SetTemporary(storageKey, fmt.Sprintf("%v", resultValue), maxDuration)

	return resultValue, nil
}
