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

	t.Run("check empty dsn", func(t *testing.T) {
		_, err := NewDB("")
		assert.Error(t, err)
	})

	t.Run("check valid dsn", func(t *testing.T) {
		t.Skip("Требуется запущенная локальная база данных PostgresSQL")
		// dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
		// db, err := NewDB(dsn)
		// require.NoError(t, err)
		// assert.NotNil(t, db)
		// db.Close()
	})
}

func TestConn_Ping(t *testing.T) {
	t.Skip("Требуется мок или реальная база данных")
}

func TestConn_QueryRow(t *testing.T) {
	t.Skip("Требуется мок или реальная база данных")
}

func TestConn_Query(t *testing.T) {
	t.Skip("Требуется мок или реальная база данных")
}

func TestConn_Begin(t *testing.T) {
	t.Skip("Требуется мок или реальная база данных")
}

func TestConn_Exec(t *testing.T) {
	t.Skip("Требуется мок или реальная база данных")
}
