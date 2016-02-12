package main

import (
	"fmt"
	"os"
)

func ListenAndServe(addr string) error {
	fmt.Printf("would serve on %s\n", addr)
	return nil
}

func main() {
	err := ListenAndServe("localhost:8200")
	if err != nil {
		fmt.Printf("Error serving: %s", err)
		os.Exit(1)
	}
}
