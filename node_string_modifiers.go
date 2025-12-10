package libparser

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

type ModifierName string

const (
	MOD_NOT          ModifierName = "not"
	MOD_TO_LOWER     ModifierName = "to_lower"
	MOD_TO_UPPER     ModifierName = "to_upper"
	MOD_TO_SNAKE     ModifierName = "to_snake"
	MOD_TO_KEBAB     ModifierName = "to_kebab"
	MOD_TO_CAMEL     ModifierName = "to_camel"
	MOD_TO_PASCAL    ModifierName = "to_pascal"
	MOD_TO_DELIMITED ModifierName = "to_delimited"
	MOD_DOES_EXIST   ModifierName = "does_exist"
	MOD_IS_EMPTY     ModifierName = "is_empty"
	MOD_IS_FILE      ModifierName = "is_file"
	MOD_IS_DIR       ModifierName = "is_dir"
	MOD_IS_SYMLINK   ModifierName = "is_symlink"
	MOD_LENGTH       ModifierName = "length"
	MOD_QUOTED       ModifierName = "quoted"
	MOD_TRIM         ModifierName = "trim"
	MOD_TRIM_PREFIX  ModifierName = "trim_prefix"
	MOD_TRIM_SUFFIX  ModifierName = "trim_suffix"
	MOD_PAD          ModifierName = "pad"
	MOD_PAD_LEFT     ModifierName = "pad_left"
	MOD_PAD_RIGHT    ModifierName = "pad_right"
	MOD_HAS_PREFIX   ModifierName = "has_prefix"
	MOD_HAS_SUFFIX   ModifierName = "has_suffix"
	MOD_SLICE        ModifierName = "slice"
	MOD_REVERSE      ModifierName = "reverse"
	MOD_INVERT       ModifierName = "invert"
)

var AllModifiers = []ModifierName{
	MOD_NOT,
	MOD_TO_LOWER,
	MOD_TO_UPPER,
	MOD_TO_SNAKE,
	MOD_TO_KEBAB,
	MOD_TO_CAMEL,
	MOD_TO_PASCAL,
	MOD_TO_DELIMITED,
	MOD_DOES_EXIST,
	MOD_IS_EMPTY,
	MOD_IS_FILE,
	MOD_IS_DIR,
	MOD_IS_SYMLINK,
	MOD_LENGTH,
	MOD_QUOTED,
	MOD_TRIM,
	MOD_TRIM_PREFIX,
	MOD_TRIM_SUFFIX,
	MOD_PAD,
	MOD_PAD_LEFT,
	MOD_PAD_RIGHT,
	MOD_HAS_PREFIX,
	MOD_HAS_SUFFIX,
	MOD_SLICE,
	MOD_REVERSE,
	MOD_INVERT,
}

type StringModifier struct {
	Name ModifierName
	Args []*NodeString
	Call func(Locals, string) string
}

func (modifier StringModifier) String() string {
	if len(modifier.Args) == 0 {
		return string(modifier.Name)
	}

	var builder strings.Builder
	for _, arg := range modifier.Args {
		builder.WriteString(" " + arg.String())
	}

	return fmt.Sprintf("%s %s", modifier.Name, builder.String())
}

func (modifier StringModifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("(%s)", modifier.String()))
}

// ————————————————————————————————

func SortNotModifierToEnd(modifiers []StringModifier) []StringModifier {
	var output = make([]StringModifier, len(modifiers))
	i := 0
	j := len(output) - 1
	for _, modifier := range modifiers {
		if modifier.Name == MOD_NOT {
			output[j] = modifier
			j--
		} else {
			output[i] = modifier
			i++
		}
	}
	return output
}

func GetModifier(name ModifierName, args []*NodeString) (StringModifier, error) {
	mod := StringModifier{
		Name: name,
		Args: args,
		Call: nil,
	}

	switch mod.Name {

	case MOD_NOT:
		mod.Call = func(_ Locals, in string) string {
			switch in {
			case boolToString(true), "true", "TRUE":
				return boolToString(false)
			default:
				return boolToString(true)
			}
		}

	case MOD_TO_LOWER:
		mod.Call = func(_ Locals, in string) string {
			return strings.ToLower(in)
		}

	case MOD_TO_UPPER:
		mod.Call = func(_ Locals, in string) string {
			return strings.ToUpper(in)
		}

	case MOD_TO_SNAKE:
		mod.Call = func(_ Locals, in string) string {
			return strcase.ToSnake(in)
		}

	case MOD_TO_KEBAB:
		mod.Call = func(_ Locals, in string) string {
			return strcase.ToKebab(in)
		}

	case MOD_TO_CAMEL:
		mod.Call = func(_ Locals, in string) string {
			return strcase.ToLowerCamel(in)
		}

	case MOD_TO_PASCAL:
		mod.Call = func(_ Locals, in string) string {
			return strcase.ToCamel(in)
		}

	case MOD_TO_DELIMITED:
		if len(mod.Args) > 0 && len(mod.Args[0].Segments) > 0 {
			mod.Call = func(locals Locals, in string) string {
				str, _ := mod.Args[0].Eval(locals)
				if len(str) == 0 {
					return in
				}
				return strcase.ToDelimited(in, uint8(str[0]))
			}
		}

	case MOD_DOES_EXIST:
		mod.Call = func(_ Locals, in string) string {
			_, err := os.Stat(in)
			return boolToString(!os.IsNotExist(err))
		}

	case MOD_IS_EMPTY:
		mod.Call = func(_ Locals, in string) string {
			return boolToString(len(in) == 0)
		}

	case MOD_IS_FILE:
		mod.Call = func(_ Locals, in string) string {
			stat, err := os.Stat(in)
			if os.IsNotExist(err) {
				return boolToString(false)
			}
			return boolToString(stat.Mode().IsRegular())
		}

	case MOD_IS_DIR:
		mod.Call = func(_ Locals, in string) string {
			stat, err := os.Stat(in)
			if os.IsNotExist(err) {
				return boolToString(false)
			}
			return boolToString(stat.IsDir())
		}

	case MOD_IS_SYMLINK:
		mod.Call = func(_ Locals, in string) string {
			stat, err := os.Lstat(in)
			if os.IsNotExist(err) {
				return boolToString(false)
			}
			return boolToString(stat.Mode()&fs.ModeSymlink != 0)
		}

	case MOD_LENGTH:
		mod.Call = func(_ Locals, in string) string {
			return fmt.Sprint(len(in))
		}

	case MOD_QUOTED:
		mod.Call = func(_ Locals, in string) string {
			return fmt.Sprintf("%q", in)
		}

	case MOD_TRIM:
		mod.Call = func(locals Locals, in string) string {
			for _, arg := range mod.Args {
				value, err := arg.Eval(locals)
				if err != nil {
					return in
				}
				in = strings.Trim(in, value)
			}
			return in
		}

	case MOD_TRIM_PREFIX:
		mod.Call = func(locals Locals, in string) string {
			for _, arg := range mod.Args {
				value, err := arg.Eval(locals)
				if err != nil {
					return in
				}
				in = strings.TrimPrefix(in, value)
			}
			return in
		}

	case MOD_TRIM_SUFFIX:
		mod.Call = func(locals Locals, in string) string {
			for _, arg := range mod.Args {
				value, err := arg.Eval(locals)
				if err != nil {
					return in
				}
				in = strings.TrimSuffix(in, value)
			}
			return in
		}

	case MOD_PAD:
		switch len(mod.Args) {
		case 1:
			mod.Call = func(locals Locals, in string) string {
				value, err := mod.Args[0].Eval(locals)
				if err != nil {
					return in
				}
				number, err := strconv.Atoi(value)
				if err != nil {
					return in
				}
				return padString(in, number, number)
			}
		case 2:
			mod.Call = func(locals Locals, in string) string {
				value_1, err := mod.Args[0].Eval(locals)
				if err != nil {
					return in
				}
				value_2, err := mod.Args[1].Eval(locals)
				if err != nil {
					return in
				}
				number_1, err := strconv.Atoi(value_1)
				if err != nil {
					return in
				}
				number_2, err := strconv.Atoi(value_2)
				if err != nil {
					return in
				}
				return padString(in, number_1, number_2)
			}
		}

	case MOD_PAD_LEFT:
		if len(mod.Args) > 0 {
			mod.Call = func(locals Locals, in string) string {
				value, err := mod.Args[0].Eval(locals)
				if err != nil {
					return in
				}
				number, err := strconv.Atoi(value)
				if err != nil {
					return in
				}
				return padString(in, number, 0)
			}
		}

	case MOD_PAD_RIGHT:
		if len(mod.Args) > 0 {
			mod.Call = func(locals Locals, in string) string {
				value, err := mod.Args[0].Eval(locals)
				if err != nil {
					return in
				}
				number, err := strconv.Atoi(value)
				if err != nil {
					return in
				}
				return padString(in, 0, number)
			}
		}

	case MOD_HAS_PREFIX:
		if len(mod.Args) > 0 {
			mod.Call = func(locals Locals, in string) string {
				value, err := mod.Args[0].Eval(locals)
				if err != nil {
					return in
				}
				return boolToString(strings.HasPrefix(in, value))
			}
		}

	case MOD_HAS_SUFFIX:
		if len(mod.Args) > 0 {
			mod.Call = func(locals Locals, in string) string {
				value, err := mod.Args[0].Eval(locals)
				if err != nil {
					return in
				}
				return boolToString(strings.HasSuffix(in, value))
			}
		}

	case MOD_SLICE:
		if len(mod.Args) > 1 {
			mod.Call = func(locals Locals, in string) string {
				value_1, err := mod.Args[0].Eval(locals)
				if err != nil {
					return in
				}
				value_2, err := mod.Args[1].Eval(locals)
				if err != nil {
					return in
				}
				number_1, err := strconv.Atoi(value_1)
				if err != nil {
					return in
				}
				number_2, err := strconv.Atoi(value_2)
				if err != nil {
					return in
				}
				return sliceString(in, number_1, number_2)
			}
		}

	case MOD_REVERSE:
		// https://stackoverflow.com/a/10030772
		mod.Call = func(_ Locals, in string) string {
			runes := []rune(in)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}

	case MOD_INVERT:
		// https://stackoverflow.com/a/38234154
		mod.Call = func(_ Locals, in string) string {
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

	default:
		return mod, fmt.Errorf(
			"unknown string expansion modifier %q with parameters %q",
			mod.Name,
			mod.Args,
		)
	}

	return mod, nil
}

func boolToString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func padString(in string, left, right int) string {
	return strings.Repeat(" ", max(0, left)) + in + strings.Repeat(" ", max(0, right))
}

func sliceString(in string, from, to int) string {
	if from < 0 {
		from = max(0, len(in)+from)
	}
	if to < 0 {
		to = max(0, len(in)+to)
	}
	if from >= len(in) {
		return ""
	}
	return in[from:min(to, len(in))]
}
