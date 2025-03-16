package main

import (
	"dfxluna/go-summarize/internal"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	{
		Name: "transcribe-file",
		UsageText: `transcribe-file filepath whisperWebServiceURL
Sends the file at filepath to the whisper-asr-web-service server and prints the transcription to stdout`,
		Args: true,
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 2 {
				return fmt.Errorf("incorrect usage")
			}
			host, err := url.Parse(ctx.Args().Slice()[1])
			if err != nil {
				return err
			}
			host.Port()
			trans, err := internal.NewWhisperWebserviceTranscriber(host)
			if err != nil {
				return err
			}

			result, err := trans.TranscribeFile(ctx.Context, ctx.Args().First())
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	},
}

func main() {
	app := &cli.App{
		Commands: commands,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// func parseArgs(args cli.Args) (string, *url.URL, error) {
// 	if args.Len() != 2 {
// 		return "", nil, cli.Exit(fmt.Sprintf("arg error\nArgs:%v", args.Slice()), exitBadArgs)
// 	}
// 	path, err := filepath.Abs(args.Get(0))
// 	if err != nil {
// 		return "", nil, cli.Exit(err.Error(), exitBadArgs)
// 	}

// 	oaiURL, err := url.Parse(args.Get(1))
// 	if err != nil {
// 		return "", nil, cli.Exit(err.Error(), exitBadArgs)
// 	}
// 	if len(oaiURL.Path) == 0 {
// 		return "", nil, cli.Exit("bad url", exitBadArgs)
// 	}

// 	return path, oaiURL, nil
// }
