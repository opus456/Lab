package main

import "fmt"

// 数组
func change(a [10]int){
	for i:=0;i<len(a);i+=1{
		fmt.Print(a[i]," ")
	}
	// 这里的还是属于值传递所以传入的时候不能修改原来数组的值
	a[0] = 100
}

func main(){
	// 数组的初始化是必须要指定长度
	arr := [10] int {}
	arr = [10] int {1,2,3,4}

	// for循环的一种写法
	for i:=0;i<len(arr);i++{
		fmt.Print(arr[i]," ")
	}
	fmt.Print("\n")

	// for 循环的另外一种写法
	for index,value := range(arr){
		fmt.Println(index,value)
	}
	// 或者是如果不需要使用index的话也可以使用_来匿名index
	for _,v := range(arr){
		fmt.Print(v," ")
	}
	fmt.Print("\n")
	change(arr)

}