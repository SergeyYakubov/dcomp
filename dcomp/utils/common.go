package utils

func StringInArray(val string, array []string) bool {
	for _, s := range array {
		if s == val {
			return true
		}
	}
	return false
}
