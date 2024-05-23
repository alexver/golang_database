package database

import (
	"testing"

	"github.com/alexver/golang_database/internal/query"
	"github.com/stretchr/testify/require"
)

func TestPingProcessor_Name(t *testing.T) {
	t.Run("Test Name", func(t *testing.T) {
		p := NewPingProcessor()
		require.Equal(t, "PING", p.Name())
	})
}

func TestPingProcessor_Suports(t *testing.T) {
	t.Run("Test Supports true", func(t *testing.T) {
		p := NewPingProcessor()
		q := query.CreateQuery("PING", []string{})
		require.True(t, p.Suports(q))
	})

	t.Run("Test Supports false", func(t *testing.T) {
		p := NewPingProcessor()
		q := query.CreateQuery("PONG", []string{})
		require.False(t, p.Suports(q))
	})
}

func TestPingProcessor_Process(t *testing.T) {
	t.Run("Test Process", func(t *testing.T) {
		p := NewPingProcessor()
		q := query.CreateQuery("PING", []string{})
		result, err := p.Process(q)
		require.Equal(t, "[ok] PONG", result)
		require.Nil(t, err)
	})
}
