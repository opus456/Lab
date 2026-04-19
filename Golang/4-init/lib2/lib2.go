package lib2

import "fmt"

// lib2对外的接口函数
func Lib2Test(s string) string{
	fmt.Println("Lib2 Testing ... ")
	return s
}

func init(){
	fmt.Println("Lib2 initing ... ")
}