package main

import (
	"fmt"
)

// func swap(a int, b int){
// 	var temp int
// 	temp = a
// 	a = b
// 	b = temp
// }

func swap(a *int, b *int) {
	var temp int
	temp = *a
	*a = *b
	*b = temp
}

func main() {
	a := 10
	b := 20

	swap(&a, &b)
	fmt.Println("a = ", a)
	fmt.Println("b = ", b)

	// 指针指向地址
	var p *int
	p = &a
	fmt.Println(p)
	fmt.Println(&a)

	// 二级指针
	var pp **int
	pp = &p
	fmt.Println(pp)
	fmt.Println(&p)
}
