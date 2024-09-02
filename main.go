package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// Read audio as bytes
// Make an api call

const (
	exitBadArgs = 1
)

func main() {
	app := &cli.App{
		Name: "gsm",
		Usage: `gsm filePath oaiURL
filePath is a path to an audio file to summarize
url is a url that points to an openAI api server`,
		Args: true,
		Action: func(ctx *cli.Context) error {
			path, oaiURL, err := parseArgs(ctx.Args())
			if err != nil {
				return err
			}
			fmt.Printf("Using filepath: %s\nUsing oaiURL: %s\n", path, oaiURL.Path)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseArgs(args cli.Args) (string, *url.URL, error) {
	if args.Len() != 2 {
		return "", nil, cli.Exit(fmt.Sprintf("arg error\nArgs:%v", args.Slice()), exitBadArgs)
	}
	path, err := filepath.Abs(args.Get(0))
	if err != nil {
		return "", nil, cli.Exit(err.Error(), exitBadArgs)
	}
	oaiURL, err := url.Parse(args.Get(1))
	if err != nil {
		return "", nil, cli.Exit(err.Error(), exitBadArgs)
	}
	if len(oaiURL.Path) == 0 {
		return "", nil, cli.Exit("bad url", exitBadArgs)
	}

	return path, oaiURL, nil
}
