package utils

func GetBiggestId(db map[int]any) int {
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
