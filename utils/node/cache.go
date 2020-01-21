package nodeUtils

import "dijan/core/cache"

func GetCacheNumber() int {
	keys := 0
	scanner := cache.Conn.NewScanner()
	for scanner.Scan() {
		scanner.Key()
		scanner.Value()
		keys += 1
	}
	return keys
}
