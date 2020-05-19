package atoi

const MAX_INT int32 = int32(^uint32(0) >> 1)
const MIN_INT int32 = ^MAX_INT

type Type string

var TypeEnum = struct {
	Blank, Digital, Operator, Other, Start Type
}{
	"Blank", "Digital", "Operator", "Other", "Start",
}

type Element struct {
	typo  Type
	value int
}

func ParseType(r rune) *Element {
	var typo Type
	var value int
	if r == '+' {
		typo = TypeEnum.Operator
		value = 1
	} else if r == '-' {
		typo = TypeEnum.Operator
		value = -1
	} else if r == ' ' {
		typo = TypeEnum.Blank
	} else if r >= '0' && r <= '9' {
		typo = TypeEnum.Digital
		value = int(r) - int('0')
	} else {
		typo = TypeEnum.Other
	}

	return &Element{typo, value}
}

func fsm(runes []rune) int64 {
	state := TypeEnum.Start
	var product int64 = 0
	var factor int64 = 10
	var sign int64 = 1

loop:
	for _, v := range runes {
		e := ParseType(v)
		switch state {
		case TypeEnum.Start:
			{
				state = e.typo
				if e.typo == TypeEnum.Operator {
					sign = int64(e.value)
				}
				if e.typo == TypeEnum.Digital {
					product = product*factor + int64(e.value)
					state = e.typo
				}
			}
		case TypeEnum.Blank:
			{
				state = e.typo
				if e.typo == TypeEnum.Operator {
					sign = int64(e.value)
				}
				if e.typo == TypeEnum.Digital {
					product = product*factor + int64(e.value)
					state = e.typo
				}
			}
		case TypeEnum.Operator:
			{
				if e.typo == TypeEnum.Operator {
					return 0
				}
				if e.typo == TypeEnum.Blank {
					return 0
				}
				if e.typo == TypeEnum.Digital {
					product = product*factor + int64(e.value)
					state = e.typo
				}
				if e.typo == TypeEnum.Other {
					return 0
				}
			}
		case TypeEnum.Digital:
			{
				if e.typo == TypeEnum.Operator {
					state = TypeEnum.Other
				}
				if e.typo == TypeEnum.Blank {
					state = TypeEnum.Other
				}
				if e.typo == TypeEnum.Digital {
					tmp := product*factor + int64(e.value)
					if tmp < product { // check overflow
						if sign == 1 {
							return int64(MAX_INT)
						} else {
							return int64(MIN_INT)
						}
					}
					product = tmp
				}
				if e.typo == TypeEnum.Other {
					state = TypeEnum.Other
				}
			}
		case TypeEnum.Other:
			{
				break loop
			}
		}
	}
	return product * sign
}

func atoi(str string) int {
	result := fsm([]rune(str))

	if result > int64(MAX_INT) {
		result = int64(MAX_INT)
	}
	if result < int64(MIN_INT) {
		result = int64(MIN_INT)
	}
	return int(result)
}

func trim(str string) string {
	s := []rune(str)
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			return string(s[i:])
		}
	}
	return ""
}
