package main

import "fmt"

func main() {
	// slice的一些操作追加和截取
	var s []int = make([]int, 3, 5) // 这里可以指定slice的容量capability
	fmt.Println("s = ", s, "len = ", len(s), "cap = ", cap(s))
	// 追加新的元素
	s = append(s, 1, 2)
	fmt.Println("s = ", s, "len = ", len(s), "cap = ", cap(s))

	// 如果在cap已经满了的情况下再append新的元素的话会直接让cap = 2*cap
	s = append(s, 3)
	fmt.Println("s = ", s, "len = ", len(s), "cap = ", cap(s))

	// 截取
	fmt.Println("====================================")
	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[0:2] // s2和s1是共享的同一个内存
	s2[0] = 100
	fmt.Println(s1[:2]) //[0:2)
	fmt.Println(s2)     // 发现和s1的前面两个元素是一样的
	// 如果想要开辟一个新的空间的话使用copy
	s_c := make([]int, 3)
	copy(s_c, s1)
	fmt.Println(s_c)
	// 常用的一个操作删除指定位置的元素
	x := []int{1, 2, 3, 4, 5}
	i := 2
	copy(x[i:], x[i+1:]) // 这一段实现{3,4,5}->{4,5,5}也就是说用后面的元素把前面的这个一个给覆盖掉了
	x = x[:len(x)-1]     // 这里再删除最后的一个多余的元素就可以保证删除指定元素
	fmt.Println(x)
}
