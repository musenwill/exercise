package calculator

import (
	"errors"
	"fmt"
	"github.com/musenwill/exercise/stack"
	"strconv"
)

type runeType string

var runeTypeEnum = struct {
	digital, operator, blank, unknown runeType
}{
	"digital", "operator", "blank", "unknown",
}

func checkType(c string) runeType {
	if c >= "0" && c <= "9" {
		return runeTypeEnum.digital
	}
	if c == "+" || c == "-" || c == "*" || c == "/" || c == "(" || c == ")" {
		return runeTypeEnum.operator
	}
	if c == " " || c == "\n" || c == "\t" || c == "\r" {
		return runeTypeEnum.blank
	}
	return runeTypeEnum.unknown
}

type item struct {
	t        runeType
	value    int64
	name     string
	priority int
}

func createOperator(c string) (*item, error) {
	if checkType(c) != runeTypeEnum.operator {
		return nil, errors.New("invalid operator type")
	}

	priority := 0
	switch c {
	case "+":
		priority = 1
	case "-":
		priority = 1
	case "*":
		priority = 2
	case "/":
		priority = 2
	case "(":
		priority = 3
	case ")":
		priority = 3
	}

	return &item{
		t:        runeTypeEnum.operator,
		value:    0,
		name:     c,
		priority: priority,
	}, nil
}

func createNumber(c string) (*item, error) {
	value, err := strconv.Atoi(c)
	if err != nil {
		return nil, err
	}

	return &item{
		t:        runeTypeEnum.digital,
		value:    int64(value),
		name:     c,
		priority: 0,
	}, nil
}

func postOrderExpression(expression string) (*stack.Stack, error) {
	sk := stack.CreateStack(make([]interface{}, 0))
	operatorStack := stack.CreateStack(make([]interface{}, 0))

	for i := 0; i < len(expression); i++ {
		c := string(expression[i])
		t := checkType(c)
		switch t {
		case runeTypeEnum.unknown:
			return nil, fmt.Errorf("invalid character in expression at %d", i)
		case runeTypeEnum.blank:
			continue
		case runeTypeEnum.operator:
			{
				o, _ := createOperator(c)
				if operatorStack.Empty() {
					operatorStack.Push(o)
				} else {
					if operatorStack.Peek().(*item).priority < o.priority {
						operatorStack.Push(o)
					} else {
						for !operatorStack.Empty() && operatorStack.Peek().(*item).priority >= o.priority {
						}
					}
				}
			}
		}
	}

	return sk, nil
}
