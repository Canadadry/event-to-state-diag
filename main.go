package main

import (
	"app/cmd/diag"
	"app/cmd/matrix"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "failed", err)
		os.Exit(1)
	}
}

func run(args []string) error {

	actions := map[string]func([]string) error{
		diag.Name:   diag.Run,
		matrix.Name: matrix.Run,
	}
	if len(args) == 0 || actions[args[0]] == nil {
		return fmt.Errorf("app action [args]\n actions:%s,%s", diag.Name, matrix.Name)
	}

	return actions[args[0]](args[1:])
}
