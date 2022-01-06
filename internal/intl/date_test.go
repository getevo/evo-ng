package intl

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	fmt.Println(Date("2006/06/08").Time())
	fmt.Println(Date("2006/06/08 10:45:30").Time())
	fmt.Println(Date(2006, 6, 8).Time())
	fmt.Println(Date(2006, 6, 8, 10, 45, 30).Time())
	fmt.Println(Date(2006, 6, 8, 10, 45, 30, 20).Time())
	var l, _ = time.LoadLocation("Asia/Tehran")
	fmt.Println(Date(2006, 6, 8, 10, 46, 30, 20, *l).Time())
	fmt.Println(Date(2006, 6, 8, 10, 46, 30, 20, l).Time())
	fmt.Println(Date().Time())
	//t.Error() // to indicate test failed
}
