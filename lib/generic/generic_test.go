package generic

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGeneric(t *testing.T) {
	var string = "Hello World"

	fmt.Printf("Parse(%+v).SizeInBytes() => %+v \n", "4mb", Parse("4mb").SizeInBytes())
	fmt.Printf("Parse(%+v).SizeInBytes() => %+v \n", "4GB", Parse("4GB").SizeInBytes())
	fmt.Printf("Parse(%+v).SizeInBytes() => %+v \n", "10", Parse("10").SizeInBytes())
	fmt.Printf("Parse(%+v).SizeInBytes() => %+v \n", "10B", Parse("10B").SizeInBytes())

	fmt.Printf("Parse(%+v).ByteCount() => %+v \n", "4294967296 ", Parse("4294967296").ByteCount())
	fmt.Printf("Parse(%+v).ByteCount() => %+v \n", "4398046511104", Parse("4398046511104").ByteCount())
	fmt.Printf("Parse(%+v).ByteCount() => %+v \n", "10737418240", Parse("10737418240").ByteCount())

	fmt.Printf("Parse(%+v).String() => %+v \n", string, Parse(string).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n", string, Parse(string).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n", string, Parse(string).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n", string, Parse(string).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n", string, Parse(string).Bool())
	fmt.Printf("Parse(%+v).Uint64() => %+v \n", "12222222222", Parse("12222222222").Uint64())
	fmt.Printf("Parse(%+v).Int64() => %+v \n", "-12222222222", Parse("-12222222222").Int64())

	fmt.Printf("Parse(%+v).String() => %+v \n", nil, Parse(nil).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n", nil, Parse(nil).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n", nil, Parse(nil).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n", nil, Parse(nil).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n", string, Parse(nil).Bool())

	fmt.Printf("Parse(%+v).String() =>  %+v \n", 6.3, Parse(6.3).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n", 6.3, Parse(6.3).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n", 6.3, Parse(6.3).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n", 6.3, Parse(6.3).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n", string, Parse(6.3).Bool())

	fmt.Printf("Parse(%+v).String() => %+v \n", 78, Parse(78).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n", 78, Parse(78).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n", 78, Parse(78).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n", 78, Parse(78).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n", 78, Parse(78).Bool())
	fmt.Printf("Parse(%+v).Uint64() => %+v \n", 78, Parse(78).Uint64())
	fmt.Printf("Parse(%+v).Int64() => %+v \n", 78, Parse(78).Int64())

	fmt.Printf("Parse(%+v).Bool() => %+v \n", "Yes", Parse("Yes").Bool())
	fmt.Printf("Parse(%+v).Bool() => %+v \n", "True", Parse("True").Bool())
	fmt.Printf("Parse(%+v).Bool() => %+v \n", 1, Parse(1).Bool())

	fmt.Printf("Parse(%+v).Bool() => %+v \n", 1, Parse(1).Bool())
	fmt.Printf("Parse(%+v).Bool() => %+v \n", 1, Parse(1).Bool())

	var date time.Time
	var err error
	var duration time.Duration
	date, err = Parse("10/13/2021").Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n", "10/13/2021", date, err)

	date, err = Parse("03.31.2014").Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n", "03.31.2014", date, err)

	date, err = Parse(time.Now().Unix()).Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n", time.Now().Unix(), date, err)

	date, err = Parse(time.Now().UnixNano()).Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n", time.Now().UnixNano(), date, err)

	duration, err = Parse("6h12m16s").Duration()
	fmt.Printf("Parse(%+v).Duration() => %+v,%+v \n", "6h12m16s", duration, err)

	duration, err = Parse(6 * time.Hour).Duration()
	fmt.Printf("Parse(%+v).Duration() => %+v minutes,%+v \n", 6*time.Hour, duration.Minutes(), err)
}

func TestString(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X: sample,
	}
	var array = []string{"a", "b", "c"}
	fmt.Println(sample, "=>", ToString(sample))
	fmt.Println(&sample, "=>", ToString(&sample))

	fmt.Println(array, "=>", ToString(array))
	fmt.Println(object, "=>", ToString(&object))

	fmt.Println(123, "=>", ToString(123))

}
func TestTypeOf(t *testing.T) {
	var s = "string"
	fmt.Println(TypeOf(s).Indirect().Is(reflect.String))
	fmt.Println(TypeOf(&s).Is(reflect.String))
	fmt.Println(TypeOf(&s).Indirect().Is(reflect.String))
	fmt.Println(TypeOf(&s).Is("*string"))
	fmt.Println(TypeOf(&s).Is(&s))

}
