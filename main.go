package main

import (
	"fmt"
)

const requestURL = "https://api.github.com/users/"

func main() {
	var username string
	// design a stdin and stdout interface.
	// So that a user could can cat input.txt | go run main.go > output.txt
	// and the program will run and output the result to output.txt
	// This is useful for testing.
	fmt.Println("Enter username: ") // breaks jq pipe
	fmt.Scanf("%s", &username)
	var data GitHubAPI
	if err := data.GETUserData(username); err != nil {
		fmt.Println(err)
		return
	}
	AnyToJSON(data)
}
