package main

import (
	"fmt"
	"reflect"
)

func reflectNum(arg interface{}){
	fmt.Println("type of num is",reflect.TypeOf(arg)) // 直接得到pair里面的type
	fmt.Println("value of num is",reflect.ValueOf(arg)) // 得到pair里面的value
}

type User struct{
	Id int
	Name string
	Age int
}

func (this *User) Call(){
	fmt.Println("User`s function is Called")
}

func (this *User) PrintMsg(msg string){
	fmt.Println("you have received a new message:",msg)
}

// 实现复杂数据结构的reflect 
func FieldAndMethod(user interface{}){
	typeOfUser := reflect.TypeOf(user)
	valueOfUser := reflect.ValueOf(user)
	fmt.Println("type of user is",typeOfUser)
	fmt.Println("value of user is",valueOfUser)

	// 通过反射获取type里面的字段的所有的信息
	for i:=0 ;i<typeOfUser.NumField();i++{
		field:=typeOfUser.Field(i)
		value:=valueOfUser.Field(i).Interface()

		fmt.Printf("%s: %v %v\n",field.Name,field.Type,value)
	}

	// 通过反射获取复杂数据结构中的所有的函数的调用方法
	for i:=0;i<typeOfUser.NumMethod();i++{
		m:=typeOfUser.Method(i)
		fmt.Printf("%s: %v\n",m.Name,m.Type)
		m.Func.Call([]reflect.Value{valueOfUser}) //调用这个Call方法
	}
	method:=valueOfUser.MethodByName("Call")
	method.Call(nil)
}

// 实现用反射机制实现函数方法的调用
func operateFunc(obj interface{},method string,args ...interface{}){
	value:=reflect.ValueOf(obj)
	m:=value.MethodByName(method)
	if(!m.IsValid()){
		fmt.Println("Method not found",method)
		return
	}
	inputs := make([]reflect.Value,len(args))
	for i,arg:=range args{
		inputs[i]=reflect.ValueOf(arg)
	}
	m.Call(inputs)
}

func main(){
	// 简单数据类型
	var num float64 = 1.23456
	reflectNum(num)

	user := User{007, "fengnuan", 18}

	FieldAndMethod(&user)
}