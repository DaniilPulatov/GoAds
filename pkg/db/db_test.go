package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	t.Run("check invalid dsn", func(t *testing.T) {
		_, err := NewDB("invalid-dsn")
		assert.Error(t, err)
	})
	t.Run("check valid dsn", func(t *testing.T) {
		pool, err := NewDB("postgres://user:1234@localhost:5432/ads_db")
		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}
