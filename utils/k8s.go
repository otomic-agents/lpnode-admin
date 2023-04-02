package utils

import (
	"strings"
)

func ParseServiceNameFromDeplayMessage(message string) string {
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		lineField := strings.Split(line, " ")
		if len(lineField) <= 1 {
			continue
		}
		nameField := lineField[0]
		if strings.HasPrefix(nameField, "service/") {
			return strings.Replace(nameField, "service/", "", 1)
		}
	}
	return ""
}
