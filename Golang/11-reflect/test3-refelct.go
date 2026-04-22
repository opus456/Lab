package main

import (
	"fmt"
	"reflect"
)

func reflectNum(arg interface{}){
	fmt.Println("type of num is",reflect.TypeOf(arg)) // 直接得到pair里面的type
	fmt.Println("value of num is",reflect.ValueOf(arg)) // 得到pair里面的value
}

func main(){
	// 简单数据类型
	var num float64 = 1.23456
	reflectNum(num)


	
}