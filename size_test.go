package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSize(t *testing.T) {
	cases := []struct {
		input string
		ok    int64
		ko    bool
	}{
		{input: "", ko: true},
		{input: "   ", ko: true},
		{input: "KB", ko: true},
		{input: "1xz", ko: true},

		{input: "42", ok: 42},

		{input: "42KB", ok: 42 * 1024},
		{input: "42Kb", ok: 42 * 1000},
		{input: "42kb", ok: 42 * 1000},

		{input: "42MB", ok: 42 * 1024 * 1024},
		{input: "42Mb", ok: 42 * 1000 * 1000},

		{input: "42GB", ok: 42 * 1024 * 1024 * 1024},
		{input: "42Gb", ok: 42 * 1000 * 1000 * 1000},

		{
			input: "1Gb 2MB 3kb 4",
			ok: 1*1000*1000*1000 +
				2*1024*1024 +
				3*1000 +
				4,
		},
	}

	for _, cas := range cases {
		size, err := parseSize(cas.input)
		if cas.ko {
			t.Log(err)
			require.Error(t, err, "parsing '%s' should have failed", cas.input)
			continue
		}
		require.NoError(t, err, "'%s' should not have generated parse error %v", cas.input, err)
		require.Equal(t, cas.ok, size, "input '%s' should parse to %v, instead got %v", cas.input, cas.ok, size)
	}
}
