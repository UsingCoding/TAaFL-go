package main

import (
	"github.com/urfave/cli/v2"
	stdlog "log"
	"os"
)

func main() {
	err := runApp(os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(args []string) error {
	app := cli.App{
		Name: "slr-runner",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lexer",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
			},
			&cli.StringFlag{
				Name:    "grammar",
				Aliases: []string{"g"},
			},
		},
		EnableBashCompletion: true,
		Action:               executeAction,
	}

	return app.Run(args)
}
