package lib1

import "fmt"

// lib1的对外的方法函数, 函数名首字母大写表示对外开放函数
func Lib1Test(a int){
	fmt.Println("Lib1 value is ",a)
}

func init(){
	fmt.Println("Lib1 initing ... ")
}