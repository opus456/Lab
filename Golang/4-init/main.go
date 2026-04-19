package main

import (
	"fmt"
	_ "Golang/4-init/lib1" // 这种写法可以导入这个lib1这个包但是可以不使用只执行init
	mylib2 "Golang/4-init/lib2" // 相当于给这个lib2这个包取一个别名
)

func main(){
	// a := 10
	s := "init function"
	// lib1.Lib1Test(a)
	r2:= mylib2.Lib2Test(s)
	fmt.Println("lib2 return value is ",r2)
}