package main

import "fmt"

// 定义一个类
// 注意首字母大写说明可以暴露给其他的包
type Hero struct{
	Name string
	Ad int 
	Level int
}

/*func (this Hero) Show(){
	fmt.Println("Name is",this.Name)
	fmt.Println("Ad is",this.Ad)
	fmt.Println("Level is",this.Level)
}

func (this Hero) GetName() string{
	return this.Name
}

func (this Hero) SetName(newName string){
	// 这里不能修改因为这是值传递只是拷贝了一个副本
	this.Name = newName
} */

func (this *Hero) Show(){
	fmt.Println("Name is",this.Name)
	fmt.Println("Ad is",this.Ad)
	fmt.Println("Level is",this.Level)
}

func (this *Hero) GetName() string{
	return this.Name
}

func (this *Hero) SetName(newName string){
	// 这里可以修改
	this.Name = newName
}

func main(){
	hero := Hero{"ljq",100,10}
	hero.Show()

	newName := "fengnuan"
	hero.SetName(newName) 

	heroName := hero.GetName()
	fmt.Println(heroName)
}