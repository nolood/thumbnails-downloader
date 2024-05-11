package main

import (
	thumbsgrpc "cli-thumbs/internal/clients/thumbs/grpc"
	"cli-thumbs/internal/config"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Грязюка

func main() {
	cfg := config.MustLoad()

	thumbsClient, err := thumbsgrpc.New(context.Background(), cfg.Clients.Thumbs.Address, cfg.Clients.Thumbs.Timeout, cfg.Clients.Thumbs.RetriesCount)
	if err != nil {
		panic(err)
	}

	async := flag.Bool("async", false, "You can download a lot of files concurrently")

	flag.Parse()

	arguments := flag.Args()

	if len(arguments) == 0 {
		fmt.Println("Please provide at least one URL.")
		return
	}

	if *async {
		var wg sync.WaitGroup
		for _, link := range arguments {
			wg.Add(1)
			go func(link string, wg *sync.WaitGroup) {
				defer wg.Done()
				log.Println(link)
				url, err := thumbsClient.Download(context.Background(), link)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					return
				}
				downloadAndSave(url)
			}(link, &wg)
		}
		wg.Wait()
		return
	}

	for _, link := range arguments {
		url, err := thumbsClient.Download(context.Background(), link)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		downloadAndSave(url)
	}
}

func downloadAndSave(url string) {
	fileName := "thumb_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
		return
	}

	file, err := os.Create("./images/" + fileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Error copying content to file: %v\n", err)
		return
	}
}
