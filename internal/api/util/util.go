package util

import (
	"fmt"
	"strings"
)

func ExtractNameFromUri(uri string) string {
	return strings.Split(uri, "/")[3]
}

func ExtractTypeFromUri(uri string) string {
	fmt.Println(uri)
	return strings.Split(uri, "/")[2]
}
