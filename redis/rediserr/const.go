package rediserr

import "errors"

var (
	NOT_INTEGER_ERROR = errors.New("ERR value is not an integer or out of range")

	SYNTAX_ERROR = errors.New("ERR syntax error")
)
