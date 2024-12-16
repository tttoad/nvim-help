package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMM(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		tag    string
		expect string
	}{
		{
			name:   "test name",
			str:    "AA string",
			tag:    "json",
			expect: "AA string `json:\"a_a\"`",
		},
		{
			name:   "",
			str:    "ServiceName int `json:\"service_name\"`",
			tag:    "yaml",
			expect: "ServiceName int `json:\"service_name\" yaml:\"serviceName\"`",
		},
		{
			name:   "",
			str:    "AgentManager int `json:\"agent_manager\"`",
			tag:    "json",
			expect: "AgentManager int `json:\"agent_manager\"`",
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(test.expect, addTags(test.str, test.tag))
		})
	}
}
