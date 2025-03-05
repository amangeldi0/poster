package config

import (
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)

	t.Run("file not exists", func(t *testing.T) {
		filePath := "./not-found-file"
		os.Clearenv()
		err := os.Setenv(PathKey, filePath)
		assert.Nil(err, "setting environment variable should succeed")

		_, err = os.Stat(filePath)
		assert.NotNil(err, "test config file not must exist")

	})

	t.Run("success load", func(t *testing.T) {
		filePath := "./testconfigs/config_test.yaml"
		os.Clearenv()
		err := os.Setenv(PathKey, filePath)

		assert.Nil(err, "setting environment variable should succeed")

		_, err = os.Stat(filePath)
		assert.Nil(err, "test config file must exist")

		_, err = New()

		assert.Nil(err, "config must load")
	})

	t.Run("failure load", func(t *testing.T) {
		filePath := "incorrect-path"
		os.Clearenv()
		err := os.Setenv(PathKey, filePath)

		assert.Nil(err, "setting environment variable should succeed")

		_, err = os.Stat(filePath)
		assert.Nil(err, "test config file not must exist")

		_, err = New()

		assert.NotNil(err, "config not must load")
	})

	t.Run("failure config file", func(t *testing.T) {
		filePath := "./testconfigs/bad_config_test.yaml"
		os.Clearenv()
		err := os.Setenv(PathKey, filePath)

		assert.Nil(err, "setting environment variable should succeed")

		_, err = os.Stat(filePath)

		cfg, err := New()

		assert.NotNil(err, "config must fail to load")

		if err != nil {
			assert.Contains(err.Error(), "cannot read config:", "Error message should indicate config reading failure")
		}

		assert.Nil(cfg, "Config object must be nil on failure")
	})

	t.Run("invalid mailer port", func(t *testing.T) {
		filePath := "./testconfigs/bad_config_test2.yaml"

		os.Clearenv()
		err := os.Setenv(PathKey, filePath)

		assert.Nil(err, "setting environment variable should succeed")

		_, err = os.Stat(filePath)

		cfg, err := New()

		assert.NotNil(err, "config must fail to load")

		if err != nil {
			assert.Contains(err.Error(), "invalid mailer port", "Error message should indicate converting mailer port failure")
		}

		assert.Nil(cfg, "Config object must be nil on failure")
	})

}
