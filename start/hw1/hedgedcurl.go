package main

import (
	"context"
	"flag"
	"fmt"
	"hedgedcurl/service"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	var timeout int

	flag.IntVar(&timeout, "timeout", 15, "Таймаут запросов в секундах")
	flag.IntVar(&timeout, "t", 15, "Таймаут запросов в секундах (сокращенно)")
	flag.Parse()

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	urls := flag.Args()
	sender := service.NewSender(urls)
	responseChannel := sender.SendRequests(ctx)
	select{
	case response := <- responseChannel:
		httpStatus := response.StatusCode
		var headers strings.Builder
		for key, val := range response.Header {
			for _, v := range val {
				headers.WriteString(fmt.Sprintf("%s: %s\n", key, v))
			}
		}
		headersStr := headers.String()
		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		fmt.Printf("Полученный код ответа: %d заголовки: %s\n Тело: \n%s", httpStatus, headersStr, string(body))
		
		
	case <- ctx.Done():
		fmt.Print("Timeout Exceed")
		os.Exit(228)
	}
}
