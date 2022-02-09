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
	Inf     int //最大值或者最小值
}

func (border *Border) Greater(value float64) bool {
	if border.Inf == negativeInf {
		return false
	}
	if border.Inf == positiveInf {
		return true
	}
	if border.Include {
		return border.Value >= value
	}

	return border.Value > value
}

func (border *Border) Less(value float64) bool {
	if border.Inf == negativeInf {
		return true
	}
	if border.Inf == positiveInf {
		return false
	}
	if border.Include {
		return border.Value <= value
	}
	return border.Value < value
}

/**
(-inf 和 -inf 语义是一样的
*/
func ParserBorder(s string) (*Border, error) {

	//是否为 inf/+inf 和 -inf

	if strings.ToLower(s) == "-inf" || strings.ToLower(s) == "(-inf" {
		return &Border{
			Value:   0,
			Include: true,
			Inf:     negativeInf,
		}, nil
	}

	if strings.ToLower(s) == "inf" || strings.ToLower(s) == "+inf" || strings.ToLower(s) == "(inf" || strings.ToLower(s) == "(+inf" {
		return &Border{
			Value:   0,
			Include: true,
			Inf:     positiveInf,
		}, nil
	}

	if s[0] == '(' {
		value, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return nil, errors.New("ERR min or max is not a float")
		}
		return &Border{
			Value:   value,
			Include: false,
			Inf:     0,
		}, nil
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.New("ERR min or max is not a float")
	}
	return &Border{
		Inf:     0,
		Value:   value,
		Include: true,
	}, nil
}
