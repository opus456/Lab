package main

import "fmt"

func printMap(cityMap map[string]string) {
	// 引用传递
	for key, value := range cityMap {
		fmt.Println("Country",key,"`s capital is",value)
	}
}

func changeMap(cityMap map[string]string){
	cityMap["俄罗斯"] = "莫斯科"
}

func main() {
	cityMap := make(map[string]string)
	cityMap["China"] = "Beijing"
	cityMap["America"] = "NewYork"
	cityMap["England"] = "London"

	printMap(cityMap)
	fmt.Print("\n")
	delete(cityMap,"America") // 删除指定的键值对

	changeMap(cityMap)
	printMap(cityMap)

}