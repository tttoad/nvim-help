package utils_test

import (
	"os"
	"testing"

	"nvim-help/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestConvertEnvPlace(t *testing.T) {
	tests := []struct {
		name   string
		source string
		expect string
	}{
		{
			name:   "convert",
			source: "${GOPATH}/bin",
			expect: os.Getenv("GOPATH") + "/bin",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, utils.ConvertEnvPlace(test.source))
		})
	}
}
