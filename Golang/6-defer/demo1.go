package main

import "fmt"

// Golang中的defer关键字的使用和教学
// defer的函数的执行是用栈空间存储的

func func1(){
	fmt.Println("func1 over")
}

func func2(){
	fmt.Println("func2 over")
}

func deferFunc() int{
	fmt.Println("defer function called..")
	return 0
}

func returnFunc() int{
	fmt.Println("return function called..")
	return 0
}

// 对比defer 和return的先后 通过结果可以看到defer是最后被执行的还是在return的后面被执行
func returnAndDefer() int{
	defer deferFunc()
	return returnFunc()
}
func main(){
	
	// defer func1() // 这里的func1会被后执行因为defer的处理是压栈也就是说明在func1是在栈底
	// defer func2() // func2在栈顶

	fmt.Println("test defer keyword")

	returnAndDefer()

}