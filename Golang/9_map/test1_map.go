package main

import "fmt"

// Golang中的map的声明和定义的方式

func main(){

	// 第一种声明方式
	var myMap1 map[int]string // [key]value
	if(myMap1 == nil){
		fmt.Println("myMap1 is nil: ",myMap1)
	}
	// 需要给map分配空间才可以 
	myMap1 = make(map[int]string,3)
	myMap1[1] = "python"
	myMap1[2] = "c++"
	myMap1[3] = "javascript"
	fmt.Println(myMap1)
	// 如果超出分配的空间的话会自己扩容
	myMap1[4] = "Golang"
	fmt.Println(myMap1)

	// 第二种声明方式在声明的时候直接分配内存
	myMap2 := make(map[int]string) // 这可以不用分配大小 , 会自己带上几个大小
	mySlice := make([]int,3) // 但是slice就需要手动的分配大小
	mySlice[0] = 1
	fmt.Println(mySlice)
	for i:=1;i<=3;i++{
		myMap2[i] = "i"
	}
	fmt.Println(myMap2)

	// 第三种声明方式直接声明的时候赋值
	myMap3 := map[int]string{
		1:"python",
		2:"Golang",
		3:"javaScript",
	}
	for key,v := range myMap3{
		fmt.Println(key,v)
	}
	fmt.Println(myMap3[1])
}