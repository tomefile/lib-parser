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
	MOD_REVERSE,
	MOD_INVERT,
}

type StringModifier struct {
	Name ModifierName
	Args []string
	Call func(string) string
}

func (modifier StringModifier) String() string {
	if len(modifier.Args) == 0 {
		return string(modifier.Name)
	}

	return fmt.Sprintf("%s %s", modifier.Name, strings.Join(modifier.Args, " "))
}

func (modifier StringModifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s(%s)", modifier.Name, strings.Join(modifier.Args, " ")))
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

func GetModifier(name ModifierName, args []string) (StringModifier, error) {
	mod := StringModifier{
		Name: name,
		Args: args,
		Call: nil,
	}

	switch mod.Name {

	case MOD_NOT:
		mod.Call = func(in string) string {
			switch in {
			case boolToString(false), "false", "FALSE":
				return boolToString(true)
			default:
				return boolToString(false)
			}
		}

	case MOD_TO_LOWER:
		mod.Call = strings.ToLower

	case MOD_TO_UPPER:
		mod.Call = strings.ToUpper

	case MOD_TO_SNAKE:
		mod.Call = strcase.ToSnake

	case MOD_TO_KEBAB:
		mod.Call = strcase.ToKebab

	case MOD_TO_CAMEL:
		mod.Call = strcase.ToLowerCamel

	case MOD_TO_PASCAL:
		mod.Call = strcase.ToCamel

	case MOD_TO_DELIMITED:
		if len(mod.Args) > 0 {
			mod.Call = func(in string) string {
				return strcase.ToDelimited(in, uint8(mod.Args[0][0]))
			}
		}

	case MOD_DOES_EXIST:
		mod.Call = func(in string) string {
			_, err := os.Stat(in)
			return boolToString(!os.IsNotExist(err))
		}

	case MOD_IS_EMPTY:
		mod.Call = func(in string) string {
			return boolToString(len(in) == 0)
		}

	case MOD_IS_FILE:
		mod.Call = func(in string) string {
			stat, err := os.Stat(in)
			if os.IsNotExist(err) {
				return boolToString(false)
			}
			return boolToString(stat.Mode().IsRegular())
		}

	case MOD_IS_DIR:
		mod.Call = func(in string) string {
			stat, err := os.Stat(in)
			if os.IsNotExist(err) {
				return boolToString(false)
			}
			return boolToString(stat.IsDir())
		}

	case MOD_IS_SYMLINK:
		mod.Call = func(in string) string {
			stat, err := os.Lstat(in)
			if os.IsNotExist(err) {
				return boolToString(false)
			}
			return boolToString(stat.Mode()&fs.ModeSymlink != 0)
		}

	case MOD_LENGTH:
		mod.Call = func(in string) string {
			return fmt.Sprint(len(in))
		}

	case MOD_QUOTED:
		mod.Call = func(in string) string {
			return fmt.Sprintf("%q", in)
		}

	case MOD_TRIM:
		mod.Call = func(in string) string {
			for _, arg := range mod.Args {
				in = strings.Trim(in, arg)
			}
			return in
		}

	case MOD_TRIM_PREFIX:
		mod.Call = func(in string) string {
			for _, arg := range mod.Args {
				in = strings.TrimPrefix(in, arg)
			}
			return in
		}

	case MOD_TRIM_SUFFIX:
		mod.Call = func(in string) string {
			for _, arg := range mod.Args {
				in = strings.TrimSuffix(in, arg)
			}
			return in
		}

	case MOD_PAD:
		switch len(mod.Args) {
		case 1:
			number, err := strconv.Atoi(mod.Args[0])
			if err == nil {
				padding := strings.Repeat(" ", max(0, number))
				mod.Call = func(in string) string {
					return padding + in + padding
				}
			}
		case 2:
			number_1, err_1 := strconv.Atoi(mod.Args[0])
			number_2, err_2 := strconv.Atoi(mod.Args[1])
			if err_1 == nil && err_2 == nil {
				mod.Call = func(in string) string {
					return strings.Repeat(" ", max(0, number_1)) +
						in +
						strings.Repeat(" ", max(0, number_2))
				}
			}
		}

	case MOD_PAD_LEFT:
		if len(mod.Args) > 0 {
			number, err := strconv.Atoi(mod.Args[0])
			if err == nil {
				padding := strings.Repeat(" ", max(0, number))
				mod.Call = func(in string) string {
					return padding + in
				}
			}
		}

	case MOD_PAD_RIGHT:
		if len(mod.Args) > 0 {
			number, err := strconv.Atoi(mod.Args[0])
			if err == nil {
				padding := strings.Repeat(" ", max(0, number))
				mod.Call = func(in string) string {
					return in + padding
				}
			}
		}

	case MOD_REVERSE:
		// https://stackoverflow.com/a/10030772
		mod.Call = func(in string) string {
			runes := []rune(in)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}

	case MOD_INVERT:
		// https://stackoverflow.com/a/38234154
		mod.Call = func(in string) string {
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
