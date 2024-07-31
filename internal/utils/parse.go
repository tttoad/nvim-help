package utils

import (
	"os"
	"regexp"
	"strings"
)

// ConvertEnvPlace ${GOPATH}/bin -> /Users/***/go/bin
func ConvertEnvPlace(str string) string {
	r, _ := regexp.Compile(`\${(.+)}`)
	for _, s := range r.FindStringSubmatch(str) {
		es := os.Getenv(s)
		if es != "" {
			str = strings.ReplaceAll(str, "${"+s+"}", es)
		}
	}
	return str
}
