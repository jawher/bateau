package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		input string
		ok    time.Duration
		ko    bool
	}{
		{input: "", ko: true},
		{input: "   ", ko: true},
		{input: "d", ko: true},
		{input: "1", ko: true},
		{input: "42ms", ok: 42 * time.Millisecond},
		{input: "42s", ok: 42 * time.Second},
		{input: "42m", ok: 42 * time.Minute},
		{input: "42h", ok: 42 * time.Hour},
		{input: "42d", ok: 42 * 24 * time.Hour},
		{input: "42M", ok: 42 * 30 * 24 * time.Hour},
		{input: "42months", ok: 42 * 30 * 24 * time.Hour},
		{input: "42y", ok: 42 * 365 * 24 * time.Hour},
		{input: "1y 2M 3w 4d 5h 6m 7s 8ms ",
			ok: 365*24*time.Hour +
				2*30*24*time.Hour +
				3*7*24*time.Hour +
				4*24*time.Hour +
				5*time.Hour +
				6*time.Minute +
				7*time.Second +
				8*time.Millisecond,
		},
	}

	for _, cas := range cases {
		dur, err := parseDuration(cas.input)
		if cas.ko {
			t.Log(err)
			require.Error(t, err, "parsing '%s' should have failed", cas.input)
			continue
		}
		require.NoError(t, err, "'%s' should not have generated parse error %v", cas.input, err)
		require.Equal(t, cas.ok, dur, "input '%s' should parse to %v, instead got %v", cas.input, cas.ok, dur)
	}
}
