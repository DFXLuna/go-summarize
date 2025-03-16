package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
)

const (
	whisperEndpoint = "asr"
	audioFile       = "audio_file"
)

type Transcriber interface {
	TranscribeFile(ctx context.Context, filepath string) (string, error)
}

type whisperASRWebservice struct {
	Endpoint *url.URL
}

func NewWhisperWebserviceTranscriber(host *url.URL) (*whisperASRWebservice, error) {
	endpoint := host.JoinPath(whisperEndpoint)
	return &whisperASRWebservice{
		Endpoint: endpoint,
	}, nil
}

func (w *whisperASRWebservice) TranscribeFile(ctx context.Context, filepath string) (string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var bs bytes.Buffer
	form := multipart.NewWriter(&bs)
	field, err := form.CreateFormFile(audioFile, path.Base(filepath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(field, file)
	if err != nil {
		return "", err
	}
	form.Close()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.Endpoint.String(), &bs)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-type", form.FormDataContentType())

	q := req.URL.Query()
	q.Add("output", "text")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("error making request: Status code: %d, status: %s", resp.StatusCode, resp.Status)
	}

	respBs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBs), nil
}
