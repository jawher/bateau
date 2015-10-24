package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseDuration(input string) (time.Duration, error) {
	return (&durationParser{&parser{input: input}}).parse()
}

type durationParser struct {
	*parser
}

func (p *durationParser) parse() (res time.Duration, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v\n%s\n%s^", r, p.input, strings.Repeat(" ", p.pos))
		}
	}()

	err = nil
	var d int64 = 0
	empty := true

	p.eatWs()

	for !p.eof() {
		d += p.parseNum() * p.parseUnit()
		empty = false
		p.eatWs()
	}
	if empty {
		panic("no value")
	}
	res = time.Duration(1000 * 1000 * d)

	return
}

func (p *durationParser) parseNum() int64 {
	start := p.pos
	for ; p.pos < len(p.input); p.pos++ {
		c := p.input[p.pos]
		if c < '0' || c > '9' {
			break
		}
	}
	if p.pos == start {
		panic("was expecting a numeric value")
	}
	res, err := strconv.ParseInt(p.input[start:p.pos], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("value %s is too large", p.input[start:p.pos]))
	}
	return res
}

var durationUnitMultipliers = map[string]int64{
	"ms":     1,
	"s":      1000,
	"m":      60 * 1000,
	"h":      60 * 60 * 1000,
	"d":      24 * 60 * 60 * 1000,
	"w":      7 * 24 * 60 * 60 * 1000,
	"M":      30 * 24 * 60 * 60 * 1000,
	"months": 30 * 24 * 60 * 60 * 1000,
	"y":      365 * 24 * 60 * 60 * 1000,
}

var durationUnits = []string{"months", "ms", "s", "m", "h", "d", "w", "M", "y"}

func (p *durationParser) parseUnit() int64 {
	if p.pos >= len(p.input) {
		panic("was expecting a duration unit")
	}

	for _, u := range durationUnits {
		if strings.HasPrefix(p.input[p.pos:], u) {
			p.pos += len(u)
			return durationUnitMultipliers[u]
		}
	}
	panic("was expecting a duration unit")
}
