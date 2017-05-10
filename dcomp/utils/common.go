package utils

import "time"

func StringInArray(val string, array []string) bool {
	for _, s := range array {
		if s == val {
			return true
		}
	}
	return false
}

func TimeToString(val time.Time) string {
	return val.Format(time.RFC3339)
}

func StringToTime(val string) (time.Time, error) {
	return time.Parse(time.RFC3339, val)
}
