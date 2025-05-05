package utils

import "strings"

func Includes(slice []string, value string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}
