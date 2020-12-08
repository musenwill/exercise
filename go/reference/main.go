package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	users := []User{User{Name: "joe", Age: 18}}
	for _, u := range users {
		u.Age = 20
	}
	fmt.Println(users)
}
