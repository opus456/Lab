package main

import ("fmt")


// 这里是函数的闭包使用方法也就是说在外部函数定义的变量相当于是变得半永久了内部函数调用之后的状态会被存储起来
func counter() func() int{
	count := 0
	return func() int{
		count++
		return count
	}
}

func main(){
	c := counter()
	fmt.Println(c())
	fmt.Println(c())
	fmt.Println(c())

}