package generic

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"reflect"
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

func Parse(i interface{}) Value {
	return Value{
		Input: i,
	}
}

type Type struct {
	input interface{}
	iType reflect.Type
}

type Value struct {
	Input interface{}
}

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

func (v Value) String() string {
	var value = v.direct()
	return fmt.Sprint(value)
}

func (v Value) Int() int {
	i, _ := strconv.Atoi(v.String())
	return i
}

func (v Value) Uint64() uint64 {
	i, _ := strconv.ParseUint(v.String(), 0, 64)
	return i
}

func (v Value) Int64() int64 {
	i, _ := strconv.ParseInt(v.String(), 0, 64)
	return i
}

func (v Value) Float() float64 {
	i, _ := strconv.ParseFloat(v.String(), 64)
	return i
}

func (v Value) Bool() bool {
	var s = strings.ToLower(v.String())
	if len(s) > 0 {
		if s == "1" || s == "true" || s == "yes" {
			return true
		}
	}
	return false
}

func (v Value) Time() (time.Time, error) {
	return dateparse.ParseAny(v.String())
}

func (v Value) Duration() (time.Duration, error) {
	return time.ParseDuration(v.String())
}

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

func TypeOf(input interface{}) *Type {
	var el = Type{
		input: input,
		iType: reflect.TypeOf(input),
	}
	return &el
}

func (t *Type) Is(input interface{}) bool {
	switch v := input.(type) {
	case reflect.Kind:
		return v == t.iType.Kind()
	case string:
		return t.iType.String() == v
	}
	return TypeOf(input).iType.String() == t.iType.String()
}

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
