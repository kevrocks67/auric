package util

import (
	"strings"
)

func ExtractNameFromUri(uri string) string {
	return strings.Split(uri, "/")[3]
}

func ExtractTypeFromUri(uri string) string {
	return strings.Split(uri, "/")[2]
}
