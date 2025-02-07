package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

var VERSION_HEADERS = map[string]string{
	"editor-version":       "Neovim/0.10.4",
	"ditor-plugin-version": "lazycopilot/v0.0.1",
	"User-Agent":           "lazycopilot/v0.0.1",
	"sec-fetch-site":       "none",
	"sec-fetch-mode":       "no-cors",
	"sec-fetch-dest":       "empty",
	"priority":             "u=4, i",
}

type Headers map[string]string

type HttpOptions struct {
	Url     string
	Method  string
	Headers *Headers
	Body    interface{}
}

type HttpResponse struct {
	*http.Response
}

func HttpRequest(ctx context.Context, opts HttpOptions) (*HttpResponse, error) {
	var body io.Reader
	if opts.Body != nil {
		data, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, opts.Url, body)
	if err != nil {
		return nil, err
	}

	if opts.Headers != nil {
		for key, value := range *opts.Headers {
			req.Header.Set(key, value)
		}
	}
	for key, value := range VERSION_HEADERS {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	httpResp := &HttpResponse{res}

	return httpResp, nil
}

func (r *HttpResponse) StringDecode() (string, error) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *HttpResponse) JsonDecode(v any) error {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
