package main

import (
	"fmt";

)

// pair是指对于一个interface变量无论存储什么东西, 底层都是由两个部分组成(value,type)也就是说interface是一个包含两个指针的结构体
// 对于static变量pair不会起作用
// 对于concrete变量类型会跟着pair走

func main(){
	var i interface{}
	i = 10
	fmt.Printf("type of i is %T\n",i)
	i = "Golang"
	fmt.Printf("type of i is %T\n",i) // 这里会发现没有报错看似可以直接转换过来数据类型,但是其实i一直是interface只是interface里面的pair变了

	// example2
	var a string
	// pair<statictype:string,value:"ljq"
	a = "ljq"

	var t interface{}
	t = a
	str,ok := t.(string)
	if(ok){
		fmt.Println(str)
	}else {
		fmt.Printf("type of t is %T\n",t)
		fmt.Printf("type of str is %T\n",str)
	}

}