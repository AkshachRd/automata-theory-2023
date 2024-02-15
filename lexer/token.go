package main

const DIGITS = "0123456789"

const TT_INT = "INT"
const TT_FLOAT = "FLOAT"
const TT_PLUS = "PLUS"
const TT_MINUS = "MINUS"
const TT_MUL = "MUL"
const TT_DIV = "DIV"
const TT_LPAREN = "LPAREN"
const TT_RPAREN = "RPAREN"

type Token struct {
	Type  string
	Value interface{}
}

func NewToken(tokenType string, value ...interface{}) *Token {
	var val interface{}
	if len(value) > 0 {
		val = value[0]
	}
	return &Token{tokenType, val}
}
