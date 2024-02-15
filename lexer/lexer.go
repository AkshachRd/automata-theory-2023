package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
	text        []rune
	pos         uint64
	currentChar *rune
}

func NewLexer(text string) *Lexer {
	lexer := &Lexer{text: []rune(text)}
	lexer.currentChar = &lexer.text[lexer.pos]
	return lexer
}

func (l *Lexer) Advance() {
	l.pos++
	if l.pos < uint64(len(l.text)) {
		l.currentChar = &l.text[l.pos]
	} else {
		l.currentChar = nil
	}
}

func (l *Lexer) MakeTokens() ([]*Token, error) {
	var tokens []*Token

	for l.currentChar != nil {
		switch {
		case strings.ContainsRune(" \t", *l.currentChar):
			l.Advance()
		case strings.ContainsRune(DIGITS, *l.currentChar):
			token, err := l.MakeNumber()
			if err != nil {
				return make([]*Token, 0), err
			}

			tokens = append(tokens, token)
		case *l.currentChar == '+':
			tokens = append(tokens, NewToken(TT_PLUS))
			l.Advance()
		case *l.currentChar == '-':
			tokens = append(tokens, NewToken(TT_MINUS))
			l.Advance()
		case *l.currentChar == '*':
			tokens = append(tokens, NewToken(TT_MUL))
			l.Advance()
		case *l.currentChar == '/':
			tokens = append(tokens, NewToken(TT_DIV))
			l.Advance()
		case *l.currentChar == '(':
			tokens = append(tokens, NewToken(TT_LPAREN))
			l.Advance()
		case *l.currentChar == ')':
			tokens = append(tokens, NewToken(TT_RPAREN))
			l.Advance()
		default:
			char := l.currentChar
			l.Advance()
			return make([]*Token, 0), fmt.Errorf("illigal char: %q", *char)
		}
	}

	return tokens, nil
}

func (l *Lexer) MakeNumber() (*Token, error) {
	numStr := ""
	dotCount := 0

	for l.currentChar != nil && strings.ContainsRune(DIGITS+".", *l.currentChar) {
		if *l.currentChar == '.' {
			if dotCount == 1 {
				break
			}
			dotCount++
			numStr += "."
		} else {
			numStr += string(*l.currentChar)
		}
		l.Advance()
	}

	if dotCount == 0 {
		num, err := strconv.Atoi(numStr)
		return NewToken(TT_INT, num), err
	}

	num, err := strconv.ParseFloat(numStr, 32)
	return NewToken(TT_FLOAT, float32(num)), err
}
