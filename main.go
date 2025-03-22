package main

import (
	"dfxluna/go-summarize/internal"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/samber/lo"
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
	{
		Name: "summarize-file",
		UsageText: `summarize-file filepath whisperWebServiceURL ollamaURL
Sends the file at filepath to the whisper-asr-web-service server, chunks the resulting transcription,
converts the chunks into embeddings using the ollama server, and uses the resulting embeddings to aid in 
summarization of the chunks
`,
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

			fmt.Println("Transcribing ", ctx.Args().First())
			result, err := trans.TranscribeFile(ctx.Context, ctx.Args().First())
			if err != nil {
				return err
			}
			chunker := internal.NewEnglishChunker()
			chunks, err := chunker.Chunk(result, internal.ChunkOpts{Overlap: 1})
			if err != nil {
				return err
			}

			// lo.ForEach(chunks, func(chunk internal.Chunk, i int) {
			// 	fmt.Printf("Chunk %d: %+v\n", i, chunk)
			// })
			url, err := url.Parse("http://192.168.1.6:11434")
			if err != nil {
				return err
			}

			em, err := internal.NewOllamaEmbedder(ctx.Context, url, "nomic-embed-text")
			if err != nil {
				return err
			}

			embeddings := lo.Map(chunks, func(c internal.Chunk, i int) []float64 {
				fmt.Println("Getting embeddings for chunk ", i)
				embedding, err := em.GetEmbeddings(ctx.Context, strings.Join(c.Sentences, " "))
				if err != nil {
					panic(err)
				}
				return embedding
			})
			fmt.Printf("found  %d embeddings\n", len(embeddings))
			distances := make([][]float64, len(embeddings))

			lo.ForEach(embeddings, func(e []float64, i int) {
				distances[i] = make([]float64, len(embeddings))
				lo.ForEach(embeddings, func(f []float64, j int) {
					distances[i][j], err = internal.CosineDistance(embeddings[i], embeddings[j])
					if err != nil {
						panic(err)
					}
				})
				fmt.Printf("Row %d: %.2f\n", i, distances[i])
			})

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
