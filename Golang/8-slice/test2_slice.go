package  main

import "fmt"

func test(a []int){
	for _,v := range(a){
		fmt.Print(v," ")
	}
	fmt.Print("\n")
	// 这里是引用传递可以直接修改
	a[0] = 100
}

func slice(){
	// 动态数组只需要不指定长度就可以
	// arr := []int {1,2,3,4}
	// fmt.Println(len(arr))
	// for _,v := range(arr){
	// 	fmt.Print(v," ")
	// }
	// fmt.Print("\n")
	// test(arr)

	// fmt.Println(arr)

	// 声明动态数组的四种方法
	// 1.直接使用:=声明如果后面没有值的话就会
	// arr1 := []int {1,2,3,4}
	// 这里必须要一个{}因为:=是要给他赋值一个东西(这里是一个空切片)不能是nil空
	// arr1 := []int{}
	
	// 2. 使用var进行赋值--声明了一个slice但是并没有任何的空间这种声明是nil
	var arr2 []int
	// arr2 = make([]int, 3) 

	// 3. 直接第二种的一步到位
	// var arr3 []int = make([]int,3)
	
	// 4. 最常用的
	// arr4 := make([]int,3)
	// fmt.Println("len is ",len(arr4),"arr4 is ",arr4)
	if(arr2==nil){
		fmt.Println(arr2, "is nil")
	} else{
		fmt.Println(arr2)
	}
}