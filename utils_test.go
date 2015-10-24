package main

import (
	"testing"

	"time"

	"github.com/jawher/bateau/query"
	"github.com/stretchr/testify/require"
)

func TestIntCompare(t *testing.T) {
	require.True(t, intCompare(1, query.EQ, "1"))
	require.False(t, intCompare(1, query.EQ, "2"))

	require.True(t, intCompare(2, query.GT, "1"))
	require.False(t, intCompare(1, query.GT, "1"))
}

func TestStrCompare(t *testing.T) {
	require.True(t, strCompare("test", query.EQ, "test"))
	require.False(t, strCompare("test", query.EQ, "niet"))

	require.True(t, strCompare("test", query.LIKE, "test"))
	require.True(t, strCompare("test", query.LIKE, "tes"))
	require.True(t, strCompare("test", query.LIKE, "est"))
	require.True(t, strCompare("tEsT", query.LIKE, "eSt"))
	require.False(t, strCompare("test", query.LIKE, "niet"))
}

func TestDurationCompare(t *testing.T) {
	base := time.Now()
	originalDurationBaseTime := durationBaseTime
	durationBaseTime = func() time.Time {
		return base
	}
	defer func() {
		durationBaseTime = originalDurationBaseTime
	}()

	require.True(t, durationCompare(base.Add(-1*time.Hour), query.EQ, "1h"))
	require.True(t, durationCompare(base.Add(-1*time.Hour), query.GT, "1m"))

	require.False(t, durationCompare(base.Add(-1*time.Hour), query.EQ, "1h 1m"))
	require.False(t, durationCompare(base.Add(-1*time.Hour), query.GT, "1h 1s"))
}

func TestSizeCompare(t *testing.T) {
	require.True(t, sizeCompare(42, query.EQ, "42"))
	require.True(t, sizeCompare(42*1000, query.EQ, "42kb"))
	require.True(t, sizeCompare(42*1024, query.EQ, "42KB"))
	require.False(t, sizeCompare(42*1024, query.EQ, "42Mb"))

	require.True(t, sizeCompare(42*1024, query.GT, "42"))
	require.True(t, sizeCompare(42*1024, query.GT, "42kb"))
	require.False(t, sizeCompare(42*1024, query.GT, "42MB"))
}

func TestSliceCompare(t *testing.T) {
	require.True(t, sliceCompare([]string{"test"}, query.EQ, "test"))
	require.True(t, sliceCompare([]string{"niet", "test"}, query.EQ, "test"))
	require.False(t, sliceCompare([]string{"niet", "test"}, query.EQ, "tEst"))

	require.True(t, sliceCompare([]string{"niet", "test"}, query.LIKE, "Est"))
	require.False(t, sliceCompare([]string{"niet", "test"}, query.LIKE, "42"))

}
