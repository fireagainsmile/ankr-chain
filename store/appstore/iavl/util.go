package iavl

import (
	"fmt"
	"strings"
)

func containPrefix(key string, prefix string) string {
	return prefix + key
}

func stripKeyPrefix(key string, prefix string) (string, error) {
	strSlice := strings.Split(key, ":")
	if strSlice[0] != prefix {
		return "", fmt.Errorf("invalid key with prefix: key = %s, prefix = %s", key, prefix)
	}

	return strSlice[1], nil
}

func prefixEndBytes(prefix []byte) []byte {
	if prefix == nil {
		return nil
	}

	end := make([]byte, len(prefix))
	copy(end, prefix)

	for {
		if end[len(end)-1] != byte(255) {
			end[len(end)-1]++
			break
		} else {
			end = end[:len(end)-1]
			if len(end) == 0 {
				end = nil
				break
			}
		}
	}

	return end
}
