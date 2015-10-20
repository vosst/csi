package dmesg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaskLoglevelExtractsCorrectBits(t *testing.T) {
	v := uint(LOG_USER) << 3
	v |= uint(LOG_ERR)

	assert.EqualValues(t, LOG_ERR, MaskLoglevel(v))
}
