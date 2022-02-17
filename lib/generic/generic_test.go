package generic_test

import (
	"fmt"
	"github.com/getevo/evo-ng/lib/generic"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestGeneric(t *testing.T) {
	var string = "Hello World"

	assert.Equal(t, uint64(4294967296), generic.Parse("4mb").SizeInBytes())
	assert.Equal(t, uint64(4398046511104), generic.Parse("4GB").SizeInBytes())
	assert.Equal(t, uint64(10737418240), generic.Parse("10").SizeInBytes())
	assert.Equal(t, uint64(10), generic.Parse("10B").SizeInBytes())

	assert.Equal(t, generic.Parse(4294967296).ByteCount(), "4.0 GB")
	assert.Equal(t, generic.Parse("4398046511104").ByteCount(), "4.0 TB")
	assert.Equal(t, generic.Parse("10737418240").ByteCount(), "10.0 GB")

	assert.Equal(t, string, generic.Parse(string).String())
	assert.Equal(t, 0, generic.Parse(string).Int())
	assert.Equal(t, float64(0), generic.Parse(string).Float())
	assert.Equal(t, false, generic.Parse(string).IsNil())
	assert.Equal(t, false, generic.Parse(string).Bool())

	assert.Equal(t, uint64(12222222222), generic.Parse(12222222222).Uint64())
	assert.Equal(t, int64(-12222222222), generic.Parse("-12222222222").Int64())

	assert.Equal(t, "<nil>", generic.Parse(nil).String())
	assert.Equal(t, 0, generic.Parse(nil).Int())
	assert.Equal(t, float64(0), generic.Parse(nil).Float())
	assert.Equal(t, true, generic.Parse(nil).IsNil())
	assert.Equal(t, false, generic.Parse(nil).Bool())

	assert.Equal(t, "6.3", generic.Parse(6.3).String())
	assert.Equal(t, uint64(6), generic.Parse(6.3).Uint64())
	assert.Equal(t, int64(6), generic.Parse(6.3).Int64())
	assert.Equal(t, 6, generic.Parse(6.3).Int())
	assert.Equal(t, float64(6.3), generic.Parse(6.3).Float())
	assert.Equal(t, false, generic.Parse(6.3).IsNil())
	assert.Equal(t, false, generic.Parse(6.3).Bool())

	assert.Equal(t, "78", generic.Parse(78).String())
	assert.Equal(t, uint64(78), generic.Parse(78).Uint64())
	assert.Equal(t, int64(78), generic.Parse(78).Int64())
	assert.Equal(t, 78, generic.Parse(78).Int())
	assert.Equal(t, float64(78), generic.Parse(78).Float())
	assert.Equal(t, false, generic.Parse(78).IsNil())
	assert.Equal(t, false, generic.Parse(78).Bool())

	assert.Equal(t, true, generic.Parse("Yes").Bool())
	assert.Equal(t, true, generic.Parse("True").Bool())
	assert.Equal(t, true, generic.Parse("yes").Bool())
	assert.Equal(t, true, generic.Parse(1).Bool())

	var date time.Time
	var err error
	var duration time.Duration
	date, err = generic.Parse("10/13/2021").Time()
	assert.NoError(t, err)
	assert.Equal(t, "2021-10-13 00:00:00 +0000 UTC", date.String())

	date, err = generic.Parse("03.31.2014").Time()
	assert.NoError(t, err)
	assert.Equal(t, "2014-03-31 00:00:00 +0000 UTC", date.String())

	date, err = generic.Parse(time.Now().Unix()).Time()
	assert.NoError(t, err)
	assert.Equal(t, time.Now().Format("2006-01-02 15:04:05"), date.Format("2006-01-02 15:04:05"))

	date, err = generic.Parse(time.Now().UnixNano()).Time()
	assert.NoError(t, err)
	assert.Equal(t, time.Now().Format("2006-01-02 15:04:05"), date.Format("2006-01-02 15:04:05"))

	duration, err = generic.Parse("6h12m16s").Duration()
	assert.NoError(t, err)
	assert.Equal(t, "6h12m16s", duration.String())

	duration, err = generic.Parse(6 * time.Hour).Duration()
	assert.NoError(t, err)
	assert.Equal(t, "6h0m0s", duration.String())

}

func TestString(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X: sample,
	}
	var array = []string{"a", "b", "c"}
	fmt.Println(sample, "=>", generic.ToString(sample))
	fmt.Println(&sample, "=>", generic.ToString(&sample))

	fmt.Println(array, "=>", generic.ToString(array))
	fmt.Println(object, "=>", generic.ToString(&object))

	fmt.Println(123, "=>", generic.ToString(123))

}
func TestTypeOf(t *testing.T) {
	var s = "string"
	fmt.Println(generic.TypeOf(s).Indirect().Is(reflect.String))
	fmt.Println(generic.TypeOf(&s).Is(reflect.String))
	fmt.Println(generic.TypeOf(&s).Indirect().Is(reflect.String))
	fmt.Println(generic.TypeOf(&s).Is("*string"))
	fmt.Println(generic.TypeOf(&s).Is(&s))

}
