package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/interpreter/cmd"
)

func main() {
	app := &cli.App{
		Name:  "repl",
		Usage: "A REPL for the Monkey programming language",
		Action: func(ctx *cli.Context) error {
			user, err := user.Current()
			if err != nil {
				return err
			}

			fmt.Printf("Hello %s! This is the Monkey programming language REPL!\n", user.Username)
			fmt.Printf("Feel free to type in Monkey statements!\n")

			return cmd.StartRepl(os.Stdin, os.Stdout)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
