package generic

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"gopkg.in/yaml.v3"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Invalid reflect.Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)

// Parse parse input
//  @param i
//  @return Value
func Parse(i interface{}) Value {
	return Value{
		Input: i,
	}
}

// Type lib structure to keep input and its type
type Type struct {
	input interface{}
	iType reflect.Type
}

// Value wraps over interface
type Value struct {
	Input interface{}
}

// IsNil returns if the value is nil
//  @receiver v
//  @return bool
func (v Value) IsNil() bool {
	return !reflect.ValueOf(v.Input).IsValid()
}

func (v Value) direct() interface{} {
	ref := reflect.ValueOf(v.Input)
	if !ref.IsValid() {
		return nil
	}
	if ref.Type().Kind() == reflect.Ptr {
		return ref.Elem().Interface()
	}
	return v.Input
}

// ParseJSON parse json value into struct
//  @receiver v
//  @param in
//  @return error
func (v Value) ParseJSON(in interface{}) error {
	var value = v.direct()
	return json.Unmarshal([]byte(fmt.Sprint(value)), in)
}

// String return value as string
//  @receiver v
//  @return string
func (v Value) String() string {
	var value = v.direct()
	return fmt.Sprint(value)
}

// Int return value as integer
//  @receiver v
//  @return int
func (v Value) Int() int {
	i, _ := strconv.Atoi(v.String())
	return i
}

// Uint64 return value as uint64
//  @receiver v
//  @return uint64
func (v Value) Uint64() uint64 {
	i, _ := strconv.ParseUint(v.String(), 0, 64)
	return i
}

// Int64 return value as int64
//  @receiver v
//  @return int64
func (v Value) Int64() int64 {
	i, _ := strconv.ParseInt(v.String(), 0, 64)
	return i
}

// Float return value as float64
//  @receiver v
//  @return float64
func (v Value) Float() float64 {
	i, _ := strconv.ParseFloat(v.String(), 64)
	return i
}

// Bool return value as bool
//  @receiver v
//  @return bool
func (v Value) Bool() bool {
	var s = strings.ToLower(v.String())
	if len(s) > 0 {
		if s == "1" || s == "true" || s == "yes" {
			return true
		}
	}
	return false
}

// Time return value as time.Time
//  @receiver v
//  @return time.Time
//  @return error
func (v Value) Time() (time.Time, error) {
	return dateparse.ParseAny(v.String())
}

// Duration return value as time.Duration
//  @receiver v
//  @return time.Duration
//  @return error
func (v Value) Duration() (time.Duration, error) {
	return time.ParseDuration(v.String())
}

// ToString cast anything to string
//  @param v
//  @return string
func ToString(v interface{}) string {
	ref := reflect.ValueOf(v)
	if !ref.IsValid() {
		return ""
	}
	if ref.Type().Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	switch ref.Kind() {
	case reflect.String:
		return ref.Interface().(string)
	case reflect.Struct, reflect.Slice:
		if v, ok := ref.Interface().([]byte); ok {
			return string(v)
		}
		if v, ok := ref.Interface().(Value); ok {
			return v.String()
		}
		b, _ := json.Marshal(ref.Interface())
		return string(b)
	default:
		return fmt.Sprint(ref.Interface())
	}
}

// TypeOf return type of input
//  @param input
//  @return *Type
func TypeOf(input interface{}) *Type {
	var el = Type{
		input: input,
		iType: reflect.TypeOf(input),
	}
	return &el
}

// Is checks if input and given type are equal
//  @receiver t
//  @param input
//  @return bool
func (t *Type) Is(input interface{}) bool {
	switch v := input.(type) {
	case reflect.Kind:
		return v == t.iType.Kind()
	case string:
		return t.iType.String() == v
	}
	return TypeOf(input).iType.String() == t.iType.String()
}

// Indirect get type of object considering if it is pointer
//  @receiver t
//  @return *Type
func (t *Type) Indirect() *Type {
	if t.Is(reflect.Ptr) {
		var el = Type{
			input: reflect.ValueOf(t.input).Elem().Interface(),
		}
		el.iType = reflect.TypeOf(el.input)
		return &el
	}
	return t
}

var sizeRegex = regexp.MustCompile(`(?m)^(\d+)\s*([kmgte]{0,1}b){0,1}$`)

func (v Value) SizeInBytes() uint64 {
	var s = strings.ToLower(strings.TrimSpace(fmt.Sprint(v.Input)))
	var match = sizeRegex.FindAllStringSubmatch(s, 1)
	if len(match) == 1 && len(match[0]) == 3 {
		var base, _ = strconv.ParseUint(match[0][1], 10, 64)
		switch match[0][2] {
		case "kb":
			base *= 1024
		case "", "mb":
			base *= 1024 * 1024 * 1024
		case "gb":
			base *= 1024 * 1024 * 1024 * 1024
		case "tb":
			base *= 1024 * 1024 * 1024 * 1024 * 1024
		case "eb":
			base *= 1024 * 1024 * 1024 * 1024 * 1024 * 1024
		}
		return base
	}
	return 0
}

func (v Value) ByteCount() string {
	var b, _ = strconv.ParseUint(strings.TrimSpace(fmt.Sprint(v.Input)), 10, 64)
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func (v *Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Input)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.Input)
}

func (v *Value) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(v.Input)
}

func (v *Value) UnmarshalYAML(data []byte) error {
	return yaml.Unmarshal(data, &v.Input)
}

func (v *Value) Scan(value interface{}) error {
	switch cast := value.(type) {
	case string:
		v.Input = cast
	case []byte:
		v.Input = string(cast)
	default:
		v.Input = cast
	}
	return nil
}

// Value return drive.Value value, implement driver.Valuer interface of gorm
func (v Value) Value() (driver.Value, error) {
	return ToString(v.Input), nil
}
