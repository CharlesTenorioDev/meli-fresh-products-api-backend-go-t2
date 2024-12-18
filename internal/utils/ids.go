package utils

func GetBiggestId[T any](db map[int]T) int {
	if len(db) == 0 {
		return 1
	}
	biggest := -1
	for key := range db {
		if key > biggest {
			biggest = key
		}
	}
	return biggest
}
