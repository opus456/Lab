package main

import (
	"fmt"
	"reflect"
)

// 定义User类和他的方法
type User struct {
	Id   int
	Name string
	Age  int
}

// 测试三种函数
func (this *User) CallUser() {
	fmt.Println("Called user")
}
func (this *User) PrintMsg(msg string) {
	fmt.Print("You received a msg:", msg)
}
func (this *User) Info(t ...string) []interface{} {
	// 用一个slice来存储返回值
	res := make([]interface{}, 0, len(t))
	// 拿到类的反射值--这里Value实际上就是拿到了字段名和值
	val := reflect.ValueOf(this).Elem()
	for _, name := range t {
		field := val.FieldByName(name) // 拿到对应的字段
		if !field.IsValid() {
			res = append(res, nil)
			continue
		}
		res = append(res, field.Interface())
		fmt.Println(field.Type()) // 查看字段名字
	}
	return res
}

func main() {
	usr := User{007, "ljq", 20}
	info := usr.Info("Name", "Age", "xxx")
	for _, i := range info {
		fmt.Println(i)
	}

}
