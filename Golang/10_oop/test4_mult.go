package main

import "fmt"

// 多态的实现(使用接口 )
type Animal interface{
	Sleep()
	Sound() string
}

type Cat struct{
	Color string
}
// 类中重写接口方法即可实现多态
func (this *Cat) Sleep(){
	fmt.Println(this.Color,"Cat is sleeping")
}
func (this *Cat) Sound() string{
	return "喵喵喵"
}

type Dog struct{
	Color string
}

func (this *Dog) Sleep(){
	fmt.Println(this.Color,"dog is sleeping")
}
func (this *Dog) Sound() string{
	return "汪汪汪"
}

func Sleep(a Animal){
	a.Sleep()
}
func Sound(a Animal) string{
	return a.Sound()
}

func main(){
	dog := Dog{"yellow"}
	cat := Cat{"black"}

	cat.Sleep()
	dog.Sleep()
	cat_voice := cat.Sound()
	dog_voice := dog.Sound()
	fmt.Println(cat_voice)
	fmt.Println(dog_voice)
}