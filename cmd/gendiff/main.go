package main

import (
	"code"
	"code/formatters"
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v3"
)

// TODO: !!!!Добавьте в README.md аскинему с примером работы пакета. ++С шага 7
func main() {
	cmd := &cli.Command{
		Name:                   "gendiff",
		Usage:                  "compares two configuration files and shows a difference",
		ArgsUsage:              "<config1> <config2>",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Usage:   "output format",
				Aliases: []string{"f"},
				Value:   "stylish",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() != 2 {
				return fmt.Errorf("only two files required to be compared")
			}

			filepath1, filepath2, format := cmd.Args().Get(0), cmd.Args().Get(1), ""

			if formatters.IsFormatAvailable(cmd.String("format")) {
				format = cmd.String("format")
			} else {
				return formatters.ErrUnsupportedFormat
			}

			ret, err := code.GenDiff(filepath1, filepath2, format)
			if err != nil {
				return err
			}

			fmt.Printf("%s", ret)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
