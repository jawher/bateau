package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseSize(input string) (int64, error) {
	return (&sizeParser{&parser{input: input}}).parse()
}

type sizeParser struct {
	*parser
}

func (p *sizeParser) parse() (res int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v\n%s\n%s^", r, p.input, strings.Repeat(" ", p.pos))
		}
	}()

	err = nil
	empty := true

	for !p.eof() {
		p.eatWs()
		res += p.parseNum() * p.parseUnit()
		empty = false
	}
	if empty {
		panic("no value")
	}

	return
}

func (p *sizeParser) parseNum() int64 {
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

var sizeUnitMultipliers = map[string]int64{
	"KB": 1024,
	"kb": 1000,
	"Kb": 1000,
	"MB": 1024 * 1024,
	"Mb": 1000 * 1000,
	"GB": 1024 * 1024 * 1024,
	"Gb": 1000 * 1000 * 1000,
}

var sizeUnits = []string{"KB", "Kb", "kb", "MB", "Mb", "GB", "Gb"}

func (p *sizeParser) parseUnit() int64 {
	if p.pos >= len(p.input) {
		return 1
	}

	for _, u := range sizeUnits {
		if strings.HasPrefix(p.input[p.pos:], u) {
			p.pos += len(u)
			return sizeUnitMultipliers[u]
		}
	}
	panic("was expecting a size unit")
}
