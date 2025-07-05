package main

import (
	"fmt"
	"unicode"
)

const (
	Id uint8 = iota // any sequence of characters not containing whitespace, :, [, ] or '
	String
	Special
)

type token struct {
	Type uint8
	text string
}

func scan(sink chan token, runes []rune) {
	var pos uint = 0

	var next = func() {
		pos++
	}

	var current = func() rune {
		if pos < uint(len(runes)) {
			return runes[pos]
		} else {
			return 0
		}
	}

	var emit = func(tt uint8, from, to uint) {
		sink <- token{Type: tt, text: string(runes[from:to])}
	}

	var currentIsAlphaOrDash = func() bool {
		return unicode.IsLetter(current()) || unicode.IsDigit(current()) || current() == '-'
	}

	var readString = func() {
		var start = pos
		for {
			next()
			if current() == '\'' {
				emit(String, start, pos)
				next()
				break
			} else if current() == 0 {
				panic(fmt.Sprintf("Runaway string, starting at %d", start))
			}
		}
	}

	var readComment = func() {
		for {
			next()
			if current() == '\n' || current() == 0 {
				break
			}
		}
	}

	var readId = func() {
		var start = pos
		for {
			next()
			if !currentIsAlphaOrDash() {
				emit(Id, start, pos)
				break
			}
		}
	}

	for current() != 0 {
		if currentIsAlphaOrDash() {
			readId()
		} else if current() == '\'' {
			readString()
		} else if current() == '#' {
			readComment()
		} else if unicode.IsSpace(current()) {
			next()
		} else {
			emit(Special, pos, pos+1)
			next()
		}
	}
	close(sink)
}
