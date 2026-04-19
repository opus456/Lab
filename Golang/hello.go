package main
import(
	"fmt"
)
// 对三种Print进行比较
// main函数
func main(){
	var a,b = 99,"I wanna study Golang"
	fmt.Print("hello" , "Golang") // 直接存粹的输出 , 不会自动换行 , 两个参数之间也不会自己加空格
	fmt.Print("\n")

	// Printf 是可以自己手动高度定制格式化的, %v就输出原文,如果是字符串的话可以用%s , 输出类型的话使用%T
	fmt.Printf("a = %v\n",a) 
	fmt.Printf("type of a is %T \n",a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("b = %s\n", b)

	fmt.Println("a = ",a, "b = ",b) //会在最后自动换行光标移动到下一行 , 并且两个参数之间会自动的加一个空格
}