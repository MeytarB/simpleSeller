package main

import (
	"fmt"
	"simpleSeller/handlers"
)

func main() {
	fmt.Println("hello")
	simpleServer := handlers.NewHandler()
	fmt.Println("Server listening on port 8080")
	simpleServer.InitServer()

}
