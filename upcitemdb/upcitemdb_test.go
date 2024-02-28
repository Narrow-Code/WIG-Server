package upcitemdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBarcodePossitive(t *testing.T) {
	value := GetBarcode("1234")
	assert.Equal(t, value, 0)
}
