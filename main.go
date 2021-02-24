package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "keylight",
		Usage: "A simple CLI to control your Elgato Keylights",
		Commands: []*cli.Command{
			switchCommand,
			listCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
