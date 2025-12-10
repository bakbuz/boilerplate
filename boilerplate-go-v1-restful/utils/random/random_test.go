package random_test

import (
	"testing"

	"codegen/utils/random"

	"github.com/stretchr/testify/assert"
)

func Test_RandomStr_6(t *testing.T) {
	actual := random.Str(6)
	assert.Len(t, actual, 6)
}

func Test_RandomStr_16(t *testing.T) {
	actual := random.Str(16)
	assert.Len(t, actual, 16)
}

func Test_RandomStr_64(t *testing.T) {
	actual := random.Str(64)
	assert.Len(t, actual, 64)
}
