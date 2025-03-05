package sl

import (
	"errors"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestErr(t *testing.T) {
	assert := assert2.New(t)

	customErr := errors.New("custom error")
	attr := Err(customErr)

	t.Run("Not nil", func(t *testing.T) {
		assert.NotNil(attr)
	})

	t.Run("Correct key", func(t *testing.T) {
		assert.Equal("error", attr.Key)
	})

	t.Run("Correct value", func(t *testing.T) {
		assert.Equal(customErr.Error(), attr.Value.String())
	})
}
