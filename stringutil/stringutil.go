package stringutil

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var rxSpace = regexp.MustCompile(`[\s\-]+`)
var rxHexadecimal = regexp.MustCompile(`^[0-9a-fA-F]+$`)
var DefaultThousandsSeparator = `,`
var DefaultDecimalSeparator = `.`

var TimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC1123,
	time.RFC1123Z,
	time.Kitchen,
	`2006-01-02 15:04:05 -0700 MST`,
	`2006-01-02 15:04:05Z07:00`,
	`2006-01-02 15:04:05`,
	`2006-01-02 15:04`,
	`2006-01-02`,
}

type SiPrefix int

const (
	None  SiPrefix = 0
	Kilo           = 1
	Mega           = 2
	Giga           = 3
	Tera           = 4
	Peta           = 5
	Exa            = 6
	Zetta          = 7
	Yotta          = 8
)

func (self SiPrefix) String() string {
	switch self {
	case Kilo:
		return `K`
	case Mega:
		return `M`
	case Giga:
		return `G`
	case Tera:
		return `T`
	case Peta:
		return `P`
	case Exa:
		return `E`
	case Zetta:
		return `Z`
	case Yotta:
		return `Y`
	default:
		return ``
	}
}

type ConvertType int

const (
	Invalid ConvertType = iota
	String
	Boolean
	Float
	Integer
	Time
)

func IsInteger(in interface{}) bool {
	inV := reflect.ValueOf(in)

	switch inV.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true

	default:
		if asStr, err := ToString(in); err == nil {
			if _, err := strconv.Atoi(asStr); err == nil {
				return true
			}
		}
	}

	return false
}

func IsFloat(in interface{}) bool {
	inV := reflect.ValueOf(in)

	switch inV.Kind() {
	case reflect.Float32, reflect.Float64:
		return true

	default:
		if asStr, err := ToString(in); err == nil {
			if _, err := strconv.ParseFloat(asStr, 64); err == nil {
				return true
			}
		}
	}

	return false
}

func IsNumeric(in interface{}) bool {
	return IsFloat(in)
}

func IsBoolean(in string) bool {
	in = strings.ToLower(in)

	return (IsBooleanTrue(in) || IsBooleanFalse(in))
}

func IsBooleanTrue(in string) bool {
	in = strings.ToLower(in)

	switch in {
	case `true`, `yes`, `on`:
		return true
	}

	return false
}

func IsBooleanFalse(in string) bool {
	in = strings.ToLower(in)

	switch in {
	case `false`, `no`, `off`:
		return true
	}

	return false
}

func IsTime(in string) bool {
	if DetectTimeFormat(in) != `` {
		return true
	}

	return false
}

func DetectTimeFormat(in string) string {
	for _, layout := range TimeFormats {
		if _, err := time.Parse(layout, in); err == nil {
			return layout
		}
	}

	return ``
}

func ToString(in interface{}) (string, error) {
	switch reflect.TypeOf(in).Kind() {
	case reflect.Int:
		return strconv.FormatInt(int64(in.(int)), 10), nil
	case reflect.Int8:
		return strconv.FormatInt(int64(in.(int8)), 10), nil
	case reflect.Int16:
		return strconv.FormatInt(int64(in.(int16)), 10), nil
	case reflect.Int32:
		return strconv.FormatInt(int64(in.(int32)), 10), nil
	case reflect.Int64:
		return strconv.FormatInt(in.(int64), 10), nil
	case reflect.Uint:
		return strconv.FormatUint(uint64(in.(uint)), 10), nil
	case reflect.Uint8:
		return strconv.FormatUint(uint64(in.(uint8)), 10), nil
	case reflect.Uint16:
		return strconv.FormatUint(uint64(in.(uint16)), 10), nil
	case reflect.Uint32:
		return strconv.FormatUint(uint64(in.(uint32)), 10), nil
	case reflect.Uint64:
		return strconv.FormatUint(in.(uint64), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(float64(in.(float32)), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(in.(float64), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(in.(bool)), nil
	case reflect.String:
		return in.(string), nil
	case reflect.Struct:
		if stringFn := reflect.ValueOf(in).MethodByName(`String`); stringFn != reflect.Zero(reflect.TypeOf(stringFn)) {
			return fmt.Sprintf("%v", in), nil
		}

		fallthrough
	default:
		return "", fmt.Errorf("Unable to convert type '%T' to string", in)
	}
}

func ToStringSlice(in interface{}) ([]string, error) {
	values := make([]string, 0)

	if in != nil {
		if v, ok := in.([]string); ok {
			return v, nil
		}

		inV := reflect.ValueOf(in)

		if inV.IsValid() {
			if inV.Kind() == reflect.Ptr {
				inV = inV.Elem()
			}

			if inV.IsValid() {
				switch inV.Kind() {
				case reflect.Array, reflect.Slice:
					for i := 0; i < inV.Len(); i++ {
						if indexV := inV.Index(i); indexV.IsValid() {
							if v, err := ToString(indexV.Interface()); err == nil {
								values = append(values, v)
							} else {
								return nil, err
							}
						} else {
							return nil, fmt.Errorf("Element %d in slice is invalid", i)
						}
					}

				default:
					if v, err := ToString(in); err == nil {
						values = append(values, v)
					} else {
						return nil, err
					}
				}
			} else {
				return nil, fmt.Errorf("Cannot parse value pointed to by given value.")
			}

		} else {
			return nil, fmt.Errorf("Cannot parse given value.")
		}
	}

	return values, nil
}

func ToByteString(in interface{}, formatString ...string) (string, error) {
	if asBytes, err := ConvertToInteger(in); err == nil {
		for i := 0; i < 9; i++ {
			if converted := (float64(asBytes) / math.Pow(1024, float64(i))); converted < 1024 {
				prefix := SiPrefix(i)
				f := `%g`

				if len(formatString) > 0 {
					f = formatString[0]
				}

				return fmt.Sprintf(f+"%sB", converted, prefix.String()), nil
			}
		}

		return fmt.Sprintf("%fB", asBytes), nil
	} else {
		return ``, err
	}
}

func GetSiPrefix(input string) (SiPrefix, error) {
	switch input {
	case "", "b", "B":
		return None, nil
	case "k", "K":
		return Kilo, nil
	case "m", "M":
		return Mega, nil
	case "g", "G":
		return Giga, nil
	case "t", "T":
		return Tera, nil
	case "p", "P":
		return Peta, nil
	case "e", "E":
		return Exa, nil
	case "z", "Z":
		return Zetta, nil
	case "y", "Y":
		return Yotta, nil
	default:
		return None, fmt.Errorf("Unrecognized SI unit '%s'", input)
	}
}

func ToBytes(input string) (float64, error) {
	//  handle -ibibyte suffixes like KiB, GiB
	if strings.HasSuffix(input, "ib") || strings.HasSuffix(input, "iB") {
		input = input[0 : len(input)-2]

		//  handle input that puts the 'B' in the suffix; e.g.: Kb, GB
	} else if len(input) > 2 && IsInteger(string(input[len(input)-3])) && (input[len(input)-1] == 'b' || input[len(input)-1] == 'B') {
		input = input[0 : len(input)-1]
	}

	if prefix, err := GetSiPrefix(string(input[len(input)-1])); err == nil {
		if v, err := strconv.ParseFloat(input[0:len(input)-1], 64); err == nil {
			return v * math.Pow(1024, float64(prefix)), nil
		} else {
			return 0, err
		}
	} else {
		if v, err := strconv.ParseFloat(input, 64); err == nil {
			return v, nil
		} else {
			return 0, fmt.Errorf("Unrecognized input string '%s'", input)
		}
	}
}

func ConvertTo(toType ConvertType, inI interface{}) (interface{}, error) {
	if in, err := ToString(inI); err == nil {
		switch toType {
		case Float:
			if inS, ok := inI.(string); ok {
				if inS == `` {
					return float64(0), nil
				}
			}

			return strconv.ParseFloat(in, 64)
		case Integer:
			if inTime, ok := inI.(time.Time); ok {
				return int64(inTime.UnixNano()), nil
			} else if inS, ok := inI.(string); ok {
				if inS == `` {
					return int64(0), nil
				} else if layout := DetectTimeFormat(inS); layout != `` {
					if tm, err := time.Parse(layout, inS); err == nil {
						return tm.UnixNano(), nil
					} else {
						return nil, err
					}
				}
			}

			return strconv.ParseInt(in, 10, 64)
		case Boolean:
			if IsBooleanTrue(in) {
				return true, nil
			} else if IsBooleanFalse(in) {
				return false, nil
			} else {
				return nil, fmt.Errorf("Cannot convert '%s' into a boolean value", in)
			}
		case Time:
			inS := fmt.Sprintf("%v", in)

			switch inS {
			case `now`:
				return time.Now(), nil
			default:
				// handle time zero values
				tmS := strings.Map(func(r rune) rune {
					switch r {
					case '-', ':', ' ', 'T', 'Z':
						return '0'
					}

					return r
				}, inS)

				if v, err := strconv.ParseInt(tmS, 10, 64); err == nil && v == 0 {
					return time.Time{}, nil
				}

				for _, format := range TimeFormats {
					if tm, err := time.Parse(format, strings.TrimSpace(in)); err == nil {
						return tm, nil
					}
				}

				return nil, fmt.Errorf("Cannot convert '%s' into a date/time value", in)
			}

		default:
			return in, nil
		}
	} else {
		return nil, err
	}
}

func ConvertToInteger(in interface{}) (int64, error) {
	if v, err := ConvertTo(Integer, in); err == nil {
		return v.(int64), nil
	} else {
		return int64(0), err
	}
}

func ConvertToFloat(in interface{}) (float64, error) {
	if v, err := ConvertTo(Float, in); err == nil {
		return v.(float64), nil
	} else {
		return float64(0.0), err
	}
}

func ConvertToString(in interface{}) (string, error) {
	if v, err := ConvertTo(String, in); err == nil {
		return v.(string), nil
	} else {
		return ``, err
	}
}

func ConvertToBool(in interface{}) (bool, error) {
	if v, err := ConvertTo(Boolean, in); err == nil {
		return v.(bool), nil
	} else {
		return false, err
	}
}

func ConvertToTime(in interface{}) (time.Time, error) {
	switch in.(type) {
	case time.Time:
		return in.(time.Time), nil
	default:
		if v, err := ConvertTo(Time, in); err == nil {
			return v.(time.Time), nil
		} else {
			return time.Time{}, err
		}
	}
}

func Autotype(in interface{}) interface{} {
	for _, ctype := range []ConvertType{
		Boolean,
		Time,
		Integer,
		Float,
		String,
	} {
		if value, err := ConvertTo(ctype, in); err == nil {
			return value
		}
	}

	return in
}

func IsSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		}

		return true
	}

	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}

	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

func TokenizeFunc(in string, tokenizer func(rune) bool, partfn func(part string) []string) []string {
	// split on word-separating characters (and discard them), or on capital
	// letters (preserving them)
	parts := strings.FieldsFunc(in, tokenizer)
	out := make([]string, 0)

	for _, part := range parts {
		partOut := partfn(part)

		if partOut != nil {
			for _, v := range partOut {
				if v != `` {
					out = append(out, v)
				}
			}
		}
	}

	return out
}

func Camelize(in string) string {
	return strings.Join(TokenizeFunc(in, IsSeparator, func(part string) []string {
		part = strings.TrimSpace(part)
		part = strings.Title(part)
		return []string{part}
	}), ``)
}

func Underscore(in string) string {
	in = rxSpace.ReplaceAllString(in, `_`)
	out := make([]rune, 0)
	runes := []rune(in)

	sepfn := func(i int) bool {
		return i >= 0 && i < len(runes) && unicode.IsLower(runes[i])
	}

	for i, r := range runes {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)

			if i > 0 && runes[i-1] != '_' && (sepfn(i-1) || sepfn(i+1)) {
				out = append(out, '_')
			}
		}

		out = append(out, r)
	}

	return string(out)
}

// Returns whether the letters (Unicode Catgeory 'L') in a given string are
// homogenous in case (all upper-case or all lower-case).
//
func IsMixedCase(in string) bool {
	var hasLower bool
	var hasUpper bool

	for _, c := range in {
		if unicode.IsLetter(c) {
			if unicode.IsLower(c) {
				hasLower = true

				if hasUpper {
					return true
				}
			} else if unicode.IsUpper(c) {
				hasUpper = true

				if hasLower {
					return true
				}
			}
		}
	}

	return false
}

// Returns whether the given string is a hexadecimal number. If the string is
// prefixed with "0x", the prefix is removed first. If length is greater than 0,
// the length of the input (excluding prefix) is checked as well.
//
func IsHexadecimal(in string, length int) bool {
	in = strings.TrimPrefix(in, `0x`)

	if IsMixedCase(in) {
		return false
	}

	if rxHexadecimal.MatchString(in) {
		if length <= 0 {
			return true
		} else if len(in) == length {
			return true
		}
	}

	return false
}


func Thousandify(in interface{}, separator string, decimal string) string {
	if separator == `` {
		separator = DefaultThousandsSeparator
	}

	if decimal == `` {
		decimal = DefaultDecimalSeparator
	}

	if inStr, err := ToString(in); err == nil {
		if IsNumeric(in) {
			var buffer []rune

			lastIndexBeforeDecimal := strings.Index(inStr, decimal) - 1
			decimalAndAfter := strings.Index(inStr, decimal)

			if lastIndexBeforeDecimal < 0 {
				lastIndexBeforeDecimal = len(inStr) - 1
			}

			j := 0

			for i := lastIndexBeforeDecimal; i >= 0; i-- {
				j++
				buffer = append([]rune{rune(inStr[i])}, buffer...)

				if j == 3 && i > 0 && !(i == 1 && inStr[0] == '-') {
					buffer = append([]rune(separator), buffer...)
					j = 0
				}
			}

			if decimalAndAfter >= 0 {
				for _, r := range inStr[decimalAndAfter:] {
					buffer = append(buffer, rune(r))
				}
			}

			return string(buffer[:])
		}else{
			return inStr
		}
	}else{
		return ``
	}
}
