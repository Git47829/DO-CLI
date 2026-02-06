package main

import (
	"fmt"
	"os"

	"mathOps/internal/todo"
	"mathOps/internal/ui"
)

func main() {
	store := todo.NewStore("todo.json")
	if err := ui.Run(store); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
