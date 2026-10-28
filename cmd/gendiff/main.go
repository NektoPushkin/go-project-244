package main

import (
	"code"
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v3"
)

// TODO: !!!!Добавьте в README.md аскинему с примером работы пакета.
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
				return fmt.Errorf("both files required to be compared")
			}

			filepath1, filepath2, format := cmd.Args().Get(0), cmd.Args().Get(1), ""

			switch cmd.String("format") {
			case "stylish":
				fallthrough
			case "plain":
				format = cmd.String("format")
			default:
				return code.ErrUnsupportedFormat
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
