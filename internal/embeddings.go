package internal

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/JexSrs/go-ollama"
)

const (
	success = "success"
)

var (
	ErrModelPull         = errors.New("could not pull model")
	ErrVectorLenMismtach = errors.New("v1 and v2 are not the same size")
)

type Embedder interface {
	GetEmbeddings(ctx context.Context, text string) ([]float64, error)
}

type OllamaEmbedder struct {
	ollamaUrl    *url.URL
	model        string
	ollamaClient *ollama.Ollama
}

func NewOllamaEmbedder(ctx context.Context, ollamaURL *url.URL, model string) (*OllamaEmbedder, error) {
	oc := ollama.New(*ollamaURL)

	pmr, err := oc.Models.Pull(oc.Models.Pull.WithModel(model))
	if err != nil {
		return nil, err
	}
	if !strings.Contains(pmr.Status, success) {
		return nil, fmt.Errorf("%w: %s", ErrModelPull, pmr.Status)
	}

	return &OllamaEmbedder{
		ollamaUrl:    ollamaURL,
		model:        model,
		ollamaClient: oc,
	}, nil
}

func (e *OllamaEmbedder) GetEmbeddings(ctx context.Context, text string) ([]float64, error) {
	embeds, err := e.ollamaClient.GenerateEmbeddings(
		e.ollamaClient.GenerateEmbeddings.WithModel(e.model),
		e.ollamaClient.GenerateEmbeddings.WithPrompt(text),
		e.ollamaClient.GenerateEmbeddings.WithKeepAlive("20m"),
	)
	if err != nil {
		return []float64{}, err
	}
	return embeds.Embedding, nil
}
