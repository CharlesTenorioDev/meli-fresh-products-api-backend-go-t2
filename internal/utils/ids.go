package utils

// Returns the biggest id of a int map
// if len(map) == 0, returns 1
func GetBiggestId[K int, V any](db map[K]V) int {
	if len(db) == 0 {
		return 1
	}

	biggest := -1

	for key := range db {
		k2 := int(key)
		if k2 > biggest {
			biggest = k2
		}
	}

	return biggest
}
