package hash

import (
	"fmt"
	"testing"
)




func TestFNV32(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",FNV32(sample))
	fmt.Println(&sample,"=>",FNV32(&sample))

	fmt.Println(array,"=>",FNV32(array))
	fmt.Println(object,"=>",FNV32(&object))

	fmt.Println(123,"=>",FNV32(123))
}


func TestFNV32a(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",FNV32a(sample))
	fmt.Println(&sample,"=>",FNV32a(&sample))

	fmt.Println(array,"=>",FNV32a(array))
	fmt.Println(object,"=>",FNV32a(&object))

	fmt.Println(123,"=>",FNV32a(123))
}


func TestFNV64(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",FNV64(sample))
	fmt.Println(&sample,"=>",FNV64(&sample))

	fmt.Println(array,"=>",FNV64(array))
	fmt.Println(object,"=>",FNV64(&object))

	fmt.Println(123,"=>",FNV64(123))
}


func TestFNV64a(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",FNV64a(sample))
	fmt.Println(&sample,"=>",FNV64a(&sample))

	fmt.Println(array,"=>",FNV64a(array))
	fmt.Println(object,"=>",FNV64a(&object))

	fmt.Println(123,"=>",FNV64a(123))
}


func TestMD5(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",MD5(sample))
	fmt.Println(&sample,"=>",MD5(&sample))

	fmt.Println(array,"=>",MD5(array))
	fmt.Println(object,"=>",MD5(&object))

	fmt.Println(123,"=>",MD5(123))
}


func TestSHA1(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",SHA1(sample))
	fmt.Println(&sample,"=>",SHA1(&sample))

	fmt.Println(array,"=>",SHA1(array))
	fmt.Println(object,"=>",SHA1(&object))

	fmt.Println(123,"=>",SHA1(123))
}


func TestSHA256(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",SHA256(sample))
	fmt.Println(&sample,"=>",SHA256(&sample))

	fmt.Println(array,"=>",SHA256(array))
	fmt.Println(object,"=>",SHA256(&object))

	fmt.Println(123,"=>",SHA256(123))
}


func TestSHA512(t *testing.T) {
	var sample = "Hello World"
	var object = struct {
		X string
	}{
		X:sample,
	}
	var array = []string{"a","b","c"}
	fmt.Println(sample,"=>",SHA512(sample))
	fmt.Println(&sample,"=>",SHA512(&sample))

	fmt.Println(array,"=>",SHA512(array))
	fmt.Println(object,"=>",SHA512(&object))

	fmt.Println(123,"=>",SHA512(123))
}