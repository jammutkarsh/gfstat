package main

import (
	"fmt"
	"os"
)

const requestURL = "https://api.github.com/users/"

func main() {
	var username string
	// design a stdin and stdout interface.
	// So that a user could can cat input.txt | go run main.go > output.txt
	// and the program will run and output the result to output.txt
	// This is useful for testing.
	// fmt.Println("Enter username: ") // breaks jq pipe
	fmt.Scanf("%s", &username)
	var (
		data GitHubAPI
		c    CurrentUser
	)
	if err := data.GETUserData(username); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.setMetadata(data.Username, data.HTMLURL, data.Type)
	if err := data.GETFollowing(&c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := data.GETFollowers(&c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	AnyToJSON(c)
}
