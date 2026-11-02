package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Sender struct{
	urls [] string
}

func NewSender(sliceUrls []string) *Sender {
	return &Sender{
		urls: sliceUrls,
	}
}



func (s *Sender) SendRequests(ctx context.Context) <- chan http.Response{
	out := make(chan http.Response, 1)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	for _, url := range s.urls {
		go func(targetURL string){
			response, err := httpGetRequest(ctxWithCancel, targetURL)
			if err != nil {
				return
			}
			body, err := io.ReadAll(response.Body)
			response.Body.Close()
			response.Body = io.NopCloser(bytes.NewReader(body))
			if err != nil {
				return
			}
			select {
			case out <- *response:
				cancel()
			default:
				response.Body.Close()
			}
		}(url)
	}
	return out
}

func httpGetRequest(ctx context.Context, url string) (*http.Response, error){
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Ошибка на этапе создания HTTP запроса: %w", err )
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		if errors.Is(err, context.Canceled){
			return nil, fmt.Errorf("Ошибка отмены запроса: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded){
			return nil, fmt.Errorf("Ошибка таймаута запроса: %w", err)
		}
		return nil, fmt.Errorf("Ошибка HTTP запроса: %w", err)
	}
	return response, nil
}

