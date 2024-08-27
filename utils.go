package main

func lstContains[T comparable](lst []T, item T) bool {
	for _, element := range lst {
		if element == item {
			return true
		}
	}
	return false
}

func lstIdxOf[T comparable](lst []T, item T) int {
	for i, element := range lst {
		if element == item {
			return i
		}
	}
	return -1
}
