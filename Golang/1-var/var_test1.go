// 各种的变量的声明方式
package main

import (
	"fmt"
)

// 只有用var声明的变量才可以用作全局变量
var g int
var gg int = 999
var ggg = 999

func main() {

	// 方法1: 声明变量 , 默认是0
	var a int
	fmt.Println("a =", a)
	fmt.Printf("type of a = %T \n", a)

	// 方法2: 声明的时候声明数据类型并且赋值
	var b int = 100
	fmt.Println("b =", b)
	fmt.Printf("type of b = %T\n", b)

	// 方法3: 声明的时候省略数据类型会根据值进行判断
	var c = "daldsad"
	fmt.Println("c =", c)
	fmt.Printf("type of c is %T", c)

	// 方法4(最常用但是不能申明全局的变量): :=
	d := 100
	e := ":="
	fmt.Println("d = ", d)
	fmt.Println("e = ", e)
	fmt.Println("type of d is", d)
	fmt.Println("type of e is", e)

	// =====
	fmt.Println("g = ", g)
	fmt.Printf("type of g is %T \n", g)
	fmt.Println("gg= ", gg)
	fmt.Printf("type of gg is %T \n", gg)
	fmt.Println("ggg = ", ggg)
	fmt.Printf("type of ggg is %T \n", ggg)

	// 还有一种方法
	var(
		var_a int = 100
		var_b string = "dadsada"
	)
	fmt.Println(var_a)
	fmt.Println(var_b)

	// 注意在go里面变量如果在声明的时候没有给值都是会有默认的返回值的
	var(
		string_s string
		int_i int
		bool_b bool
	)
	fmt.Println("string_s  = ",string_s)
	fmt.Println("int_i  = ",int_i)
	fmt.Println("bool_b  = ",bool_b)
}
