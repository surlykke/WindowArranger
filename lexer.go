package main

import (
	"regexp"
)

type TokenType string

const (
	String     TokenType = "string"
	Identifier           = "identifier"
	Other                = "other"
)

type Token struct {
	Text string
	Type TokenType
	Line int
	Col  int
}

func Tokenize(bytes []byte) chan Token {
	var tokens = make(chan Token)
	go scan(bytes, tokens)
	return tokens	
}

func scan(bytes []byte, tokenSink chan Token) {
	defer close(tokenSink)

	var identifier = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\-]*`) // Assume monitornames never start with a number
	var singleQuoteString = regexp.MustCompile(`^'.*?'`)
	var comment = regexp.MustCompile(`^\#.*`)


	var pos int = 0
	var linenumber int = 1
	var linestart int = 0

	var emit = func(ttype TokenType, length int) {
		if ttype == String {
			tokenSink <- Token{Text: string(bytes[pos+1 : pos+length-1]), Type: ttype, Line: linenumber, Col: pos - linestart}
		} else {
			tokenSink <- Token{Text: string(bytes[pos : pos+length]), Type: ttype, Line: linenumber, Col: pos - linestart}
		}
		pos = pos + length
	}

	for {
		for pos < len(bytes) && (' ' == bytes[pos] || '\t' == bytes[pos] || '\n' == bytes[pos]) {
			if '\n' == bytes[pos] {
				linenumber = linenumber + 1
				linestart = pos + 1
			}
			pos = pos + 1
		}

		if pos >= len(bytes) {
			return
		}
		var match []int
		if match = identifier.FindIndex(bytes[pos:]); match != nil {
			emit(Identifier, match[1])
		} else if match = singleQuoteString.FindIndex(bytes[pos:]); match != nil {
			emit(String, match[1])
		} else if match = comment.FindIndex(bytes[pos:]); match != nil {
			pos = pos + match[1]
		} else {
			emit(Other, 1)
		}
	}
}
