package border

import (
	"errors"
	"strconv"
	"strings"
)

const (
	positiveInf = 1
	negativeInf = -1
)

/*
	如果包含 ( 那么 include 为 false
*/
type Border struct {
	Value   float64
	Include bool
	Inf     int
}

func ParserBorder(s string) (*Border, error) {

	//是否为 inf/+inf 和 -inf
	if strings.ToLower(s) == "(-inf" {
		return &Border{
			Value:   0,
			Include: false,
			Inf:     -1,
		}, nil
	}
	if strings.ToLower(s) == "-inf" {
		return &Border{
			Value:   0,
			Include: true,
			Inf:     -1,
		}, nil
	}

	if strings.ToLower(s) == "inf" || strings.ToLower(s) == "+inf" {
		return &Border{
			Value:   0,
			Include: true,
			Inf:     1,
		}, nil
	}

	if strings.ToLower(s) == "(inf" || strings.ToLower(s) == "(+inf" {
		return &Border{
			Value:   0,
			Include: false,
			Inf:     1,
		}, nil
	}

	value, err := strconv.ParseFloat(s, 64)
	if err == nil {
		if value >= 0 {
			return &Border{
				Value:   value,
				Include: true,
				Inf:     1,
			}, nil
		} else {
			return &Border{
				Value:   value,
				Include: true,
				Inf:     -1,
			}, nil
		}
	}

	//是否包含 (
	hasParentheses := strings.Contains(s, "(")
	if hasParentheses {
		value, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return nil, err
		}

		if value >= 0 {
			return &Border{
				Value:   value,
				Include: false,
				Inf:     1,
			}, nil
		} else {
			return &Border{
				Value:   value,
				Include: false,
				Inf:     -1,
			}, nil
		}
	}

	return nil, errors.New("(error) ERR min or max is not a float")
}
