package main

import "fmt"

type Book struct {
	title string
	auth  string
}

func changeBook (book Book){
	// 值传递所以在这里修改根本没有用
	book.auth = "djaldsa"
}

func changeBook2 (book *Book){
	book.auth = "asjdiasdsa"
}

func main() {
	var book1 Book
	book1.title = "Golang"
	book1.auth = "ljq"
	fmt.Println(book1)
	changeBook(book1)
	fmt.Println(book1)

	changeBook2(&book1)
	fmt.Println(book1)
}