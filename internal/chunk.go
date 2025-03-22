package internal

import (
	"github.com/samber/lo"
	"github.com/sentencizer/sentencizer"
)

const (
	en = "en"
)

type Chunker interface {
	Chunk(text string, opts ChunkOpts) ([]Chunk, error)
}

type ChunkOpts struct {
	Overlap int
}

type Chunk struct {
	Sentences []string
}

type EnglishChunker struct {
	segmenter sentencizer.Segmenter
}

func NewEnglishChunker() *EnglishChunker {
	return &EnglishChunker{
		segmenter: sentencizer.NewSegmenter(en),
	}
}

func (c *EnglishChunker) Chunk(text string, opts ChunkOpts) ([]Chunk, error) {
	sentences := c.segmenter.Segment(text)
	chunks := lo.Map(sentences, func(_ string, i int) Chunk {
		lowerbound := i
		upperbound := i
		if i-opts.Overlap >= 0 {
			lowerbound = i - opts.Overlap
		}
		if i+opts.Overlap <= len(sentences) {
			upperbound = i + opts.Overlap + 1
		}
		return Chunk{
			Sentences: sentences[lowerbound:upperbound],
		}
	})
	return chunks, nil
}
