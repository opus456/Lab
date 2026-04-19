package main

import(
	"fmt"
)

// 函数以及其的返回值

// 方法1: 指定返回的参数的类型如果没有返回类型的话可以直接不用指定返回类型
func foo1(a string , b int) int{
	fmt.Println("a = ",a)
	fmt.Println("b = ",b)
	
	c := 100
	return c
}

// 方法2: 返回多个参数--匿名
func foo2(a int ,b int) (int ,int){
	return a,b
}

// 方法3: 返回多个参数--指定名字
func foo3(a int,b int) (r1 int, r2 int){
	fmt.Println("---function3---")
	
	// 到直接给返回的变量赋值
	r1 = a
	r2 = b
	return 
}

func foo4(a int, b int ) (r1,r2 int){
	fmt.Println("---foo4---")
	r1 = a
	r2 = b
	return r1,r2
}

func main(){
	var s string = "hello go"
	var a int = 99
	var b = 520
	r1 := foo1(s,a)
	aa,bb := foo2(a,b)

	fmt.Println("function1 return value is ",r1)
	fmt.Println("function2 return value is ",aa, "function2 return value2 is",bb)
	aaa,bbb := foo3(a,b)
	fmt.Println(aaa,bbb)
	aaaa,bbbb := foo4(a,b)
	fmt.Println(aaaa,bbbb)
}	