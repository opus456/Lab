package main

import "fmt"

// 空接口数据类型可以用来作为any可以占位用来当int string float struct等

func myFunc(arg interface{}){
	fmt.Printf("arg is ",arg," it`s type is %T ",arg,"\n")

	// 断言用来区分数据结构到底是什么
	value,ok := arg.(string)
	if(!ok){
		fmt.Printf("arg`type is not string but %T ",arg ,"\n")
	} else{
		fmt.Println("arg is string")
		fmt.Println(value)
	}
}
type Man struct{
	name string
}

func main(){
	book:=Man{"fengnuan"}
	myFunc(book)
	myFunc("fengnuan")
	myFunc(100)
	myFunc('c')
}