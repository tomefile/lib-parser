package libparser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

type StringModifier func(string) string

func (modifier StringModifier) MarshalJSON() ([]byte, error) {
	return json.Marshal("modifier()")
}

func GetModifier(name string, args []string) StringModifier {
	switch name {

	case "to_lower":
		return strings.ToLower

	case "to_upper":
		return strings.ToUpper

	case "to_snake":
		return strcase.ToSnake

	case "to_screaming_snake":
		return strcase.ToScreamingSnake

	case "to_kebab":
		return strcase.ToKebab

	case "to_screaming_kebab":
		return strcase.ToScreamingKebab

	case "to_camel":
		return strcase.ToLowerCamel

	case "to_pascal":
		return strcase.ToCamel

	case "to_delimited":
		if len(args) == 0 {
			return nil
		}
		return func(in string) string {
			return strcase.ToDelimited(in, uint8(args[0][0]))
		}

	case "length":
		return func(in string) string {
			return fmt.Sprint(len(in))
		}

	case "quoted":
		return func(in string) string {
			return fmt.Sprintf("%q", in)
		}

	case "trim":
		return func(in string) string {
			for _, arg := range args {
				in = strings.Trim(in, arg)
			}
			return in
		}

	case "trim_prefix":
		return func(in string) string {
			for _, arg := range args {
				in = strings.TrimPrefix(in, arg)
			}
			return in
		}

	case "trim_suffix":
		return func(in string) string {
			for _, arg := range args {
				in = strings.TrimSuffix(in, arg)
			}
			return in
		}

	case "pad":
		switch len(args) {
		case 1:
			number, err := strconv.Atoi(args[0])
			if err == nil {
				padding := strings.Repeat(" ", max(0, number))
				return func(in string) string {
					return padding + in + padding
				}
			}
		case 2:
			number_1, err_1 := strconv.Atoi(args[0])
			number_2, err_2 := strconv.Atoi(args[1])
			if err_1 == nil && err_2 == nil {
				return func(in string) string {
					return strings.Repeat(" ", max(0, number_1)) +
						in +
						strings.Repeat(" ", max(0, number_2))
				}
			}
		default:
			return nil
		}

	case "pad_left":
		if len(args) == 0 {
			return nil
		}
		number, err := strconv.Atoi(args[0])
		if err == nil {
			padding := strings.Repeat(" ", max(0, number))
			return func(in string) string {
				return padding + in
			}
		}

	case "pad_right":
		if len(args) == 0 {
			return nil
		}
		number, err := strconv.Atoi(args[0])
		if err == nil {
			padding := strings.Repeat(" ", max(0, number))
			return func(in string) string {
				return in + padding
			}
		}

	case "reverse":
		// https://stackoverflow.com/a/10030772
		return func(in string) string {
			runes := []rune(in)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}

	case "invert":
		// https://stackoverflow.com/a/38234154
		return func(in string) string {
			return strings.Map(func(char rune) rune {
				switch {
				case unicode.IsLower(char):
					return unicode.ToUpper(char)
				case unicode.IsUpper(char):
					return unicode.ToLower(char)
				}
				return char
			}, in)
		}

	}

	return nil
}
