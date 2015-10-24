package main

import (
	"strings"

	"fmt"

	"strconv"

	"time"

	"github.com/jawher/bateau/query"
)

func intCompare(value int, op query.Operator, pattern string) bool {
	ipattern, err := strconv.Atoi(pattern)
	if err != nil {
		panic(fmt.Sprintf("'%s' is not a numeric", op))
	}
	switch op {
	case query.EQ:
		return value == ipattern
	case query.GT:
		return value > ipattern
	default:
		panic(fmt.Sprintf("Unsupported operator %s", op))
	}
}

func strCompare(value string, op query.Operator, pattern string) bool {
	switch op {
	case query.EQ:
		return value == pattern
	case query.LIKE:
		return like(value, pattern)
	default:
		panic(fmt.Sprintf("Unsupported operator %s", op))
	}
}

var durationBaseTime = func() time.Time { return time.Now() }

func durationCompare(value time.Time, op query.Operator, pattern string) bool {
	duration, err := parseDuration(pattern)
	if err != nil {
		panic(err)
	}
	v := durationBaseTime().Sub(value)
	switch op {
	case query.EQ:
		return v == duration
	case query.GT:
		return v.Nanoseconds() > duration.Nanoseconds()
	default:
		panic(fmt.Sprintf("Unsupported operator %s", op))
	}
}

func sizeCompare(value int64, op query.Operator, pattern string) bool {
	against, err := parseSize(pattern)
	if err != nil {
		panic(err)
	}
	switch op {
	case query.EQ:
		return value == against
	case query.GT:
		return value > against
	default:
		panic(fmt.Sprintf("Unsupported operator %s", op))
	}
}

func sliceCompare(values []string, op query.Operator, pattern string) bool {
	for _, value := range values {
		switch op {
		case query.EQ:
			if value == pattern {
				return true
			}
		case query.LIKE:
			if like(value, pattern) {
				return true
			}
		default:
			panic(fmt.Sprintf("Unsupported operator %s", op))
		}
	}
	return false
}

func like(value, pattern string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(pattern))
}

type parser struct {
	input string
	pos   int
}

func (p *parser) eatWs() {
	for ; p.pos < len(p.input); p.pos++ {
		if p.input[p.pos] != ' ' {
			return
		}
	}
}

func (p *parser) eof() bool {
	return p.pos >= len(p.input)
}
