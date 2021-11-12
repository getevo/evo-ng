package generic

import (
	"fmt"
	"testing"
	"time"
)

func TestGeneric(t *testing.T) {
	var string = "Hello World"

	fmt.Printf("Parse(%+v).String() => %+v \n",string,Parse(string).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n",string,Parse(string).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n",string,Parse(string).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n",string,Parse(string).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n",string,Parse(string).Bool())
	fmt.Printf("Parse(%+v).Uint64() => %+v \n","12222222222",Parse("12222222222").Uint64())
	fmt.Printf("Parse(%+v).Int64() => %+v \n","-12222222222",Parse("-12222222222").Int64())

	fmt.Printf("Parse(%+v).String() => %+v \n",nil,Parse(nil).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n",nil,Parse(nil).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n",nil,Parse(nil).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n",nil,Parse(nil).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n",string,Parse(nil).Bool())

	fmt.Printf("Parse(%+v).String() =>  %+v \n",6.3,Parse(6.3).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n",6.3,Parse(6.3).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n",6.3,Parse(6.3).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n",6.3,Parse(6.3).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n",string,Parse(6.3).Bool())

	fmt.Printf("Parse(%+v).String() => %+v \n",78,Parse(78).String())
	fmt.Printf("Parse(%+v).Int() => %+v \n",78,Parse(78).Int())
	fmt.Printf("Parse(%+v).Float() => %+v \n",78,Parse(78).Float())
	fmt.Printf("Parse(%+v).IsNil() => %+v \n",78,Parse(78).IsNil())
	fmt.Printf("Parse(%+v).Bool() => %+v \n",78,Parse(78).Bool())
	fmt.Printf("Parse(%+v).Uint64() => %+v \n",78,Parse(78).Uint64())
	fmt.Printf("Parse(%+v).Int64() => %+v \n",78,Parse(78).Int64())

	fmt.Printf("Parse(%+v).Bool() => %+v \n","Yes",Parse("Yes").Bool())
	fmt.Printf("Parse(%+v).Bool() => %+v \n","True",Parse("True").Bool())
	fmt.Printf("Parse(%+v).Bool() => %+v \n",1,Parse(1).Bool())


	fmt.Printf("Parse(%+v).Bool() => %+v \n",1,Parse(1).Bool())
	fmt.Printf("Parse(%+v).Bool() => %+v \n",1,Parse(1).Bool())



	var date time.Time
	var err error
	var duration time.Duration
	date,err = Parse("10/13/2021").Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n","10/13/2021",date,err )

	date,err = Parse("03.31.2014").Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n","03.31.2014", date,err)

	date,err = Parse(time.Now().Unix()).Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n",time.Now().Unix(), date,err)

	date,err = Parse(time.Now().UnixNano()).Time()
	fmt.Printf("Parse(%+v).Time() => %+v,%+v \n",time.Now().UnixNano(), date,err)


	duration,err = Parse("6h12m16s").Duration()
	fmt.Printf("Parse(%+v).Duration() => %+v,%+v \n","6h12m16s", duration,err)

	duration,err = Parse(6*time.Hour).Duration()
	fmt.Printf("Parse(%+v).Duration() => %+v minutes,%+v \n",6*time.Hour, duration.Minutes(),err)
}

func TestString(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",String(sample))
	fmt.Println(&sample,"=>",String(&sample))

	fmt.Println(array,"=>",String(array))
	fmt.Println(object,"=>",String(&object))

	fmt.Println(123,"=>",String(123))


}