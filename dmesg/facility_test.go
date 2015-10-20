package dmesg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskFacilityExtractsCorrectBits(t *testing.T) {
	v := uint(LOG_USER) << 3
	v |= uint(LOG_ERR)
	assert.EqualValues(t, LOG_USER, MaskFacility(v))
}
