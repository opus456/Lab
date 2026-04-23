package main

import (
	"fmt"
	"reflect"
)

// 对于简单的类型可以拿到他的值和类型
func main(){
	var a float64 = 1.23456

	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	fmt.Println("type of a is",t)
	fmt.Println("value of a is",v)

}