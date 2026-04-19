package main

import "fmt"

type Human struct {
	name string
	sex  string
}

func (this *Human) ShowInfo(){
	fmt.Println("name :",this.name)
	fmt.Println("sex :",this.sex)
}

func (this *Human) Walk() {
	fmt.Println(this.name,"is Walking")
}

func (this *Human) Eat() {
	fmt.Println(this.name,"is Eating")
}

type SuperMan struct{
	Human
	level int
}

func (this *SuperMan) ShowInfo(){
	fmt.Println("name :",this.name)
	fmt.Println("sex :",this.sex)
	fmt.Println("level :",this.level)
}

func (this *SuperMan) LevelUp(){
	this.level+=1
	fmt.Println("SuperMan",this.name,"`s level ++ ,now his/her level is",this.level)
}

func main() {
	human := Human{name:"zhang3",sex:"male"}
	human.ShowInfo()
	human.Walk()
	human.Eat()
	fmt.Println("=============================")
	// ==================
	superMan := SuperMan{
		Human: Human{name:"beauty",sex:"female"},
		level:10,
	}
	superMan.ShowInfo()
	superMan.LevelUp()
	superMan.Eat()
}	