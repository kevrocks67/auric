package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractNameFromGoldenUri(t *testing.T) {
	name := ExtractNameFromGoldenUri("/golden/qcow2/rocky9-base/prod")
	assert.Equal(t, name, "rocky9-base")
}

func TestExtractTypeFromGoldenUri(t *testing.T) {
	name := ExtractTypeFromGoldenUri("/golden/qcow2/rocky9-base/prod")
	assert.Equal(t, name, "qcow2")
}
