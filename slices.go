package go_slices

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func IntInSlice(num int, list []int) bool {
	for _, v := range list {
		if v == num {
			return true
		}
	}
	return false
}